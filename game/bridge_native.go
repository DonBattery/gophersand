//go:build !js

// the JS bridge is only required when running in a browser (GOOS=js, GOARCH=wasm).
package game

import "github.com/hajimehoshi/ebiten/v2"

// SetupJSBridge is a no-op on native builds.
func (g *Game) SetupJSBridge() {}

// SendToSite is a no-op on native builds.
func (g *Game) SendToSite(event string) {}

func (g *Game) SwitchFullscreen() {
	g.FullScreen = !g.FullScreen
	ebiten.SetFullscreen(g.FullScreen)
}

func (g *Game) HandleEsc() bool {
	return true
}
