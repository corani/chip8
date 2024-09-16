package main

import (
	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/corani/chip-8/internal/tui"
)

func init() {
	availableUIs.Register("tui", func(log *log.Logger, chip8 *chip8.Chip8) App {
		return tui.New(log, chip8)
	})
}
