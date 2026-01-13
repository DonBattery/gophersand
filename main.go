package main

import (
	"gophersand/game"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	VERSION = "dev"
	BUILD   = "latest"
)

func main() {
	ebiten.SetWindowTitle("GopherSand")

	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(game.NewGame(VERSION, BUILD)); err != nil {
		panic(err)
	}
}
