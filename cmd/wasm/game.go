//go:build wasm && js

package main

import (
	"fmt"
	"io"
	"syscall/js"
	"time"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
)

type gameState struct {
	canvas  js.Value
	console js.Value

	chip8  *chip8.Chip8
	logger *log.Logger
	time   time.Time
	grid   [64][32]uint8
}

func (state *gameState) init() {
	state.logger = log.New(io.Discard)
}

func (state *gameState) onRun(romName string, romData []byte) {
	state.log("name: %v", romName)
	state.log("size: %v", len(romData))

	state.chip8 = chip8.New(nil, romName, romData)
	state.time = time.Now()

	state.step()
}

func (state *gameState) step() {
	if state.chip8 == nil {
		return
	}

	state.log("step")

	state.update()
	state.draw()

	js.Global().Call("setTimeout", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			state.step()

			return nil
		},
	), runInterval)
}

func (state *gameState) update() {
	now := time.Now()
	dt := now.Sub(state.time)
	state.time = now

	state.chip8.Tick(dt)

	fb := state.chip8.Framebuffer()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if fb[x][y] == 0 {
				state.grid[x][y] = 0
			} else {
				state.grid[x][y] = 1
			}
		}
	}
}

func (state *gameState) draw() {
	ctx := state.canvas.Call("getContext", "2d")
	w := state.canvas.Get("width").Int()
	h := state.canvas.Get("height").Int()

	// get the cell size
	cellWidth := w / 64
	cellHeight := h / 32

	cellSize := cellHeight

	if cellWidth < cellHeight {
		cellSize = cellWidth
	}

	// center the grid
	offsetX := (w - cellWidth*64) / 2
	offsetY := (h - cellHeight*32) / 2

	// clear the canvas
	ctx.Set("fillStyle", "#f4f4f4")
	ctx.Call("fillRect", 0, 0, w, h)

	// draw the grid
	ctx.Set("fillStyle", "#1818baba")
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if state.grid[x][y] == 0 {
				continue
			}

			ctx.Call("fillRect", x*cellSize+offsetX, y*cellSize+offsetY, cellSize, cellSize)
		}
	}
}

func (state *gameState) log(format string, args ...any) {
	state.console.Call("log", fmt.Sprintf(format, args...))
}

func (state *gameState) Write(p []byte) (n int, err error) {
	state.console.Call("log", string(p))

	return len(p), nil
}
