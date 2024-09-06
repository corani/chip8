package chip8

import (
	"time"

	"github.com/corani/chip-8/internal/cpu"
	"github.com/corani/chip-8/internal/display"
	"github.com/corani/chip-8/internal/keyboard"
	"github.com/corani/chip-8/internal/memory"
	"github.com/corani/chip-8/internal/sound"
	"github.com/corani/chip-8/internal/timer"
)

func New() *Chip8 {
	soundTimer := timer.New()

	chip8 := &Chip8{
		memory:   memory.New(),
		display:  display.New(),
		keyboard: keyboard.New(),
		sound:    sound.New(soundTimer),
		delay:    timer.New(),
	}

	chip8.cpu = cpu.New(chip8.memory, chip8.display, chip8.keyboard, chip8.delay, soundTimer)

	return chip8
}

type Chip8 struct {
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
	c.keyboard.Tick(dt)
	c.cpu.Tick(dt)
}

func (c *Chip8) KeyPress(code uint8) {
	c.keyboard.KeyPress(code)
}

func (c *Chip8) Framebuffer() [64][32]byte {
	return c.display.Framebuffer
}
