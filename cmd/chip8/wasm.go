package main

import (
	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/corani/chip-8/internal/web"
)

func init() {
	availableUIs.Register("web", func(log *log.Logger, chip8 *chip8.Chip8) App {
		return web.New(log, chip8)
	})
}
