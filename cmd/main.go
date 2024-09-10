package main

import (
	"flag"
	"io"
	"os"
	"runtime/pprof"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/corani/chip-8/internal/tui"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	romfile    = flag.String("rom", "", "path to the rom file")
	logfile    = flag.String("log", "", "path to the log file")
)

func main() {
	flag.Parse()

	var out io.Writer

	if *logfile != "" {
		out, err := os.Create(*logfile)
		if err != nil {
			panic(err)
		}
		defer out.Close()
	} else {
		out = io.Discard
	}

	logger := log.New(io.MultiWriter(out, os.Stderr))
	logger.SetReportTimestamp(true)

	if romfile == nil {
		logger.Errorf("Usage: %s -rom <rom-file>", os.Args[0])
		os.Exit(1)
	}

	bs, err := os.ReadFile(*romfile)
	if err != nil {
		logger.Errorf("failed to load rom: %v", err)
		os.Exit(1)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			logger.Errorf("failed to create cpu profile: %v", err)
			os.Exit(1)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Errorf("failed to start cpu profile: %v", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
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
