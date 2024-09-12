package chip8

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/cpu"
	"github.com/corani/chip-8/internal/display"
	"github.com/corani/chip-8/internal/keyboard"
	"github.com/corani/chip-8/internal/memory"
	"github.com/corani/chip-8/internal/sound"
	"github.com/corani/chip-8/internal/timer"
)

func New(logger *log.Logger, romfile string, romdata []uint8) *Chip8 {
	soundTimer := timer.New()

	chip8 := &Chip8{
		logger:   logger,
		memory:   memory.New(),
		display:  display.New(logger),
		keyboard: keyboard.New(),
		sound:    sound.New(soundTimer),
		delay:    timer.New(),
	}

	chip8.memory.Load(0x000, digitSprites())
	chip8.memory.Load(0x200, romdata)

	chip8.cpu = cpu.New(logger, chip8.memory, chip8.display, chip8.keyboard, chip8.delay, soundTimer)

	return chip8
}

type Chip8 struct {
	logger   *log.Logger
	memory   *memory.Memory
	display  *display.Display
	keyboard *keyboard.Keyboard
	sound    *sound.Sound
	delay    *timer.Timer
	cpu      *cpu.CPU
}

func (c *Chip8) LoadROM(rom []uint8) {
}

func (c *Chip8) Tick(dt time.Duration) {
	c.delay.Tick(dt)
	c.sound.Tick(dt)
	c.cpu.Tick(dt)
}

func (c *Chip8) KeyDown(code uint8) {
	c.keyboard.KeyDown(code)
}

func (c *Chip8) KeyUp(code uint8) {
	c.keyboard.KeyUp(code)
}

func (c *Chip8) Framebuffer() [][]uint8 {
	return c.display.Framebuffer
}

func digitSprites() []uint8 {
	return []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
}
