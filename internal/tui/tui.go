package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
)

// TODO(daniel): `keyHold` needs to be tuned.
const keyHold = 30 * time.Millisecond

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
	app.keyDown = make(map[uint8]time.Duration)

	app.program = tea.NewProgram(app, tea.WithAltScreen(), tea.WithFPS(60))
	app.view.Grow(64*32 + 32)

	return app
}

type App struct {
	log     *log.Logger
	chip8   *chip8.Chip8
	keyMap  map[string]uint8
	program *tea.Program
	dt      time.Time
	view    strings.Builder
	keyDown map[uint8]time.Duration
}

func (app *App) Run() error {
	app.dt = time.Now()
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
			app.chip8.KeyDown(code)
			app.keyDown[code] = keyHold
		}
	}

	now := time.Now()
	dt := now.Sub(app.dt)
	app.dt = now

	// Bubbletea does not (currently) support key up events,
	// so we simulate them by checking if a key has been "held"
	// for a certain amount of time.
	for code, hold := range app.keyDown {
		if hold > 0 {
			app.keyDown[code] = hold - dt

			if app.keyDown[code] <= 0 {
				app.keyDown[code] = 0

				app.chip8.KeyUp(code)
			}
		}
	}

	app.chip8.Tick(dt)

	return app, func() tea.Msg {
		time.Sleep(16 * time.Millisecond)

		return true
	}
}

func (app *App) View() string {
	fb := app.chip8.Framebuffer()

	app.view.Reset()

	for y := 0; y < len(fb[0]); y++ {
		for x := 0; x < len(fb); x++ {
			if fb[x][y] == 0 {
				app.view.WriteRune(' ')
			} else {
				app.view.WriteRune('█')
			}
		}
		app.view.WriteRune('\n')
	}

	return app.view.String()
}
