package cpu

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/display"
	"github.com/corani/chip-8/internal/keyboard"
	"github.com/corani/chip-8/internal/memory"
	"github.com/corani/chip-8/internal/timer"
)

func New(
	l *log.Logger, m *memory.Memory, d *display.Display, k *keyboard.Keyboard,
	dt, st *timer.Timer,
) *CPU {
	return &CPU{
		logger:   l,
		memory:   m,
		display:  d,
		keyboard: k,
		delay:    dt,
		sound:    st,
		fps:      500,
		dt:       0,
		reg:      [16]uint8{},
		stack:    [16]uint16{},
		i:        0,
		pc:       0x200,
		sp:       0,
	}
}

type CPU struct {
	logger   *log.Logger
	memory   *memory.Memory
	display  *display.Display
	keyboard *keyboard.Keyboard
	delay    *timer.Timer
	sound    *timer.Timer
	fps      uint
	dt       time.Duration

	reg   [16]uint8  // general purpose registers
	stack [16]uint16 // stack
	i     uint16     // index register
	pc    uint16     // program counter
	sp    uint8      // stack pointer
}

func (cpu *CPU) Tick(dt time.Duration) {
	cpu.dt += dt

	for cpu.dt >= time.Second/time.Duration(cpu.fps) {
		cpu.dt -= time.Second / time.Duration(cpu.fps)
		cpu.tick()
	}
}

func (cpu *CPU) tick() {
	// fetch opcode
	op := cpu.memory.ReadWord(cpu.pc)
	cpu.pc += 2

	dis := fmt.Sprintf("%04x\t%04x\t", cpu.pc-2, op)

	// decode and execute opcode
	addr := op & 0x0FFF
	m := op >> 12
	x := (op & 0x0F00) >> 8
	y := (op & 0x00F0) >> 4
	n := op & 0x000F
	kk := op & 0x00FF

	switch m {
	case 0x0:
		switch addr {
		case 0x0E0:
			// 00E0: CLS
			dis += "CLS"
			cpu.display.Clear()
		case 0x0EE:
			// 00EE: RET
			dis += "RET"
			cpu.pc = cpu.stack[cpu.sp]
			cpu.sp--
		}
	case 0x1:
		// 1nnn: JP addr
		dis += fmt.Sprintf("JP   %04x", addr)

		// Jump to self is a halt
		if addr == cpu.pc-2 {
			dis += " (HALT)"
		}

		cpu.pc = addr
	case 0x2:
		// 2nnn: CALL addr
		dis += fmt.Sprintf("CALL %04x", addr)

		cpu.sp++
		cpu.stack[cpu.sp] = cpu.pc
		cpu.pc = addr
	case 0x3:
		// 3xkk: SE Vx, byte
		dis += fmt.Sprintf("SE   V%x, %02x", x, kk)

		if cpu.reg[x] == uint8(kk) {
			cpu.pc += 2
		}
	case 0x4:
		// 4xkk: SNE Vx, byte
		// If values differ, skip the next opcode
		dis += fmt.Sprintf("SNE  V%x, %02x", x, kk)

		if cpu.reg[x] != uint8(kk) {
			cpu.pc += 2
		}
	case 0x5:
		if n == 0x0 {
			// 5xy0: SE Vx, Vy
			// If values are the same, skip the next opcode
			dis += fmt.Sprintf("SE   V%x, V%x", x, y)

			if cpu.reg[x] == cpu.reg[y] {
				cpu.pc += 2
			}
		}
	case 0x6:
		// 6xkk: LD Vx, byte
		dis += fmt.Sprintf("LD   V%x, %02x", x, kk)

		cpu.reg[x] = uint8(kk)
	case 0x7:
		// 7xkk: ADD Vx, byte
		dis += fmt.Sprintf("ADD  V%x, %02x", x, kk)

		cpu.reg[x] += uint8(kk)
	case 0x8:
		switch n {
		case 0x0:
			// 8xy0: LD Vx, Vy
			dis += fmt.Sprintf("LD   V%x, V%x", x, y)

			cpu.reg[x] = cpu.reg[y]
		case 0x1:
			// 8xy1: OR Vx, Vy
			dis += fmt.Sprintf("OR   V%x, V%x", x, y)

			cpu.reg[x] |= cpu.reg[y]
		case 0x2:
			// 8xy2: AND Vx, Vy
			dis += fmt.Sprintf("AND  V%x, V%x", x, y)

			cpu.reg[x] &= cpu.reg[y]
		case 0x3:
			// 8xy3: XOR Vx, Vy
			dis += fmt.Sprintf("XOR  V%x, V%x", x, y)

			cpu.reg[x] ^= cpu.reg[y]
		case 0x4:
			// 8xy4: ADD Vx, Vy
			dis += fmt.Sprintf("ADD  V%x, V%x", x, y)

			cpu.reg[x] += cpu.reg[y]

			// set carry flag
			if cpu.reg[x] < cpu.reg[y] {
				cpu.reg[0xF] = 1
			} else {
				cpu.reg[0xF] = 0
			}
		case 0x5:
			// 8xy5: SUB Vx, Vy
			dis += fmt.Sprintf("SUB  V%x, V%x", x, y)

			// set NOT borrow flag
			if cpu.reg[x] > cpu.reg[y] {
				cpu.reg[0xF] = 1
			} else {
				cpu.reg[0xF] = 0
			}

			// subtract
			cpu.reg[x] -= cpu.reg[y]
		case 0x6:
			// 8xy6: SHR Vx {, Vy}
			dis += fmt.Sprintf("SHR  V%x {, V%x}", x, y)

			// set carry flag
			if cpu.reg[x]&0x1 == 1 {
				cpu.reg[0xF] = 1
			} else {
				cpu.reg[0xF] = 0
			}

			// shift right
			cpu.reg[x] >>= 1
		case 0x7:
			// 8xy7: SUBN Vx, Vy
			dis += fmt.Sprintf("SUBN V%x, V%x", x, y)

			// set NOT borrow flag
			if cpu.reg[y] > cpu.reg[x] {
				cpu.reg[0xF] = 1
			} else {
				cpu.reg[0xF] = 0
			}

			// subtract
			cpu.reg[x] = cpu.reg[y] - cpu.reg[x]
		case 0xE:
			// 8xyE: SHL Vx {, Vy}
			dis += fmt.Sprintf("SHL  V%x {, V%x}", x, y)

			// set carry flag
			if cpu.reg[x]&0x80 == 0x80 {
				cpu.reg[0xF] = 1
			} else {
				cpu.reg[0xF] = 0
			}

			// shift left
			cpu.reg[x] <<= 1
		}
	case 0x9:
		switch n {
		case 0x0:
			// 9xy0: SNE Vx, Vy
			dis += fmt.Sprintf("SNE  V%x, V%x", x, y)

			if cpu.reg[x] != cpu.reg[y] {
				cpu.pc += 2
			}
		}
	case 0xA:
		// Annn: LD I, addr
		dis += fmt.Sprintf("LD   I, %04x", addr)

		cpu.i = addr
	case 0xB:
		// Bnnn: JP V0, addr
		dis += fmt.Sprintf("JP   V0, %04x", addr)

		cpu.pc = uint16(cpu.reg[0]) + uint16(addr)
	case 0xC:
		// Cxkk: RND Vx, byte
		// The interpreter generates a random number from 0 to 255, which is then
		// ANDed with the value kk. The results are stored in Vx.
		dis += fmt.Sprintf("RND  V%x, %02x", x, kk)

		rnd := uint16(rand.Intn(255))
		cpu.reg[x] = uint8(rnd & kk)
	case 0xD:
		// Dxyn: DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy),
		// set VF = collision.
		dis += fmt.Sprintf("DRW  V%x, V%x, %x", x, y, n)

		if cpu.display.Blit(uint16(cpu.reg[x]), uint16(cpu.reg[y]), cpu.memory.ReadRange(cpu.i, n)) {
			cpu.reg[0xF] = 1
		} else {
			cpu.reg[0xF] = 0
		}
	case 0xE:
		switch kk {
		case 0x9E:
			// Ex9E: SKP Vx
			dis += fmt.Sprintf("SKP  V%x", x)

			if cpu.keyboard.IsKeyPressed(cpu.reg[x]) {
				cpu.pc += 2
			}
		case 0xA1:
			// ExA1: SKNP Vx
			dis += fmt.Sprintf("SKNP V%x", x)

			if !cpu.keyboard.IsKeyPressed(cpu.reg[x]) {
				cpu.pc += 2
			}
		}
	case 0xF:
		switch kk {
		case 0x07:
			// Fx07: LD Vx, DT
			dis += fmt.Sprintf("LD   V%x, DT", x)

			cpu.reg[x] = cpu.delay.Get()
		case 0x0A:
			// Fx0A: LD Vx, K
			dis += fmt.Sprintf("LD   V%x, K", x)

			code, ok := cpu.keyboard.GetKeyPress()
			if ok {
				cpu.reg[x] = code
			} else {
				// no key was pressed, reset the program counter so we
				// try again on the next tick.
				cpu.pc -= 2
			}
		case 0x15:
			// Fx15: LD DT, Vx
			dis += fmt.Sprintf("LD   DT, V%x", x)

			cpu.delay.Set(cpu.reg[x])
		case 0x18:
			// Fx18: LD ST, Vx
			dis += fmt.Sprintf("LD   ST, V%x", x)

			cpu.sound.Set(cpu.reg[x])
		case 0x1E:
			// Fx1E: ADD I, Vx
			dis += fmt.Sprintf("ADD  I, V%x", x)

			cpu.i += uint16(cpu.reg[x])
		case 0x29:
			// Fx29: LD F, Vx
			// Set I = location of sprite for digit Vx.
			// (digit sprites are 5 bytes, starting at address 0)
			dis += fmt.Sprintf("LD   F, V%x", x)

			cpu.i = uint16(cpu.reg[x] * 5)
		case 0x33:
			// Fx33: LD B, Vx
			// Write Vx as BCD to memory starting at I
			dis += fmt.Sprintf("LD   B, V%x", x)

			vx := cpu.reg[x]
			cpu.memory.WriteByte(cpu.i, vx/100)
			cpu.memory.WriteByte(cpu.i+1, (vx/10)%10)
			cpu.memory.WriteByte(cpu.i+2, vx%10)
		case 0x55:
			// Fx55: LD [I], Vx
			// Write register V0..Vx into memory starting at I
			dis += fmt.Sprintf("LD   [I], V%x", x)

			for i := uint16(0); i <= uint16(x); i++ {
				cpu.memory.WriteByte(cpu.i+i, cpu.reg[i])
			}
		case 0x65:
			// Fx65: LD Vx, [I]
			// Read memory starting at I into register v0..Vx
			dis += fmt.Sprintf("LD   V%x, [I]", x)

			for i := uint16(0); i <= uint16(x); i++ {
				cpu.reg[i] = cpu.memory.ReadByte(cpu.i + i)
			}
		}
	}

	/*
		dis += fmt.Sprintf("\n\t\ti: %04x, v:[", cpu.i)
		for i := 0; i < len(cpu.reg); i++ {
			if i > 0 {
				dis += ", "
			}
			dis += fmt.Sprintf("%x: %02x", i, cpu.reg[i])
		}
		dis += "]"
	*/

	cpu.logger.Infof(dis)
}
