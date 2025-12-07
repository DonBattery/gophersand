// cell_automata.go provides functions for working with cellular automata
package game

import (
	"image"
	"math"
	"math/rand"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

type CellAutomata struct {
	// The size of the cell automata in pixels
	Width  int
	Height int
	Size   int

	// The size of the cell automata in grid cells (used for space partitioning)
	CellSize   int
	GridWidth  int
	GridHeight int
	GridSize   int

	// Active grid cell management
	ActiveCells         []int
	ActiveCellFlags     []bool
	NextActiveCells     []int
	NextActiveCellFlags []bool

	// Tick and Processed ensures that we only update a pixel once per frame
	Tick      int
	Processed []int

	// Processors for each material type
	MaterialProcessors []func(ca *CellAutomata, pxId, x, y int) bool

	// Per pixel data
	Materials   []Material
	Pixels      []byte   // 4 byte per pixel (RGBA)
	CellBuffers [][]byte // per cell pixel buffers

	// Ebiten image for rendering
	Img     *ebiten.Image
	ImgOpts *ebiten.DrawImageOptions
}

// NewCellAutomata creates a new cellular automata, based on pixel width, height and cellSize
// the grid is used to track active cells
func NewCellAutomata(pxWidth, pxHeight, cellSize int) *CellAutomata {
	size := pxWidth * pxHeight
	gridW := int(math.Ceil(float64(pxWidth) / float64(cellSize)))
	gridH := int(math.Ceil(float64(pxHeight) / float64(cellSize)))
	gridSize := gridW * gridH

	ca := &CellAutomata{
		Width:    pxWidth,
		Height:   pxHeight,
		Size:     size,
		CellSize: cellSize,

		GridWidth:  gridW,
		GridHeight: gridH,
		GridSize:   gridSize,

		ActiveCells:         []int{},
		ActiveCellFlags:     make([]bool, gridSize),
		NextActiveCells:     []int{},
		NextActiveCellFlags: make([]bool, gridSize),

		Tick:      0,
		Processed: make([]int, size),

		MaterialProcessors: make([]func(ca *CellAutomata, pxId int, x int, y int) bool, 256),

		Materials:   make([]Material, size),
		Pixels:      make([]byte, size*4),
		CellBuffers: make([][]byte, gridSize),

		Img:     ebiten.NewImage(pxWidth, pxHeight),
		ImgOpts: &ebiten.DrawImageOptions{},
	}

	ca.MaterialProcessors[MaterialSand.GetType()] = ProcessSand
	ca.MaterialProcessors[MaterialWater.GetType()] = ProcessWater
	ca.MaterialProcessors[MaterialSmoke.GetType()] = ProcessSmoke

	baseBuffer := make([]byte, cellSize*cellSize*4)

	// for each grid cell, assign a cell buffer for rendering pixels
	for cid := 0; cid < gridSize; cid++ {
		// pixel bounds of this cell (cells on the right side and on the bottom may not be default sized)
		startX, startY := ca.CellTopLeftPxPos(cid)
		endX := min(startX+ca.CellSize-1, ca.Width-1)
		endY := min(startY+ca.CellSize-1, ca.Height-1)
		cellW := endX - startX + 1
		cellH := endY - startY + 1

		// use the base buffer by default (it is not an issue that we use this for many cells, as they will be processed sequentially)
		buff := baseBuffer

		// if the cell is not the default size, create a fitting buffer for it
		if cellW != cellSize || cellH != cellSize {
			buff = make([]byte, cellW*cellH*4)
		}
		ca.CellBuffers[cid] = buff
	}

	return ca
}

// Update updates the cellular automata
func (ca *CellAutomata) Update() {
	ca.Tick++

	// the cells we activated in the previous frame, became our active cells this frame
	ca.ActiveCells = ca.NextActiveCells
	ca.ActiveCellFlags = ca.NextActiveCellFlags
	// clear next active cells, so we can start to prepare for the next frame
	ca.NextActiveCells = []int{}
	ca.NextActiveCellFlags = make([]bool, ca.GridSize)

	// early exit if no active cells
	l := len(ca.ActiveCells)
	if l == 0 || !IsRunning {
		return
	}

	// shuffle active cells to avoid bias
	rand.Shuffle(l, func(i, j int) {
		ca.ActiveCells[i], ca.ActiveCells[j] = ca.ActiveCells[j], ca.ActiveCells[i]
	})

	// cache properties for fast access
	processors := ca.MaterialProcessors
	materials := ca.Materials
	processed := ca.Processed
	tick := ca.Tick

	// iterate over the active cells
	for _, cellId := range ca.ActiveCells {

		// pixel bounds of this cell (cells on the right and bottom edge of the grid may be not default sized)
		startX, startY := ca.CellTopLeftPxPos(cellId)
		endX := min(startX+ca.CellSize-1, ca.Width-1)
		endY := min(startY+ca.CellSize-1, ca.Height-1)

		// randomize direction to avoid bias
		dir := 1
		if rand.Intn(2) == 0 {
			dir = -1
		}

		active := false

		// check from bottom to top
		for yy := endY; yy >= startY; yy-- {

			// either check from left to right or right to left depending on direction
			var xxStart, xxEnd int
			if dir == 1 {
				xxStart, xxEnd = startX, endX
			} else {
				xxStart, xxEnd = endX, startX
			}

			for xx := xxStart; (dir == 1 && xx <= xxEnd) || (dir == -1 && xx >= xxEnd); xx += dir {

				pxId, ok := ca.PxId(xx, yy)

				// if pixel is already processed in this tick, skip it
				if !ok || processed[pxId] == tick {
					continue
				}

				processor := processors[materials[pxId].GetType()]

				// if the material does not have a processor, skip it
				if processor == nil {
					continue
				}

				// process pixel, if activity is detected mark the cell as active, and the pixel as processed
				if processor(ca, pxId, xx, yy) {
					active = true
					processed[pxId] = tick
				}
			}
		}

		if active {
			// mark the cell and neighbors active for next frame
			ca.ActivateNeighbors(cellId)
		}
	}
}

// Draw the cellular automata on to a target image, with an optional debug grid showing the active cells
func (ca *CellAutomata) Draw(target *ebiten.Image, debug bool) {
	buffers := ca.CellBuffers

	// Helper that copies a single grid cell from ca.Pixels → ca.Img.
	drawCell := func(cellId int) {
		// pixel bounds of this cell
		startX, startY := ca.CellTopLeftPxPos(cellId)
		endX := min(startX+ca.CellSize-1, ca.Width-1)
		endY := min(startY+ca.CellSize-1, ca.Height-1)

		cellW := endX - startX + 1

		buff := buffers[cellId]

		buffLineLen := cellW * 4
		worldLineLen := ca.Width * 4

		// world index for the first row
		worldStart := (startY*ca.Width + startX) * 4
		buffIndex := 0

		// Copy cell from ca.Pixels → buff (top → bottom this time)
		for yy := startY; yy <= endY; yy++ {
			copy(buff[buffIndex:buffIndex+buffLineLen],
				ca.Pixels[worldStart:worldStart+buffLineLen])

			buffIndex += buffLineLen
			worldStart += worldLineLen
		}

		// push buff into the ebiten image
		rect := image.Rect(startX, startY, endX+1, endY+1)
		sub := ca.Img.SubImage(rect).(*ebiten.Image)
		sub.WritePixels(buff)
	}

	// When the simulation is not running we still want to render the
	// entire world every frame so that any drawing / generation changes
	// are visible immediately, while CA.Update() remains paused.
	if !IsRunning {
		for cellId := 0; cellId < ca.GridSize; cellId++ {
			drawCell(cellId)
		}
	} else {
		for _, cellId := range ca.ActiveCells {
			drawCell(cellId)
		}
	}

	target.DrawImage(ca.Img, ca.ImgOpts)

	// === DEBUG GRID ===
	if debug {
		// ca.DebugGrid.Draw(target, func(x, y int) Color {
		// 	if ca.activeFlags[y*ca.GridWidth+x] {
		// 		return ColorGreen
		// 	}
		// 	return ColorInactiveCell
		// })
	}
}

// PxId returns the pixel-id of the pixel at x, y coordinates
func (ca *CellAutomata) PxId(x, y int) (int, bool) {
	if x < 0 || x >= ca.Width || y < 0 || y >= ca.Height {
		return -1, false
	}
	return y*ca.Width + x, true
}

// PxPos returns the x, y pixel-coordinates of the pixel at pixel-id
func (ca *CellAutomata) PxPos(pxId int) (x, y int) {
	return pxId % ca.Width, pxId / ca.Width
}

// InBounds returns true if the pixel-x, y coordinates are in the world
func (ca *CellAutomata) InBounds(x, y int) bool {
	return x >= 0 && x < ca.Width && y >= 0 && y < ca.Height
}

// InGrid returns true if the x, y grid-coordinates are in the grid
func (ca *CellAutomata) InGrid(x, y int) bool {
	return x >= 0 && x < ca.GridWidth && y >= 0 && y < ca.GridHeight
}

// InCell returns true if the x, y relative-coordinates are in the cell
func (ca *CellAutomata) InCell(x, y int) bool {
	return x >= 0 && x < ca.CellSize && y >= 0 && y < ca.CellSize
}

// CellId returns the cell-id of the cell at x, y grid-coordinates
func (ca *CellAutomata) CellId(x, y int) (int, bool) {
	if x < 0 || x >= ca.GridWidth || y < 0 || y >= ca.GridHeight {
		return -1, false
	}

	return y*ca.GridWidth + x, true
}

// CellTopLeftPxPos returns the x, y pixel-coordinates of the top-left pixel of the cell at cell-id
func (ca *CellAutomata) CellTopLeftPxPos(cellID int) (px, py int) {
	cx, cy := ca.CellPos(cellID)
	return cx * ca.CellSize, cy * ca.CellSize
}

// CellPos returns the x, y grid-coordinates of the cell at cell-id
func (ca *CellAutomata) CellPos(cellId int) (x, y int) {
	return cellId % ca.GridWidth, cellId / ca.GridWidth
}

// PxPosToCellPos returns the x, y grid-coordinates of the cell at pixel-coordinates
func (ca *CellAutomata) PxPosToCellPos(x, y int) (cx, cy int) {
	return x / ca.CellSize, y / ca.CellSize
}

// PxIdToCellId returns the cell-id of the cell at pixel-id
func (ca *CellAutomata) PxIdToCellId(pxId int) int {
	x, y := ca.PxPos(pxId)
	id, _ := ca.CellId(x/ca.CellSize, y/ca.CellSize)
	return id
}

// IdToCellPos returns the x, y grid-coordinates of the cell at pixel-id
func (ca *CellAutomata) IdToCellPos(pxId int) (x, y int) {
	return pxId % ca.GridWidth, pxId / ca.GridWidth
}

func (ca *CellAutomata) PosToPxId(pos *Vector) (int, bool) {
	return ca.PxId(int(pos.X), int(pos.Y))
}

func (ca *CellAutomata) SetPx(x, y int, col Color, mat Material) {
	id, ok := ca.PxId(x, y)
	if !ok {
		return
	}

	ca.Materials[id] = mat

	// set the 4 bytes of the color in the Pixels array
	*(*uint32)(unsafe.Pointer(&ca.Pixels[id*4])) = uint32(col)
}

func (ca *CellAutomata) SetPxById(pxId int, col Color, mat Material) {
	ca.Materials[pxId] = mat

	// set the 4 bytes of the color in the Pixels array
	*(*uint32)(unsafe.Pointer(&ca.Pixels[pxId*4])) = uint32(col)
}

func (ca *CellAutomata) SetRect(x, y, w, h int, getProps func() (Color, Material)) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			col, mat := getProps()
			xx := x + i
			yy := y + j
			ca.SetPx(xx, yy, col, mat)
			ca.ActivateNeighborsByPxPos(xx, yy)
		}
	}
}

func (ca *CellAutomata) Erase() {
	ca.Materials = make([]Material, ca.Size)
	ca.Pixels = make([]byte, ca.Size*4)
	ca.Tick = 0
	ca.Processed = make([]int, ca.Size)
	ca.ActivateAll()
}

func (ca *CellAutomata) SwapPx(idA, idB int) {
	// swap material
	ca.Materials[idA], ca.Materials[idB] = ca.Materials[idB], ca.Materials[idA]

	// swap color
	pa := (*uint32)(unsafe.Pointer(&ca.Pixels[idA<<2]))
	pb := (*uint32)(unsafe.Pointer(&ca.Pixels[idB<<2]))
	*pa, *pb = *pb, *pa

	// mark both, as processed
	ca.Processed[idA] = ca.Tick
	ca.Processed[idB] = ca.Tick
}

func (ca *CellAutomata) IsEmpty(x, y int) bool {
	pxId, ok := ca.PxId(x, y)
	if !ok {
		return false
	}
	return ca.Materials[pxId].IsEmpty()
}

func (ca *CellAutomata) IsEmptyPos(pos *Vector) bool {
	return ca.IsEmpty(int(pos.X), int(pos.Y))
}

// ActivateNeighbors activates a 3x3 cell neighborhood around a cell ID
func (ca *CellAutomata) ActivateNeighbors(cellId int) {
	x, y := ca.CellPos(cellId)
	for xx := x - 1; xx <= x+1; xx++ {
		for yy := y - 1; yy <= y+1; yy++ {

			cid, ok := ca.CellId(xx, yy)

			if ok && !ca.NextActiveCellFlags[cid] {
				ca.NextActiveCellFlags[cid] = true
				ca.NextActiveCells = append(ca.NextActiveCells, cid)
			}
		}
	}
}

func (ca *CellAutomata) ActivateNeighborsByPos(pos *Vector) {
	id, ok := ca.CellId(int(pos.X)/ca.CellSize, int(pos.Y)/ca.CellSize)
	if !ok {
		return
	}
	ca.ActivateNeighbors(id)
}

func (ca *CellAutomata) ActivateNeighborsByPxPos(x, y int) {
	cellId, ok := ca.CellId(x/ca.CellSize, y/ca.CellSize)
	if !ok {
		return
	}
	ca.ActivateNeighbors(cellId)
}

func (ca *CellAutomata) ActivateAll() {
	for i := 0; i < ca.GridSize; i++ {
		if !ca.NextActiveCellFlags[i] {
			ca.NextActiveCellFlags[i] = true
			ca.NextActiveCells = append(ca.NextActiveCells, i)
		}
	}
}

// //////////////
// Generator //
// ////////////
type GeneratorOptions struct {
	Seed int

	Density float64

	HighlightColor1 Color
	HighlightColor2 Color
	BaseColor       Color
	ShadowColor     Color

	TopClosed    bool
	BottomClosed bool
	LeftClosed   bool
	RightClosed  bool
}

func (ca *CellAutomata) Generate(opts GeneratorOptions) {
	mapA := [][]bool{}
	mapB := [][]bool{}

	w := ca.Width + 6
	h := ca.Height + 6

	for x := 0; x < w; x++ {
		mapA = append(mapA, make([]bool, h))
		mapB = append(mapB, make([]bool, h))
		for y := 0; y < h; y++ {
			if x < 3 || x >= w-3 || y < 3 || y >= h-3 {
				mapA[x][y] = (x < 3 && opts.LeftClosed) || (x >= w-3 && opts.RightClosed) || (y < 3 && opts.TopClosed) || (y >= h-3 && opts.BottomClosed)
			} else {
				mapA[x][y] = rand.Float64() < opts.Density
			}
			mapB[x][y] = mapA[x][y]
		}
	}

	for i := 0; i < 5; i++ {
		for x := 3; x < w-3; x++ {
			for y := 3; y < h-3; y++ {
				n := 0
				for dx := -3; dx <= 3; dx++ {
					for dy := -3; dy <= 3; dy++ {
						if !(dx == 0 && dy == 0) && mapA[x+dx][y+dy] {
							n++
						}
					}
				}
				mapB[x][y] = n >= 24
			}
		}

		mapA, mapB = mapB, mapA
	}

	for x := 0; x < ca.Width; x++ {
		for y := 0; y < ca.Height; y++ {
			if mapA[x+3][y+3] {
				ca.SetPx(x, y, StoneColors[rand.Intn(len(StoneColors))], MaterialStone)
			} else {
				ca.SetPx(x, y, ColorNull, MaterialEmpty)
			}
		}
	}

	ca.ActivateAll()
}

//////////////////////////
// Material Processors //
////////////////////////

// Each processor should decide what happens with one grain/pixel of material, and return true if there were any activity, so we can activate the cells accordingly

var (
	sandDirectionA = [][]int{
		{0, 1},
		{1, 1},
		{-1, 1},
	}

	sandDirectionB = [][]int{
		{0, 1},
		{-1, 1},
		{1, 1},
	}
)

func ProcessSand(ca *CellAutomata, pxId, x, y int) bool {
	dir := sandDirectionA
	if rand.Intn(2) == 0 {
		dir = sandDirectionB
	}

	for _, d := range dir {
		id, ok := ca.PxId(x+d[0], y+d[1])
		if ok && ca.Processed[id] != ca.Tick && ca.Materials[id].IsPenetrable() {
			ca.SwapPx(pxId, id)
			return true
		}
	}

	return false
}

var (
	// left -> right
	WaterDirectionA = [][]int{
		{0, 1},
		{1, 1},
		{1, 0},
	}

	// right -> left
	WaterDirectionB = [][]int{
		{0, 1},
		{-1, 1},
		{-1, 0},
	}
)

func ProcessWater(ca *CellAutomata, pxId, x, y int) bool {
	// check desired flow direction, with a bit of randomness
	left := ca.Materials[pxId].FaceLeft()
	if rand.Intn(20) == 0 {
		left = !left
	}
	dir := WaterDirectionA
	if left {
		dir = WaterDirectionB
	}

	// try to flow
	for _, d := range dir {
		id, ok := ca.PxId(x+d[0], y+d[1])
		if ok && ca.Processed[id] != ca.Tick && ca.Materials[id].IsFlowable() {
			ca.SwapPx(pxId, id)
			return true
		}
	}

	// check the other direction
	dir = WaterDirectionB
	if left {
		dir = WaterDirectionA
	}
	// if it is flowable flip desired direction
	for _, d := range dir {
		id, ok := ca.PxId(x+d[0], y+d[1])
		if ok && ca.Materials[id].IsFlowable() {
			ca.Materials[pxId] = ca.Materials[pxId].SetBool(0, !left)
			return true
		}
	}

	// if we cannot flow in any direction, do nothing
	return false
}

var (
	// left -> right
	SmokeDirectionA = [][]int{
		{0, -1},
		{1, -1},
		{1, 0},
	}

	// right -> left
	SmokeDirectionB = [][]int{
		{0, -1},
		{-1, -1},
		{-1, 0},
	}
)

func ProcessSmoke(ca *CellAutomata, pxId, x, y int) bool {
	// check desired raise direction, with a bit of randomness
	left := ca.Materials[pxId].FaceLeft()
	if rand.Intn(20) == 0 {
		left = !left
	}
	dir := SmokeDirectionA
	if left {
		dir = SmokeDirectionB
	}

	// try to raise
	for _, d := range dir {
		id, ok := ca.PxId(x+d[0], y+d[1])
		if ok && ca.Processed[id] != ca.Tick && ca.Materials[id].IsEmpty() {
			ca.SwapPx(pxId, id)
			return true
		}
	}

	// check the other direction
	dir = SmokeDirectionB
	if left {
		dir = SmokeDirectionA
	}
	// if it is empty flip desired direction
	for _, d := range dir {
		id, ok := ca.PxId(x+d[0], y+d[1])
		if ok && ca.Processed[id] != ca.Tick && ca.Materials[id].IsEmpty() {
			ca.Materials[pxId] = ca.Materials[pxId].SetBool(0, !left)
			return true
		}
	}

	// if we cannot raise in any direction, do nothing
	return false
}
