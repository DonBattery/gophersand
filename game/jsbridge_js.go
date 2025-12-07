//go:build js && wasm

package game

import (
	"fmt"
	"syscall/js"
)

// initJSBridge wires up the global JS function used by the web UI to send
// events into the Go game, and also prepares facilities for sending events
// back to the surrounding web page.
func initJSBridge() {
	// JS → Go: window.sendToGame("event-string")
	js.Global().Set("sendToGame", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}
		event := args[0].String()

		select {
		case EventChannel <- event:
		default:
			LogError(fmt.Sprintf("EventChannel full, dropping event: %s", event))
		}

		return nil
	}))
}

// SendUIEvent sends an event string from the Go/wasm game up to the
// containing web page (index.html).
//
// The page listens for window "message" events with a payload in the form:
//
//	{ type: "ui-cmd", payload: "..." }
func SendUIEvent(event string) {
	parent := js.Global().Get("parent")
	if !parent.Truthy() {
		return
	}

	msg := js.Global().Get("Object").New()
	msg.Set("type", "ui-cmd")
	msg.Set("payload", event)

	parent.Call("postMessage", msg, "*")
}
