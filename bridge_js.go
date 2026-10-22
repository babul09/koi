//go:build js && wasm

package main

import "syscall/js"

var configCallback js.Func

func installBrowserBridge(g *Game) {
	api := js.Global().Get("Object").New()
	configCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			if c, err := decodeMotionConfig(args[0].String()); err == nil {
				g.Config = c
			}
		}
		return nil
	})
	api.Set("setConfig", configCallback)
	js.Global().Set("motionDNA", api)
}
