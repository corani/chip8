//go:build amd64 && (windows || linux)

package main

import (
	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/corani/chip-8/internal/gui"
)

func init() {
	availableUIs.Register("gui", func(log *log.Logger, chip8 *chip8.Chip8) App {
		return gui.New(log, chip8)
	})
}
