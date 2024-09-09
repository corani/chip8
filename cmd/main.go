package main

import (
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/corani/chip-8/internal/tui"
)

func main() {
	out, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	logger := log.New(io.MultiWriter(out, os.Stderr))
	if len(os.Args) != 2 {
		logger.Errorf("Usage: %s <rom-file>", os.Args[0])
		os.Exit(1)
	}

	bs, err := os.ReadFile(os.Args[1])
	if err != nil {
		logger.Errorf("failed to load rom: %v", err)
		os.Exit(1)
	}

	// NOTE(daniel): from this point on, don't log to stderr anymore,
	// as this messes up the TUI interface.
	logger.SetOutput(out)

	chip8 := chip8.New(logger, os.Args[1], bs)

	app := tui.New(logger, chip8)

	if err := app.Run(); err != nil {
		logger.Errorf("run failed: %v", err)
		os.Exit(1)
	}
}
