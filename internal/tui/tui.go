package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
)

func New(log *log.Logger, chip8 *chip8.Chip8) *App {
	app := new(App)
	app.log = log
	app.chip8 = chip8
	app.keyMap = map[string]uint8{
		"1": 0x1, "2": 0x2, "3": 0x3, "4": 0xC,
		"q": 0x4, "w": 0x5, "e": 0x6, "r": 0xD,
		"a": 0x7, "s": 0x8, "d": 0x9, "f": 0xE,
		"z": 0xA, "x": 0x0, "c": 0xB, "v": 0xF,
	}

	app.program = tea.NewProgram(app, tea.WithAltScreen())

	return app
}

type App struct {
	log     *log.Logger
	chip8   *chip8.Chip8
	keyMap  map[string]uint8
	program *tea.Program
}

func (app *App) Run() error {
	_, err := app.program.Run()

	return err
}

func (app *App) Init() tea.Cmd {
	return func() tea.Msg {
		return tea.SetWindowTitle("chip-8")
	}
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return app, tea.Quit
		} else if code, ok := app.keyMap[msg.String()]; ok {
			app.chip8.KeyPress(code)
		}
	}
	return app, func() tea.Msg {
		return true
	}
}

func (app *App) View() string {
	var view strings.Builder

	fb := app.chip8.Framebuffer()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if fb[x][y] == 0 {
				view.WriteString(" ")
			} else {
				view.WriteString("█")
			}
		}
		view.WriteString("\n")
	}

	return view.String()
}
