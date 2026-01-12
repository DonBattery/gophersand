package game

import (
	"math/rand"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rngRingSize = 4096            // 4KB ring buffer
	rngRingMask = rngRingSize - 1 // fast modulo for power of 2
)

// MaterialProcessor is a function that can process one type of material.
// It returns true if there is potential for activity in the next update (e.g.: Sand returns false if it cannot move in any of its potential 3 directions).
type MaterialProcessor func(
	ca *CellAutomata,
	kind MaterialKind,
	mat Material,
	cid, x, y int,
) bool

// MaterialReaction handles MaterialA to MaterialB reaction (usually there is a slightly different MaterialReaction for MaterialB to MaterialA reactions)
// It returns true if the reaction succeeded (a lot of reactions are chance based)
type MaterialReaction func(
	ca *CellAutomata,
	matA, matB Material,
	cidA, cidB int,
) bool

// TurnPhase is a helper struct for slower materials
type TurnPhase struct {
	Turn2       bool // every 2 ticks
	Turn2shift1 bool // every 2 ticks, but shifted by 1
	Turn3       bool // every 3 ticks
	Turn3shift1 bool // every 3 ticks, but shifted by 1
	Turn3shift2 bool // every 3 ticks, but shifted by 2
	Turn5       bool // every 5 ticks
	Turn5shift1 bool // every 5 ticks, but shifted by 1
	Turn5shift2 bool // every 5 ticks, but shifted by 2
	Turn5shift3 bool // every 5 ticks, but shifted by 3
	Turn5shift4 bool // every 5 ticks, but shifted by 4
}

func (tp *TurnPhase) Update(tick int) {
	tp.Turn2 = tick%2 == 0
	tp.Turn2shift1 = tick%2 == 1
	tp.Turn3 = tick%3 == 0
	tp.Turn3shift1 = tick%3 == 1
	tp.Turn3shift2 = tick%3 == 2
	tp.Turn5 = tick%5 == 0
	tp.Turn5shift1 = tick%5 == 1
	tp.Turn5shift2 = tick%5 == 2
	tp.Turn5shift3 = tick%5 == 3
	tp.Turn5shift4 = tick%5 == 4
}

// CellAutomata manages a 256x256 grid of cells (represented by pixels). Each of this cell has its own material, and is processed accordingly.
// Each cell can be processed exactly once per tick. Swapping with another cell marks both as processed.
// We divide the grid into 64 32x32 tiles, to keep track of activity in the Cell Automata and only process tiles which will be potentially active.
type CellAutomata struct {
	isRunning bool

	tick int

	tp *TurnPhase

	// RNG state - 4KB ring buffer of pre-generated random bytes, consumed bit-by-bit
	rndRing [rngRingSize]byte // pre-filled random bytes
	rndIdx  int               // current position in ring buffer
	rndBits byte              // current byte for bit extraction
	rndBitN uint8             // bits remaining in rndBits (0-8)

	// Per cell information in flat arrays (index is calculated as y * width + x)
	// 4 byte per pixel (RGBA)
	pixels []byte

	// Material kind and status data
	materials []Material

	// Each cell can be processed exactly once per tick, if the corresponding entry is set to the current tick, it means the cell is already processed.
	processed []int

	// A 64 bit long bit-field indicating which tiles will be active in the next update
	wakeTiles uint64

	processors []MaterialProcessor

	reactions []MaterialReaction

	// Ebiten image and options for rendering
	img     *ebiten.Image
	imgOpts *ebiten.DrawImageOptions
}

func NewCellAutomata() *CellAutomata {
	ca := &CellAutomata{
		isRunning: true,

		tp: &TurnPhase{},

		pixels: make([]byte, WorldSize*4),

		materials: make([]Material, WorldSize),

		processed: make([]int, WorldSize),

		processors: make([]MaterialProcessor, 16),
		reactions:  make([]MaterialReaction, 256),

		img:     ebiten.NewImage(WorldWidth, WorldHeight),
		imgOpts: &ebiten.DrawImageOptions{},
	}

	// Initialize math/rand with current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Fill the ring buffer with random bytes using math/rand
	for i := 0; i < rngRingSize; i++ {
		ca.rndRing[i] = byte(rng.Intn(256))
	}

	return ca
}

/*

   Random Number Generator Methods

*/

// rngByte returns the next random byte from the ring buffer.
func (ca *CellAutomata) rngByte() byte {
	b := ca.rndRing[ca.rndIdx]
	ca.rndIdx = (ca.rndIdx + 1) & rngRingMask
	return b
}

// rngBool returns a random boolean, consuming one bit from the current byte.
// When all 8 bits are consumed, fetches the next byte from the ring.
func (ca *CellAutomata) rngBool() bool {
	if ca.rndBitN == 0 {
		ca.rndBits = ca.rngByte()
		ca.rndBitN = 8
	}
	b := (ca.rndBits & 0x80) != 0
	ca.rndBits <<= 1
	ca.rndBitN--
	return b
}

// rngChance256 returns true with probability chance/256.
// chance: 0 => 0% ; 255 => 100% ; otherwise chance/256.
func (ca *CellAutomata) rngChance256(chance uint8) bool {
	if chance == 0 {
		return false
	}
	if chance == 255 {
		return true
	}
	return ca.rngByte() < chance
}

// rngPick3 picks 0, 1, or 2 based on thresholds.
// Returns 0 if byte < t0, 1 if byte < t1, else 2.
func (ca *CellAutomata) rngPick3(t0, t1 uint8) uint8 {
	b := ca.rngByte()
	if b < t0 {
		return 0
	}
	if b < t1 {
		return 1
	}
	return 2
}

// rngPick4 picks 0, 1, 2, or 3 based on thresholds.
// Returns 0 if byte < t0, 1 if byte < t1, 2 if byte < t2, else 3.
func (ca *CellAutomata) rngPick4(t0, t1, t2 uint8) uint8 {
	b := ca.rngByte()
	if b < t0 {
		return 0
	}
	if b < t1 {
		return 1
	}
	if b < t2 {
		return 2
	}
	return 3
}

// rng012 returns 0, 1, or 2.
func (ca *CellAutomata) rng012() uint8 {
	return ca.rngPick3(85, 170)
}

// rng0123 returns 0, 1, 2, or 3.
func (ca *CellAutomata) rng0123() uint8 {
	return ca.rngPick4(64, 128, 192)
}

/*

   Material Processor and Reaction Registration Methods

*/

func (ca *CellAutomata) RegisterMaterialProcessors(processors []struct {
	kind      MaterialKind
	processor MaterialProcessor
}) {
	for _, p := range processors {
		ca.processors[p.kind] = p.processor
	}
}

// RegisterMaterialReactions registers a reaction between two material kinds in the reaction table
func (ca *CellAutomata) RegisterMaterialReactions(reactions []struct {
	matA     MaterialKind
	matB     MaterialKind
	reaction MaterialReaction
}) {
	for _, r := range reactions {
		ca.reactions[r.matA*16+r.matB] = r.reaction
	}
}

/*

   General Helper Methods

*/

// InBounds returns true if the x, y coordinates are inside the World
func (ca *CellAutomata) InBounds(x, y int) bool {
	if uint(x) > 255 || uint(y) > 255 {
		return false
	}
	return true
}

// OnEdge returns true if the cell ID is on one of the edges of the World
func (ca *CellAutomata) OnEdge(cid int) bool {
	x := cid % WorldWidth
	y := cid / WorldWidth
	return x == 0 || x == WorldWidth-1 || y == 0 || y == WorldHeight-1
}

/*

   Material Manipulation Methods

*/

// GetMaterialAt returns the material at the given x, y coordinates, it returns MaterialEmpty if the coordinates are outside of the World
func (ca *CellAutomata) GetMaterialAt(x, y int) Material {
	if !ca.InBounds(x, y) {
		return MaterialEmpty
	}
	return ca.materials[y<<8|x]
}

// HasNeighborKind returns true if the cell has a neighbor in the given MaterialKindSet
func (ca *CellAutomata) HasNeighborKind(x, y int, set MaterialKindSet) bool {
	top := ca.GetMaterialAt(x, y-1)
	right := ca.GetMaterialAt(x+1, y)
	bottom := ca.GetMaterialAt(x, y+1)
	left := ca.GetMaterialAt(x-1, y)
	return top.IsIn(set) || right.IsIn(set) || bottom.IsIn(set) || left.IsIn(set)
}

// SetCell sets the material of a cell by its cell ID, and choses a color for it based on its Life and State
func (ca *CellAutomata) SetCell(cid int, mat Material) {
	ca.materials[cid] = mat
	*(*uint32)(unsafe.Pointer(&ca.pixels[cid*4])) = uint32(mat.GetColor())
}

// SetCellAt sets the material of a cell by its x, y coordinates, and choses a color for it based on its Life and State
func (ca *CellAutomata) SetCellAt(x, y int, mat Material) {
	if !ca.InBounds(x, y) {
		return
	}
	cid := y*WorldWidth + x
	// set the material of the cell
	ca.materials[cid] = mat
	// set the 4 bytes of the color in the pixels array
	*(*uint32)(unsafe.Pointer(&ca.pixels[cid*4])) = uint32(mat.GetColor())
}

// SwapCells swaps two cells by their cell IDs (data and color), and marks both as processed
func (ca *CellAutomata) SwapCells(cidA, cidB int) {
	mats := ca.materials
	pixs := ca.pixels
	procd := ca.processed
	tick := ca.tick

	// swap the materials
	mats[cidA], mats[cidB] = mats[cidB], mats[cidA]

	// swap the colors
	pa := (*uint32)(unsafe.Pointer(&pixs[cidA<<2]))
	pb := (*uint32)(unsafe.Pointer(&pixs[cidB<<2]))
	*pa, *pb = *pb, *pa

	// mark both as processed
	procd[cidA] = tick
	procd[cidB] = tick
}

// SetCellAsProcessed sets the material of a cell by its cell ID, and choses a color for it based on its Life and State, and marks it as processed
func (ca *CellAutomata) SetCellAsProcessed(cid int, mat Material) {
	ca.materials[cid] = mat
	*(*uint32)(unsafe.Pointer(&ca.pixels[cid*4])) = uint32(mat.GetColor())
	ca.processed[cid] = ca.tick
}

/*

   Material Creation helper Methods

*/

func (ca *CellAutomata) CreateSand(cid int, burned bool) {
	status := MaterialStatusNormal
	if burned {
		status = MaterialStatusBurned
	}
	ca.SetCellAsProcessed(cid, MaterialSand.WithLife(ca.rng0123()).WithIsPenetrable(ca.rngBool()).WithStatus(status))
}

func (ca *CellAutomata) CreateRoot(cid int) {
	ca.SetCellAsProcessed(cid, MaterialRoot.WithLife(ca.rngPick4(120, 180, 200)).WithIsPenetrable(ca.rngBool()))
}

func (ca *CellAutomata) CreatePlant(cid int) {
	var l uint8 = 0
	if ca.rngBool() {
		l = 1
	}
	ca.SetCellAsProcessed(cid, MaterialPlant.WithLife(l).WithIsPenetrable(ca.rngChance256(156)).WithCanBloom(!ca.OnEdge(cid) && ca.rngChance256(13)))
}

func (ca *CellAutomata) CreateSmoke(cid int, agingSpeed uint8) {
	ca.SetCellAsProcessed(cid, MaterialSmoke.WithLife(ca.rngPick4(20*agingSpeed, 60*agingSpeed, 80*agingSpeed)).WithFaceLeft(ca.rngBool()))
}

func (ca *CellAutomata) CreateAntHill(cid int) {
	ca.SetCellAsProcessed(cid, MaterialAntHill.WithIsPenetrable(ca.rngBool()).WithLife(ca.rng0123()))
}

/*

   Reaction Methods

*/

// GetReaction returns the reaction between two MaterialKinds from the reaction table.
// nil is returned if there is no reaction between the two MAterials
func (ca *CellAutomata) GetReaction(kindA, kindB MaterialKind) MaterialReaction {
	return ca.reactions[int(kindA)*16+int(kindB)]
}

// CanReactAt checks if it is possible for material kind to react with another material at the given position
func (ca *CellAutomata) CanReactAt(kind MaterialKind, x, y int) bool {
	if !ca.InBounds(x, y) {
		return false
	}
	return ca.reactions[int(kind)*16+int(ca.materials[y<<8|x].GetKind())] != nil
}

// TryReactionAt checks if the x, y coordinates are inside the World, and the given MaterialA is able to react with MaterialB at that position.
// It returns two booleans the first indicating if it is even possible for MaterialA to react at that position,
// and the second indicating if it reacted successfully this time
func (ca *CellAutomata) TryReactionAt(cidA int, matA Material, kindA MaterialKind, x, y int) (canReact bool, reacted bool) {
	if !ca.InBounds(x, y) {
		return false, false
	}

	cidB := y<<8 | x

	// get the material at the position
	matB := ca.materials[cidB]
	kindB := matB.GetKind()

	// get the reaction, if the materials cannot react return false
	reaction := ca.GetReaction(kindA, kindB)
	if reaction == nil {
		return false, false
	}

	// if it can react, but already processed
	if ca.processed[cidB] == ca.tick {
		return true, false
	}

	// try the reaction and return the result
	return true, reaction(ca, matA, matB, cidA, cidB)
}

/*

   Activity Methods

*/

// WakeAll wakes all tiles in the CellAutomata, so it is guaranteed that everything will be processed in the next update
func (ca *CellAutomata) WakeAll() {
	ca.wakeTiles = ^uint64(0)
}

// WakenNeighborhood calculates the grid coordinates of x, y cell coordinates,
// and wakes all tiles in the 3x3 grid-neighborhood
func (ca *CellAutomata) WakenNeighborhood(x, y int) {
	tx := x / 32
	ty := y / 32

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			xx := tx + dx
			yy := ty + dy
			if xx < 0 || xx > 7 || yy < 0 || yy > 7 {
				continue
			}
			ca.wakeTiles |= 1 << (yy*8 + xx)
		}
	}
}

/*
New World Generator
*/
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

// Generate a new world
func (ca *CellAutomata) Generate(opts GeneratorOptions) {
	// Create two 2d array of booleans
	mapA := [][]bool{}
	mapB := [][]bool{}

	// The boolen maps are 3px wider in all direction than the Cell Automata (we will use this extra space to calculate borders)
	w := 262
	h := 262

	// for both mapA and mapB create the same noise, based on density
	for x := 0; x < w; x++ {
		mapA = append(mapA, make([]bool, h))
		mapB = append(mapB, make([]bool, h))
		for y := 0; y < h; y++ {
			if x < 3 || x >= w-3 || y < 3 || y >= h-3 {
				// The borders are controlled by the 4 "closed" booleans, if a border is closed it will be considered as true "wall", if it is open it will be considered as false "empty"
				mapA[x][y] = (x < 3 && opts.LeftClosed) || (x >= w-3 && opts.RightClosed) || (y < 3 && opts.TopClosed) || (y >= h-3 && opts.BottomClosed)
			} else {
				// The rest of the map is filled with random noise
				mapA[x][y] = rand.Float64() < opts.Density
			}
			// By default the two maps are the same
			mapB[x][y] = mapA[x][y]
		}
	}

	// Iterate 5 times (with each iteration the map is smoothed)
	for i := 0; i < 5; i++ {
		// Iterate over the boolean map
		for x := 3; x < w-3; x++ {
			for y := 3; y < h-3; y++ {
				// For each cell, count the "walls" in its 7x7 neighborhood (excluding itself)
				n := 0
				for dx := -3; dx <= 3; dx++ {
					for dy := -3; dy <= 3; dy++ {
						if !(dx == 0 && dy == 0) && mapA[x+dx][y+dy] {
							n++
						}
					}
				}
				// If at least the half of the neighbors are walls, the cell is considered as a wall in the next iteration
				mapB[x][y] = n >= 24
			}
		}
		// Swap the two maps
		mapA, mapB = mapB, mapA
	}

	// After we have generated the map, color it
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			if mapA[x+3][y+3] {
				ca.SetCell(y*256+x, MaterialStone.WithLife(ca.rng0123()).WithIsPenetrable(ca.rngBool()))
			} else {
				ca.SetCell(y*256+x, MaterialEmpty)
			}
		}
	}

	// Activate all cells for rendering
	ca.WakeAll()
}

/*

   Update and Draw methods

*/

// Update processes the CellAutomata for one tick
func (ca *CellAutomata) Update() {
	if !ca.isRunning {
		// If the CA is paused, the user can still change its cells, so we need to update the texture
		ca.img.WritePixels(ca.pixels)
		return
	}

	ca.tick++

	// cache variables
	tick := ca.tick
	activeTiles := ca.wakeTiles
	nextWakeTiles := uint64(0)
	procs := ca.processors
	procd := ca.processed
	mats := ca.materials

	// Update the turn phases for slower materials
	ca.tp.Update(tick)

	// Flip a coin to decide if we process tiles like [left to right, from top to bottom] or [right to left, from bottom to top]
	tileId := 0
	iterDir := 1
	if ca.rngBool() {
		tileId = 63
		iterDir = -1
	}

	// Process the tiles in the decided order
	for tileId >= 0 && tileId < 64 {
		tid := tileId
		tileId += iterDir

		// skip sleeping tiles
		if (activeTiles & (1 << tid)) == 0 {
			continue
		}

		// calculate the grid coordinates of the tile
		gridX := tid & 7     // %8
		gridY := tid >> 3    // /8
		xStart := gridX << 5 // *32
		yStart := gridY << 5 // *32
		xEnd := xStart + 32
		yEnd := yStart + 32

		// Flip a coin to decide if we sweep this tile [from left to right] or [from right to left]
		sweepDir := 1
		xxStart := xStart
		xxEnd := xEnd
		if ca.rngBool() {
			sweepDir = -1
			xxStart = xEnd - 1 // Start at last valid index
			xxEnd = xStart - 1 // End one before first valid index
		}

		// Initialize activity flags, based on these we will know which tile to activate for the next update
		isTileActive := false
		hitE, hitSE, hitS, hitSW, hitW, hitNW, hitN, hitNE := false, false, false, false, false, false, false, false

		// Sweep the tile from bottom to top
		for y := yEnd - 1; y >= yStart; y-- {
			// calculate the address of this row
			rowAddr := y << 8
			// Sweep the row in the decided direction
			for x := xxStart; x != xxEnd; x += sweepDir {
				cid := rowAddr + x
				// skip already processed cells
				if procd[cid] == tick {
					continue
				}

				// get the Material, its Kind and Processor. If there is no processor for this Material, skip it (Empty or Stone)
				mat := mats[cid]
				kind := mat.GetKind()
				processor := procs[kind]
				if processor == nil {
					continue
				}

				// process the Material, if activity is detected, flip the activity flags accordingly
				if processor(ca, kind, mat, cid, x, y) {
					isTileActive = true
					if x == xStart {
						hitW = true
						if y == yStart {
							hitNW = true
						}
						if y == yEnd-1 {
							hitSW = true
						}
					}
					if x == xEnd-1 {
						hitE = true
						if y == yStart {
							hitNE = true
						}
						if y == yEnd-1 {
							hitSE = true
						}
					}
					if y == yStart {
						hitN = true
					}
					if y == yEnd-1 {
						hitS = true
					}
				}
			}
		}

		// if any potential activity is detected in the current tile, wake it up for the next update
		// check if the activity is detected on the edges, and mark neighboring tiles (3x3 neighborhood), as wake accordingly
		if isTileActive {
			nextWakeTiles |= 1 << tid
			if hitW && gridX > 0 {
				nextWakeTiles |= 1 << (tid - 1)
				if hitNW && gridY > 0 {
					nextWakeTiles |= 1 << (tid - 9)
				}
				if hitSW && gridY < 7 {
					nextWakeTiles |= 1 << (tid + 7)
				}
			}
			if hitE && gridX < 7 {
				nextWakeTiles |= 1 << (tid + 1)
				if hitNE && gridY > 0 {
					nextWakeTiles |= 1 << (tid - 7)
				}
				if hitSE && gridY < 7 {
					nextWakeTiles |= 1 << (tid + 9)
				}
			}
			if hitN && gridY > 0 {
				nextWakeTiles |= 1 << (tid - 8)
			}
			if hitS && gridY < 7 {
				nextWakeTiles |= 1 << (tid + 8)
			}
		}
	}

	// The tiles we have detected to be potentially active will be processed in the next update
	ca.wakeTiles = nextWakeTiles

	// Update the image with the colors of the materials
	ca.img.WritePixels(ca.pixels)
}

func (ca *CellAutomata) Draw(target *ebiten.Image) {
	target.DrawImage(ca.img, ca.imgOpts)
}
