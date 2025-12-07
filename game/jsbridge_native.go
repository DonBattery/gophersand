//go:build !js

package game

// initJSBridge is a no-op on native builds; the JS bridge is only
// required when running in a browser (GOOS=js, GOARCH=wasm).
func initJSBridge() {}

// SendUIEvent is a no-op on native builds.
func SendUIEvent(event string) {}
