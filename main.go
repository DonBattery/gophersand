package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"gophersand/game"
)

const VERSION = "latest"

const BUILD = "dev"

var (
	//go:embed assets/55.ttf
	ttf_55 []byte

	DefaultFont *text.GoTextFaceSource
)

func main() {
	// create the default font, and set up the logger
	source, err := text.NewGoTextFaceSource(bytes.NewReader(ttf_55))
	if err != nil {
		fmt.Printf("Failed to load font: %v", err)
		os.Exit(1)
	}

	DefaultFont = source

	game.InitLogger(DefaultFont)

	// ebiten.SetCursorMode(ebiten.CursorModeHidden)

	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game.NewGame(VERSION, BUILD)); err != nil {
		fmt.Printf("Failed to run game: %v", err)
	}
}
