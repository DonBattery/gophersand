//go:build js && wasm

package game

import (
	"fmt"
	"syscall/js"
)

func ConsoleError(v ...interface{}) {
	js.Global().Get("console").Call("error", v...)
}

func ConsoleWarn(v ...interface{}) {
	js.Global().Get("console").Call("warn", v...)
}

func ConsoleDebug(v ...interface{}) {
	js.Global().Get("console").Call("debug", v...)
}

func ConsoleInfo(v ...interface{}) {
	js.Global().Get("console").Call("info", v...)
}

func ConsoleLog(v ...interface{}) {
	js.Global().Get("console").Call("log", v...)
}

func ConsoleErrorf(format string, v ...interface{}) {
	js.Global().Get("console").Call("error", fmt.Sprintf(format, v...))
}

func ConsoleWarnf(format string, v ...interface{}) {
	js.Global().Get("console").Call("warn", fmt.Sprintf(format, v...))
}

func ConsoleDebugf(format string, v ...interface{}) {
	js.Global().Get("console").Call("debug", fmt.Sprintf(format, v...))
}

func ConsoleInfof(format string, v ...interface{}) {
	js.Global().Get("console").Call("info", fmt.Sprintf(format, v...))
}

func ConsoleLogf(format string, v ...interface{}) {
	js.Global().Get("console").Call("log", fmt.Sprintf(format, v...))
}

// SetupJSBridge sets up the Webpage to Golang bridge, by injecting the "sendToGame" function into the global JS scope.
// Call this function with window.sendToGame("event-string");
// Internally these event strings will be buffered into the SiteEvents channel and drained upon game updates.
func (g *Game) SetupJSBridge() {
	js.Global().Set("sendToGame", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}
		event := args[0].String()

		select {
		case g.SiteEvents <- event:
		default:
			ConsoleErrorf("EventChannel full! Dropping event: %s\n", event)
		}

		return nil
	}))

	ConsoleInfo(fmt.Sprintf("GopherSand ver: %s build: %s JS bridge Initialized", g.Version, g.Build))
}

// SendToSite sends an event string from Golang to the parent Site
// Listen to these messages with:
//
//	window.addEventListener('message', (event) => {
//	    const data = event.data;
//	    if (!data || data.type !== 'game-event' || typeof data.payload !== 'string') {
//	        return;
//	    }
//		    const payload = data.payload;
//			console.log('GAME_EVENT:', payload);
//		});
func (g *Game) SendToSite(event string) {
	parent := js.Global().Get("parent")
	if !parent.Truthy() {
		return
	}

	msg := js.Global().Get("Object").New()
	msg.Set("type", "game-event")
	msg.Set("payload", event)

	parent.Call("postMessage", msg, "*")
}

func (g *Game) SwitchFullscreen() {
	g.SendToSite("game:fullscreen")
}

func (g *Game) HandleEsc() bool {
	return false
}
