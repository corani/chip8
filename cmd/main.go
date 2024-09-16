package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
)

type App interface {
	Run() error
}

type AppBuilder func(*log.Logger, *chip8.Chip8) App

type UIs map[string]AppBuilder

func (u UIs) Register(name string, builder AppBuilder) {
	u[name] = builder
}

func (u UIs) Available() []string {
	var names []string

	for name := range u {
		names = append(names, name)
	}

	return names
}

var availableUIs = UIs{}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	romfile := flag.String("rom", "", "path to the rom file")
	logfile := flag.String("log", "", "path to the log file")
	ui := flag.String("ui", "tui", fmt.Sprintf("user interface to use (%s)",
		strings.Join(availableUIs.Available(), ", ")))
	help := flag.Bool("help", false, "show this help message")
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

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

	chip8 := chip8.New(logger, os.Args[1], bs)

	var app App

	if builder, ok := availableUIs[*ui]; ok {
		logger.Infof("using user interface: %s", *ui)

		app = builder(logger, chip8)
	} else {
		logger.Errorf("unknown user interface: %s (supported: %s)",
			*ui, strings.Join(availableUIs.Available(), ", "))
		os.Exit(1)

	}

	// NOTE(daniel): from this point on, don't log to stderr anymore,
	// as this messes up the TUI interface.
	logger.SetOutput(out)

	if err := app.Run(); err != nil {
		logger.Errorf("run failed: %v", err)
		os.Exit(1)
	}
}
