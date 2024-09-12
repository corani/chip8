package gui

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func New(log *log.Logger, chip8 *chip8.Chip8) *App {
	keyMap := map[ebiten.Key]uint8{
		ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
		ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
		ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
	}
	return &App{
		logger: log,
		chip8:  chip8,
		pixels: make([]uint8, 64*32*4),
		time:   time.Now(),
		keyMap: keyMap,
	}
}

type App struct {
	logger *log.Logger
	chip8  *chip8.Chip8
	pixels []uint8
	time   time.Time
	keyMap map[ebiten.Key]uint8
}

func (app *App) Update() error {
	now := time.Now()
	dt := now.Sub(app.time)
	app.time = now

	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		if k, ok := app.keyMap[key]; ok {
			app.chip8.KeyDown(k)
		}
	}

	for _, key := range inpututil.AppendJustReleasedKeys(nil) {
		if k, ok := app.keyMap[key]; ok {
			app.chip8.KeyUp(k)
		}
	}

	app.chip8.Tick(dt)

	fb := app.chip8.Framebuffer()

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			idx := (y*64 + x) * 4
			if fb[x][y] == 0 {
				app.pixels[idx] = 0
				app.pixels[idx+1] = 0
				app.pixels[idx+2] = 0
				app.pixels[idx+3] = 255
			} else {
				app.pixels[idx] = 255
				app.pixels[idx+1] = 255
				app.pixels[idx+2] = 255
				app.pixels[idx+3] = 255
			}
		}
	}

	return nil
}

func (app *App) Draw(screen *ebiten.Image) {
	screen.WritePixels(app.pixels)
}

func (app *App) Layout(outsideWith, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}

func (app *App) Run() error {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("chip-8")

	return ebiten.RunGame(app)
}
