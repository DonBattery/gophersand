// Package game provides the top level Game object
package game

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WorldWidth           = 256
	WorldHeight          = 256
	CellAutomataCellSize = 32
)

var (
	Version string
	Build   string

	CA *CellAutomata

	FullScreen = false
	DebugInfo  = false
	DebugCA    = false
	DebugPS    = false
	DebugLines = false

	BrushSize     int
	BrushMaterial Material

	EventChannel = make(chan string, 128)

	IsRunning = false
)

type Game struct{}

func NewGame(version, build string) *Game {
	Version = version
	Build = build

	BrushSize = 1

	BrushMaterial = MaterialSand

	CA = NewCellAutomata(WorldWidth, WorldHeight, CellAutomataCellSize)

	CA.Generate(GeneratorOptions{
		Density: 0.485,
	})

	// Initialize the JS <-> Go bridge (no-op on non-web builds).
	initJSBridge()

	IsRunning = false

	return &Game{}
}

func (g *Game) Update() error {
	for {
		select {
		case event := <-EventChannel:
			g.HandleJSEvent(event)
		default:
			goto EventsDone
		}
	}

EventsDone:

	UpdateInputs()

	// if KeyEsc.Pressed {
	// 	return ebiten.Termination
	// }

	if KeyF1.Pressed {
		DebugInfo = !DebugInfo
	}

	if KeyF2.Pressed {
		DebugCA = !DebugCA
	}

	if KeyF3.Pressed {
		DebugPS = !DebugPS
	}

	// if KeyF12.Pressed {
	// 	FullScreen = !FullScreen
	// 	ebiten.SetFullscreen(FullScreen)
	// }

	if KeyN1.Pressed {
		BrushMaterial = MaterialStone
		SendUIEvent("brush_select:stone")
	}

	if KeyN2.Pressed {
		BrushMaterial = MaterialSand
		SendUIEvent("brush_select:sand")
	}

	if KeyN3.Pressed {
		BrushMaterial = MaterialWater
		SendUIEvent("brush_select:water")
	}

	if KeyN4.Pressed {
		BrushMaterial = MaterialSmoke
		SendUIEvent("brush_select:smoke")
	}

	if KeyN5.Pressed {
		CA.Generate(GeneratorOptions{
			Density: 0.485,
		})
	}

	if KeyL.Pressed {
		DebugLines = !DebugLines
	}

	if MouseWheelUp && BrushSize < 64 {
		BrushSize++
	}

	if MouseWheelDown && BrushSize > 1 {
		BrushSize--
	}

	for i := 0; i < NumberOfCursors; i++ {
		cursor := Cursors[i]

		if cursor.LeftDown {
			CA.SetRect(cursor.PosX-BrushSize/2, cursor.Posy-BrushSize/2, BrushSize, BrushSize, func() (Color, Material) {
				group := MaterialColorGroups[BrushMaterial.GetType()]
				return group[rand.Intn(len(group))], BrushMaterial.RandomDirection()
			})
		}

		if cursor.RightDown {
			CA.SetRect(cursor.PosX-BrushSize/2, cursor.Posy-BrushSize/2, BrushSize, BrushSize, func() (Color, Material) { return ColorNull, MaterialEmpty })
		}
	}

	CA.Update()

	return nil
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return WorldWidth, WorldHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	CA.Draw(screen, DebugCA)

	if DebugInfo {
		PrintData(
			screen,
			0,
			0,
			[]LogData{
				{Key: "tick", Value: fmt.Sprintf("%.2f", ebiten.ActualTPS())},
				{Key: "fps", Value: fmt.Sprintf("%.2f", ebiten.ActualFPS())},
				{Key: "mouse", Value: fmt.Sprintf("%d : %d", MouseX, MouseY)},
				{Key: "brush", Value: fmt.Sprintf("%d : %s", BrushSize, BrushMaterial)},
			},
		)

		RenderLogger(screen, 0, float64(WorldHeight-40), 4)
	}

	if DebugLines {
		Line(screen, 0, MouseY, WorldWidth-1, MouseY, ColorLime)
		Line(screen, MouseX, 0, MouseX, WorldHeight-1, ColorLightLime)
	}

	// draw cursor(s)
	for i := 0; i < NumberOfCursors; i++ {
		cursor := Cursors[i]
		Rect(screen, cursor.PosX-BrushSize/2, cursor.Posy-BrushSize/2, BrushSize, BrushSize, ColorWhite)
	}
}

func (g *Game) HandleJSEvent(event string) {
	switch event {
	case "brush_select:empty":
		BrushMaterial = MaterialEmpty
	case "brush_select:stone":
		BrushMaterial = MaterialStone
	case "brush_select:sand":
		BrushMaterial = MaterialSand
	case "brush_select:water":
		BrushMaterial = MaterialWater
	case "brush_size:plus":
		BrushSize = min(BrushSize+5, 32)
	case "brush_size:minus":
		BrushSize = max(BrushSize-5, 1)
	case "world:stop":
		IsRunning = false
	case "world:start":
		// When starting the simulation from a paused state we need to ensure
		// that there are active cells in the CA. While paused, Update() keeps
		// clearing NextActiveCells, so without re-activating the grid the
		// simulation would not resume until the user draws into the world.
		IsRunning = true
		if CA != nil {
			CA.ActivateAll()
		}
	case "world:erase":
		CA.Erase()
	case "world:gen":
		CA.Generate(GeneratorOptions{
			Density: 0.485,
		})
	}
}
