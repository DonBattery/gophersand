package game

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 256
	ScreenHeight = 256

	WorldWidth  = 256
	WorldHeight = 256
	WorldSize   = WorldWidth * WorldHeight

	CellSize   = 32
	GridWidth  = 8
	GridHeight = 8
	GridSize   = GridWidth * GridHeight

	DefaultRndSeed uint32 = 296548600

	RotateCW  = 0
	RotateCCW = 1
)

var (
	// The names of the materials, used for debugging
	MaterialKindNames = []string{
		"Empty",
		"Stone",
		"Sand",
		"Water",
		"Seed",
		"Ant",
		"Wasp",
		"Acid",
		"Fire",
		"Ice",
		"Smoke",
		"Steam",
		"Root",
		"Plant",
		"Flower",
		"AntHill",
	}

	// The names of the material statuses, used for debugging
	MaterialStatusNames = []string{
		"Normal",
		"Burned",
		"Acidic",
		"Frozen",
	}
)

type Game struct {
	Version string
	Build   string

	FullScreen bool
	DebugInfo  bool

	SiteEvents chan string

	BrushSize     int
	BrushMaterial Material
	BrushMode     uint8

	brushes []BrushActions

	ca *CellAutomata
}

func NewGame(version, build string) *Game {
	ca := NewCellAutomata()

	ca.RegisterMaterialProcessors([]struct {
		kind      MaterialKind
		processor MaterialProcessor
	}{
		{kind: MaterialKindSand, processor: ProcessSand},
		{kind: MaterialKindWater, processor: ProcessWater},
		{kind: MaterialKindSeed, processor: ProcessSeed},
		{kind: MaterialKindAntHill, processor: ProcessAntHill},
		{kind: MaterialKindAcid, processor: ProcessAcid},
		{kind: MaterialKindFire, processor: ProcessFire},
		{kind: MaterialKindIce, processor: ProcessIce},
		{kind: MaterialKindSmoke, processor: ProcessSmoke},
		{kind: MaterialKindSteam, processor: ProcessSteam},
		{kind: MaterialKindRoot, processor: ProcessRoot},
		{kind: MaterialKindPlant, processor: ProcessPlant},
		{kind: MaterialKindFlower, processor: ProcessFlower},
		{kind: MaterialKindAnt, processor: ProcessAnt},
		{kind: MaterialKindWasp, processor: ProcessWasp},
	})

	ca.RegisterMaterialReactions([]struct {
		matA     MaterialKind
		matB     MaterialKind
		reaction MaterialReaction
	}{
		// Sand
		{matA: MaterialKindSand, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindSand, matB: MaterialKindSteam, reaction: SwapReaction(240)},
		{matA: MaterialKindSand, matB: MaterialKindSmoke, reaction: SwapReaction(230)},
		{matA: MaterialKindSand, matB: MaterialKindWater, reaction: SwapReaction(180)},
		{matA: MaterialKindSand, matB: MaterialKindAcid, reaction: ReactionSandToAcid},
		{matA: MaterialKindSand, matB: MaterialKindFire, reaction: ReactionSandToFire},
		{matA: MaterialKindSand, matB: MaterialKindIce, reaction: ReactionSandToIce},

		// Water
		{matA: MaterialKindWater, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindWater, matB: MaterialKindSteam, reaction: SwapReaction(220)},
		{matA: MaterialKindWater, matB: MaterialKindSmoke, reaction: SwapReaction(200)},
		{matA: MaterialKindWater, matB: MaterialKindWasp, reaction: SwapReaction(200)},
		{matA: MaterialKindWater, matB: MaterialKindAcid, reaction: ReactionAcidToWater},
		{matA: MaterialKindWater, matB: MaterialKindFire, reaction: ReactionWaterToFire},
		{matA: MaterialKindWater, matB: MaterialKindIce, reaction: ReactionWaterToIce},
		{matA: MaterialKindWater, matB: MaterialKindAnt, reaction: SwapReaction(32)},
		{matA: MaterialKindWater, matB: MaterialKindAntHill, reaction: ReactionWaterToAntHill},
		{matA: MaterialKindWater, matB: MaterialKindStone, reaction: ReactionWaterToStone},

		// Seed
		{matA: MaterialKindSeed, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindSeed, matB: MaterialKindSteam, reaction: SwapReaction(220)},
		{matA: MaterialKindSeed, matB: MaterialKindSmoke, reaction: SwapReaction(200)},
		{matA: MaterialKindSeed, matB: MaterialKindWater, reaction: SwapReaction(40)},
		{matA: MaterialKindSeed, matB: MaterialKindAcid, reaction: ReactionSeedToAcid},
		{matA: MaterialKindSeed, matB: MaterialKindIce, reaction: ReactionSeedToIce},

		// Acid
		{matA: MaterialKindAcid, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindAcid, matB: MaterialKindSteam, reaction: SwapReaction(210)},
		{matA: MaterialKindAcid, matB: MaterialKindSmoke, reaction: SwapReaction(190)},
		{matA: MaterialKindAcid, matB: MaterialKindSand, reaction: ReactionAcidToSand},
		{matA: MaterialKindAcid, matB: MaterialKindWater, reaction: ReactionAcidToWater},
		{matA: MaterialKindAcid, matB: MaterialKindStone, reaction: ReactionAcidToStone},
		{matA: MaterialKindAcid, matB: MaterialKindSeed, reaction: ReactionAcidToSeed},
		{matA: MaterialKindAcid, matB: MaterialKindAnt, reaction: ReactionAcidToAnt},
		{matA: MaterialKindAcid, matB: MaterialKindAntHill, reaction: ReactionAcidToAntHill},
		{matA: MaterialKindAcid, matB: MaterialKindWasp, reaction: ReactionAcidToWasp},
		{matA: MaterialKindAcid, matB: MaterialKindFire, reaction: ReactionAcidToFire},
		{matA: MaterialKindAcid, matB: MaterialKindRoot, reaction: ReactionAcidToRoot},
		{matA: MaterialKindAcid, matB: MaterialKindPlant, reaction: ReactionAcidToPlant},
		{matA: MaterialKindAcid, matB: MaterialKindFlower, reaction: ReactionAcidToFlower},
		{matA: MaterialKindAcid, matB: MaterialKindIce, reaction: ReactionAcidToIce},

		// Fire
		{matA: MaterialKindFire, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindFire, matB: MaterialKindSteam, reaction: SwapReaction(100)},      // Steam rises above Fire more easily
		{matA: MaterialKindFire, matB: MaterialKindSmoke, reaction: SwapReaction(120)},      // Smoke rises above Fire more easily
		{matA: MaterialKindFire, matB: MaterialKindSand, reaction: FireBurnReaction(0, 20)}, // Burns Sand, no transformation, small chance to turn to smoke
		{matA: MaterialKindFire, matB: MaterialKindStone, reaction: FireBurnReaction(0, 0)}, // Burns Stone, no transformation
		{matA: MaterialKindFire, matB: MaterialKindWater, reaction: ReactionFireToWater},
		{matA: MaterialKindFire, matB: MaterialKindSeed, reaction: FireBurnReaction(10, 10)}, // Small chance to ignite, small chance to turn to smoke
		{matA: MaterialKindFire, matB: MaterialKindAnt, reaction: ReactionFireToAnt},
		{matA: MaterialKindFire, matB: MaterialKindAntHill, reaction: ReactionFireToAntHill},
		{matA: MaterialKindFire, matB: MaterialKindWasp, reaction: ReactionFireToWasp},
		{matA: MaterialKindFire, matB: MaterialKindAcid, reaction: ReactionFireToAcid},
		{matA: MaterialKindFire, matB: MaterialKindRoot, reaction: FireBurnReaction(200, 20)}, // High chance to ignite, small chance to turn to smoke
		{matA: MaterialKindFire, matB: MaterialKindPlant, reaction: ReactionFireToPlant},
		{matA: MaterialKindFire, matB: MaterialKindFlower, reaction: FireBurnReaction(180, 20)}, // High chance to ignite, small chance to turn to smoke
		{matA: MaterialKindFire, matB: MaterialKindIce, reaction: ReactionFireToIce},

		// Ice
		{matA: MaterialKindIce, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindIce, matB: MaterialKindSteam, reaction: ReactionIceToSteam},
		{matA: MaterialKindIce, matB: MaterialKindSmoke, reaction: SwapReaction(200)},
		{matA: MaterialKindIce, matB: MaterialKindSand, reaction: ReactionIceToSand},
		{matA: MaterialKindIce, matB: MaterialKindWater, reaction: ReactionIceToWater},
		{matA: MaterialKindIce, matB: MaterialKindSeed, reaction: ReactionIceToSeed},
		{matA: MaterialKindIce, matB: MaterialKindRoot, reaction: ReactionIceToRoot},
		{matA: MaterialKindIce, matB: MaterialKindPlant, reaction: ReactionIceToPlant},
		{matA: MaterialKindIce, matB: MaterialKindFlower, reaction: ReactionIceToFlower},
		{matA: MaterialKindIce, matB: MaterialKindWasp, reaction: ReactionIceToWasp},
		{matA: MaterialKindIce, matB: MaterialKindAcid, reaction: ReactionIceToAcid},
		{matA: MaterialKindIce, matB: MaterialKindFire, reaction: ReactionIceToFire},

		// Smoke
		{matA: MaterialKindSmoke, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindSmoke, matB: MaterialKindSteam, reaction: SwapReaction(30)},
		{matA: MaterialKindSmoke, matB: MaterialKindFire, reaction: SwapReaction(180)}, // Smoke easily rises above Fire

		// Steam
		{matA: MaterialKindSteam, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindSteam, matB: MaterialKindSmoke, reaction: SwapReaction(220)},
		{matA: MaterialKindSteam, matB: MaterialKindFire, reaction: SwapReaction(200)}, // Steam easily rises above Fire

		// Root
		{matA: MaterialKindRoot, matB: MaterialKindSeed, reaction: ReactionRootToSeed},
		{matA: MaterialKindRoot, matB: MaterialKindWater, reaction: ReactionRootToWater},
		{matA: MaterialKindRoot, matB: MaterialKindSand, reaction: ReactionRootToSand},
		{matA: MaterialKindRoot, matB: MaterialKindStone, reaction: ReactionRootToStone},
		{matA: MaterialKindRoot, matB: MaterialKindRoot, reaction: ReactionRootToRoot},
		{matA: MaterialKindRoot, matB: MaterialKindPlant, reaction: ReactionRootToPlant},
		{matA: MaterialKindRoot, matB: MaterialKindIce, reaction: ReactionRootToIce},
		{matA: MaterialKindRoot, matB: MaterialKindEmpty, reaction: RootGrowthReaction(32)},
		{matA: MaterialKindRoot, matB: MaterialKindSteam, reaction: RootGrowthReaction(16)},
		{matA: MaterialKindRoot, matB: MaterialKindSmoke, reaction: RootGrowthReaction(16)},
		{matA: MaterialKindRoot, matB: MaterialKindAntHill, reaction: RootGrowthReaction(8)},

		// Plant
		{matA: MaterialKindPlant, matB: MaterialKindSeed, reaction: ReactionPlantToSeed},
		{matA: MaterialKindPlant, matB: MaterialKindRoot, reaction: ReactionPlantToRoot},
		{matA: MaterialKindPlant, matB: MaterialKindPlant, reaction: ReactionPlantToPlant},
		{matA: MaterialKindPlant, matB: MaterialKindEmpty, reaction: PlantGrowthReaction(32)},
		{matA: MaterialKindPlant, matB: MaterialKindSteam, reaction: PlantGrowthReaction(16)},
		{matA: MaterialKindPlant, matB: MaterialKindSmoke, reaction: PlantGrowthReaction(16)},
		{matA: MaterialKindPlant, matB: MaterialKindWater, reaction: ReactionPlantToWater},
		{matA: MaterialKindPlant, matB: MaterialKindIce, reaction: ReactionPlantToIce},

		// Flower (does not moves, so no reactions)

		// Ant
		{matA: MaterialKindAnt, matB: MaterialKindEmpty, reaction: SwapReaction(220)},
		{matA: MaterialKindAnt, matB: MaterialKindAntHill, reaction: AlwaysSwap},
		{matA: MaterialKindAnt, matB: MaterialKindSteam, reaction: SwapReaction(200)},
		{matA: MaterialKindAnt, matB: MaterialKindSmoke, reaction: SwapReaction(200)},
		{matA: MaterialKindAnt, matB: MaterialKindWater, reaction: SwapReaction(48)},
		{matA: MaterialKindAnt, matB: MaterialKindSand, reaction: ReactionAntToSand},
		{matA: MaterialKindAnt, matB: MaterialKindStone, reaction: ReactionAntToStone},
		{matA: MaterialKindAnt, matB: MaterialKindSeed, reaction: AntEatReaction(8)},
		{matA: MaterialKindAnt, matB: MaterialKindRoot, reaction: AntEatReaction(16)},
		{matA: MaterialKindAnt, matB: MaterialKindPlant, reaction: AntEatReaction(64)},
		{matA: MaterialKindAnt, matB: MaterialKindFlower, reaction: AntEatReaction(96)},
		{matA: MaterialKindAnt, matB: MaterialKindAcid, reaction: ReactionAntToAcid},
		{matA: MaterialKindAnt, matB: MaterialKindFire, reaction: ReactionAntToFire},
		{matA: MaterialKindAnt, matB: MaterialKindWasp, reaction: ReactionAntToWasp},

		// AntHill
		{matA: MaterialKindAntHill, matB: MaterialKindEmpty, reaction: AlwaysSwap},

		// Wasp
		{matA: MaterialKindWasp, matB: MaterialKindEmpty, reaction: AlwaysSwap},
		{matA: MaterialKindWasp, matB: MaterialKindAntHill, reaction: SwapReaction(32)},
		{matA: MaterialKindWasp, matB: MaterialKindWater, reaction: ReactionWaspToWater},
		{matA: MaterialKindWasp, matB: MaterialKindSteam, reaction: ReactionWaspToSteam},
		{matA: MaterialKindWasp, matB: MaterialKindSmoke, reaction: ReactionWaspToSmoke},
		{matA: MaterialKindWasp, matB: MaterialKindAnt, reaction: ReactionWaspToAnt},
		{matA: MaterialKindWasp, matB: MaterialKindAcid, reaction: ReactionWaspToAcid},
		{matA: MaterialKindWasp, matB: MaterialKindFire, reaction: ReactionWaspToFire},
		{matA: MaterialKindWasp, matB: MaterialKindIce, reaction: ReactionWaspToIce},
	})

	brushes := make([]BrushActions, 16)

	brushes[MaterialKindEmpty] = BrushActions{FirstAction: brushEmpty}
	brushes[MaterialKindStone] = BrushActions{FirstAction: brushStonePass1, SecondAction: brushStonePass2}
	brushes[MaterialKindSand] = BrushActions{FirstAction: brushSand}
	brushes[MaterialKindWater] = BrushActions{FirstAction: brushWater}
	brushes[MaterialKindSeed] = BrushActions{FirstAction: brushSeed}
	brushes[MaterialKindAntHill] = BrushActions{FirstAction: brushAntHill}
	brushes[MaterialKindAcid] = BrushActions{FirstAction: brushAcid}
	brushes[MaterialKindFire] = BrushActions{FirstAction: brushFire}
	brushes[MaterialKindIce] = BrushActions{FirstAction: brushIce}
	brushes[MaterialKindSmoke] = BrushActions{FirstAction: brushSmoke}
	brushes[MaterialKindSteam] = BrushActions{FirstAction: brushSteam}
	brushes[MaterialKindRoot] = BrushActions{FirstAction: brushRoot}
	brushes[MaterialKindPlant] = BrushActions{FirstAction: brushPlant}
	brushes[MaterialKindFlower] = BrushActions{FirstAction: brushFlower}
	brushes[MaterialKindAnt] = BrushActions{FirstAction: brushAnt}
	brushes[MaterialKindWasp] = BrushActions{FirstAction: brushWasp}

	g := &Game{
		Version: version,
		Build:   build,

		FullScreen: false,
		DebugInfo:  false,

		SiteEvents: make(chan string, 128),

		BrushSize:     14,
		BrushMaterial: MaterialSand,

		brushes: brushes,

		ca: ca,
	}

	g.SetupJSBridge()

	return g
}

func (g *Game) HandleSiteEvent(event string) {
	switch event {
	// brush select
	case "brush_select:empty":
		g.BrushMaterial = MaterialEmpty
	case "brush_select:stone":
		g.BrushMaterial = MaterialStone
	case "brush_select:sand":
		g.BrushMaterial = MaterialSand
	case "brush_select:water":
		g.BrushMaterial = MaterialWater
	case "brush_select:anthill":
		g.BrushMaterial = MaterialAntHill
	// Back-compat with older UI/event names.
	case "brush_select:egg":
		g.BrushMaterial = MaterialAntHill
	case "brush_select:wasp":
		g.BrushMaterial = MaterialWasp
	case "brush_select:root":
		// Root is internal-only; ignore selection.
		return
	case "brush_select:plant":
		// Plant brush places Seeds; Seeds become Root when touching both Water and Sand.
		g.BrushMaterial = MaterialSeed
	case "brush_select:seed":
		g.BrushMaterial = MaterialSeed
	case "brush_select:ant":
		g.BrushMaterial = MaterialAnt
	case "brush_select:acid":
		g.BrushMaterial = MaterialAcid
	case "brush_select:fire":
		g.BrushMaterial = MaterialFire
	case "brush_select:ice":
		g.BrushMaterial = MaterialIce

	// brush size
	case "brush_size:8":
		g.BrushSize = 8
	case "brush_size:14":
		g.BrushSize = 14
	case "brush_size:20":
		g.BrushSize = 20
	case "brush_size:26":
		g.BrushSize = 26
	case "brush_size:32":
		g.BrushSize = 32

	// world events
	case "world:stop":
		g.ca.isRunning = false
	case "world:start":
		g.ca.isRunning = true
	case "world:erase":
		g.EraseWorld()
	case "world:gen":
		g.ca.Generate(GeneratorOptions{
			Density: 0.485,
		})
	case "world:rotate_cw":
		g.RotateWorld(RotateCW)
	case "world:rotate_ccw":
		g.RotateWorld(RotateCCW)
	case "world:debug:on":
		g.DebugInfo = true
	case "world:debug:off":
		g.DebugInfo = false
	}
}

// ApplyBrush paints a circle centered at (x, y), with a given Material and Size (diameter).
// Based on the Material's Kind, it can apply two subsequent actions on the pixels inside the circle.
func (g *Game) ApplyBrush(mat Material, x, y, size int) {
	r := size / 2

	coords := make([]struct{ x, y int }, 0, size*size)

	// single pixel case
	if r <= 0 {
		if g.ca.InBounds(x, y) {
			coords = append(coords, struct{ x, y int }{x, y})
		}
		// circle case
	} else {
		r2 := r * 2
		rr := r * r
		for i := 0; i < r2; i++ {
			for j := 0; j < r2; j++ {
				dx, dy := i-r, j-r
				// skip pixels outside the circle
				if dx*dx+dy*dy >= rr {
					continue
				}
				xx := x + i - r
				yy := y + j - r
				// only add pixels that are inside the world
				if g.ca.InBounds(xx, yy) {
					coords = append(coords, struct{ x, y int }{xx, yy})
				}
			}
		}
	}

	actions := g.brushes[mat.GetKind()]

	// Pass 1
	for _, c := range coords {
		g.ca.SetCellAt(c.x, c.y, actions.FirstAction(g.ca, c.x, c.y))
	}

	// Pass 2
	if actions.SecondAction != nil {
		for _, c := range coords {
			g.ca.SetCellAt(c.x, c.y, actions.SecondAction(g.ca, c.x, c.y))
		}
	}

	g.ca.WakenNeighborhood(x, y)
}

// RotateWorld rotates the world clockwise or counterclockwise by 90 degrees
func (g *Game) RotateWorld(dir int) {
	// Create temp buffers
	newMaterials := make([]Material, WorldSize)
	newPixels := make([]byte, WorldSize*4)

	// Map:
	//   CW:  (x,y) -> (255-y, x)
	//   CCW: (x,y) -> (y, 255-x)
	for y := 0; y < 256; y++ {
		rowOld := y * 256
		for x := 0; x < 256; x++ {
			oldCid := rowOld + x

			var newX, newY int

			switch dir {
			case RotateCW:
				newX = 255 - y
				newY = x
			case RotateCCW:
				newX = y
				newY = 255 - x
			}

			newCid := newY*256 + newX

			newMaterials[newCid] = g.ca.materials[oldCid]
			copy(newPixels[newCid*4:newCid*4+4], g.ca.pixels[oldCid*4:oldCid*4+4])
		}
	}

	// Copy back
	copy(g.ca.materials, newMaterials)
	copy(g.ca.pixels, newPixels)

	g.ca.WakeAll()
}

// EraseWorld turns every non-Empty cell into Fire
func (g *Game) EraseWorld() {
	for x := 0; x < WorldWidth; x++ {
		for y := 0; y < WorldHeight; y++ {
			cid := y<<8 | x
			if g.ca.materials[cid].GetKind() != MaterialKindEmpty {
				g.ca.materials[cid] = MaterialFire.WithLife(3).WithStatus(uint8(rand.Intn(4)))
				g.ca.SetCellAsProcessed(cid, g.ca.materials[cid])
			}
		}
	}
	g.ca.WakeAll()
}

// materialInfo returns a debug string describing the material under the first cursor.
// Format: Material name, Life, Status.
// TODO: make a magnifier tool for this
func (g *Game) MaterialInfo() string {
	info := ""

	if NumberOfCursors <= 0 {
		return "Mat: (no cursor)"
	}

	x := Cursors[0].PosX
	y := Cursors[0].PosY
	if uint(x) > 255 || uint(y) > 255 {
		return "Mat: (out of bounds)"
	}

	mat := g.ca.materials[y<<8|x]
	info = fmt.Sprintf(
		"Mat: %s  L:%d  S:%s",
		MaterialKindNames[mat.GetKind()],
		mat.GetLife(),
		MaterialStatusNames[mat.GetStatus()],
	)

	switch mat.GetKind() {
	// add flow direction
	case MaterialKindWater:
		dir := "FD: right"
		if mat.GetFaceLeft() {
			dir = "FD: left"
		}
		info += fmt.Sprintf("  %s", dir)

	// add can bloom and is penetrable
	case MaterialKindPlant:
		info += fmt.Sprintf(" CB:%t P:%t", mat.GetCanBloom(), mat.GetIsPenetrable())

	// add is penetrable
	case MaterialKindSand:
		info += fmt.Sprintf(" P:%t", mat.GetIsPenetrable())
	}

	return info
}

func (g *Game) BrushInfo() string {
	return fmt.Sprintf("Brush: %s  Size: %d", MaterialKindNames[g.BrushMaterial.GetKind()], g.BrushSize)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	// first drain the JS event queue
	for {
		select {
		case event := <-g.SiteEvents:
			g.HandleSiteEvent(event)
		default:
			goto EventsDone
		}
	}

EventsDone:

	UpdateInputs()

	if KeyEsc.Pressed {
		if g.HandleEsc() {
			return ebiten.Termination
		}
	}

	if KeyD.Pressed {
		g.DebugInfo = !g.DebugInfo
		mode := "off"
		if g.DebugInfo {
			mode = "on"
		}
		g.SendToSite(fmt.Sprintf("world:debug:%s", mode))
	}

	if KeyE.Pressed {
		g.EraseWorld()
	}

	if KeyF.Pressed {
		g.SwitchFullscreen()
	}

	if KeyG.Pressed {
		g.ca.Generate(GeneratorOptions{
			Density: 0.485,
		})
	}

	if KeyP.Pressed {
		g.ca.isRunning = !g.ca.isRunning
		mode := "stop"
		if g.ca.isRunning {
			mode = "start"
		}
		g.SendToSite(fmt.Sprintf("world:%s", mode))
	}

	if KeyN0.Pressed {
		g.BrushMaterial = MaterialEmpty
		g.SendToSite("brush_select:Empty")
	}

	if KeyN1.Pressed {
		g.BrushMaterial = MaterialStone
		g.SendToSite("brush_select:Stone")
	}

	if KeyN2.Pressed {
		g.BrushMaterial = MaterialSand
		g.SendToSite("brush_select:Sand")
	}

	if KeyN3.Pressed {
		g.BrushMaterial = MaterialWater
		g.SendToSite("brush_select:Water")
	}

	if KeyN4.Pressed {
		g.BrushMaterial = MaterialSeed
		g.SendToSite("brush_select:Seed")
	}

	if KeyN5.Pressed {
		g.BrushMaterial = MaterialFire
		g.SendToSite("brush_select:Fire")
	}

	if KeyN6.Pressed {
		g.BrushMaterial = MaterialAcid
		g.SendToSite("brush_select:Acid")
	}

	if KeyN7.Pressed {
		g.BrushMaterial = MaterialAnt
		g.SendToSite("brush_select:Ant")
	}

	if KeyN8.Pressed {
		g.BrushMaterial = MaterialWasp
		g.SendToSite("brush_select:Wasp")
	}

	if KeyN9.Pressed {
		g.BrushMaterial = MaterialIce
		g.SendToSite("brush_select:Ice")
	}

	// Mouse inputs
	if MouseWheelUp && g.BrushSize < 32 {
		g.BrushSize++
		g.SendToSite(fmt.Sprintf("brush_size:%d", g.BrushSize))
	}

	if MouseWheelDown && g.BrushSize > 1 {
		g.BrushSize--
		g.SendToSite(fmt.Sprintf("brush_size:%d", g.BrushSize))
	}

	// Cursor inputs (mouse or touch)
	for i := 0; i < NumberOfCursors; i++ {
		cursor := Cursors[i]

		if cursor.LeftDown {
			g.ApplyBrush(g.BrushMaterial, cursor.PosX, cursor.PosY, g.BrushSize)
		}

		if cursor.RightDown {
			g.ApplyBrush(MaterialEmpty, cursor.PosX, cursor.PosY, g.BrushSize)
		}
	}

	g.ca.Update()

	return nil
}

func (g *Game) Draw(target *ebiten.Image) {
	g.ca.Draw(target)

	if g.DebugInfo {

		// draw tiles, color them based on activity
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				if (g.ca.wakeTiles & (uint64(1) << uint((y<<3)+x))) != 0 {
					Rect(target, x*32, y*32, 32, 32, ColorActiveCell)
				} else {
					Rect(target, x*32, y*32, 32, 32, ColorInactiveCell)
				}
			}
		}

		// print debug info
		ebitenutil.DebugPrint(
			target,
			fmt.Sprintf(
				"FPS: %0.2f TPS: %0.2f\n%s\n%s",
				ebiten.ActualFPS(),
				ebiten.ActualTPS(),
				g.MaterialInfo(),
				g.BrushInfo(),
			),
		)
	}

	// draw cursor(s)
	for i := 0; i < NumberOfCursors; i++ {
		cursor := Cursors[i]
		Circle(target, cursor.PosX, cursor.PosY, g.BrushSize/2, ColorWhite)
		target.Set(cursor.PosX, cursor.PosY, ColorWhite)
	}
}
