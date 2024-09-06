package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/corani/chip-8/internal/tui"
)

func main() {
	logger := log.New(os.Stdout)
	chip8 := chip8.New()

	app := tui.New(logger, chip8)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
