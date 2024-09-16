//go:build wasm && js

package main

import (
	"strconv"
	"strings"
	"syscall/js"
)

const (
	FPS         int = 60
	runInterval int = 1000 / FPS
)

func main() {
	doc := js.Global().Get("document")

	// TODO(daniel): handle keyboard
	// TODO(daniel): handle sound
	state := &gameState{
		canvas:  doc.Call("getElementById", "gameCanvas"),
		console: js.Global().Get("console"),
	}

	js.Global().Set("runGame", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			if len(args) != 2 {
				panic("expected two arguments")
			}

			// TODO(daniel): there's probably a better way to pass a raw
			// blob from JS to Go WASM...
			str := args[1].Call("toString").String()

			var rom []byte

			for _, v := range strings.Split(str, ",") {
				vi, _ := strconv.Atoi(v)
				rom = append(rom, byte(vi))
			}

			state.onRun(args[0].String(), rom)

			return nil
		},
	))

	// WASM main is expected to block indefinitely.
	select {}
}
