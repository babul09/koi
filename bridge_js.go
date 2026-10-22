//go:build js && wasm

package main

import "syscall/js"

var configCallback js.Func
var resetCallback js.Func

func installBrowserBridge(g *Game) {
	api := js.Global().Get("Object").New()
	configCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			if c, err := decodeMotionConfig(args[0].String()); err == nil {
				if c.Mode != g.Config.Mode {
					g.ResetSimulation()
				}
				g.Config = c
			}
		}
		return nil
	})
	api.Set("setConfig", configCallback)
	resetCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		g.ResetSimulation()
		return nil
	})
	api.Set("resetSimulation", resetCallback)
	js.Global().Set("motionDNA", api)
}
