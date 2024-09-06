package cpu

import (
	"time"

	"github.com/corani/chip-8/internal/display"
	"github.com/corani/chip-8/internal/keyboard"
	"github.com/corani/chip-8/internal/memory"
	"github.com/corani/chip-8/internal/timer"
)

func New(m *memory.Memory, d *display.Display, k *keyboard.Keyboard, dt, st *timer.Timer) *CPU {
	return &CPU{
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
		pc:       0,
		sp:       0,
	}
}

type CPU struct {
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
	// fetch opcode
	op := cpu.memory.ReadWord(cpu.pc)
	cpu.pc += 2

	// decode opcode
	_ = op

	if op == 0x00E0 {
		// 00E0: CLS
		cpu.display.Clear()
	} else if op == 0x00EE {
		// 00EE: RET
		cpu.pc = cpu.stack[cpu.sp]
		cpu.sp--
	} else if op&0xF000 == 0x1000 {
		// 1nnn: JP addr
		cpu.pc = op & 0x0FFF
	} else if op&0xF000 == 0x2000 {
		// 2nnn: CALL addr
		cpu.sp++
		cpu.stack[cpu.sp] = cpu.pc
		cpu.pc = op & 0x0FFF
	} else if op&0xF000 == 0x3000 {
		// 3xkk: SE Vx, byte
		if cpu.reg[(op&0x0F00)>>8] == uint8(op&0x00FF) {
			cpu.pc += 2
		}
	} else if op&0xF000 == 0x4000 {
		// 4xkk: SNE Vx, byte
		if cpu.reg[(op&0x0F00)>>8] != uint8(op&0x00FF) {
			cpu.pc += 2
		}
	} else if op&0xF00F == 0x5000 {
		// 5xy0: SE Vx, Vy
		if cpu.reg[(op&0x0F00)>>8] == cpu.reg[(op&0x00F0)>>4] {
			cpu.pc += 2
		}
	} else if op&0xF000 == 0x6000 {
		// 6xkk: LD Vx, byte
		cpu.reg[(op&0x0F00)>>8] = uint8(op & 0x00FF)
	} else if op&0xF000 == 0x7000 {
		// 7xkk: ADD Vx, byte
		cpu.reg[(op&0x0F00)>>8] += uint8(op & 0x00FF)
	} else if op&0xF000 == 0x8000 {
		// 8xy0: LD Vx, Vy
		cpu.reg[(op&0x0F00)>>8] = cpu.reg[(op&0x00F0)>>4]
	} else if op&0xF00F == 0x8001 {
		// 8xy1: OR Vx, Vy
		cpu.reg[(op&0x0F00)>>8] |= cpu.reg[(op&0x00F0)>>4]
	} else if op&0xF00F == 0x8002 {
		// 8xy2: AND Vx, Vy
		cpu.reg[(op&0x0F00)>>8] &= cpu.reg[(op&0x00F0)>>4]
	} else if op&0xF00F == 0x8003 {
		// 8xy3: XOR Vx, Vy
		cpu.reg[(op&0x0F00)>>8] ^= cpu.reg[(op&0x00F0)>>4]
	} else if op&0xF00F == 0x8004 {
		// 8xy4: ADD Vx, Vy
		cpu.reg[(op&0x0F00)>>8] += cpu.reg[(op&0x00F0)>>4]
		// set carry flag
		if cpu.reg[(op&0x0F00)>>8] < cpu.reg[(op&0x00F0)>>4] {
			cpu.reg[0xF] = 1
		} else {
			cpu.reg[0xF] = 0
		}
	} else if op&0xF00F == 0x8005 {
		// 8xy5: SUB Vx, Vy
		// set NOT borrow flag
		if cpu.reg[(op&0x0F00)>>8] > cpu.reg[(op&0x00F0)>>4] {
			cpu.reg[0xF] = 1
		} else {
			cpu.reg[0xF] = 0
		}
		// subtract
		cpu.reg[(op&0x0F00)>>8] -= cpu.reg[(op&0x00F0)>>4]
	} else if op&0xF00F == 0x8006 {
		// 8xy6: SHR Vx {, Vy}
		if cpu.reg[(op&0x0F00)>>8]&0x1 == 1 {
			cpu.reg[0xF] = 1
		} else {
			cpu.reg[0xF] = 0
		}
		cpu.reg[(op&0x0F00)>>8] >>= 1
	} else if op&0xF00F == 0x8007 {
		// 8xy7: SUBN Vx, Vy
		// set NOT borrow flag
		if cpu.reg[(op&0x00F0)>>4] > cpu.reg[(op&0x0F00)>>8] {
			cpu.reg[0xF] = 1
		} else {
			cpu.reg[0xF] = 0
		}
		cpu.reg[(op&0x0F00)>>8] = cpu.reg[(op&0x00F0)>>4] - cpu.reg[(op&0x0F00)>>8]
	} else if op&0xF00F == 0x800E {
		// 8xyE: SHL Vx {, Vy}
		// set carry flag
		if cpu.reg[(op&0x0F00)>>8]&0x80 == 1 {
			cpu.reg[0xF] = 1
		} else {
			cpu.reg[0xF] = 0
		}
		cpu.reg[(op&0x0F00)>>8] <<= 1
	} else if op&0xF00F == 0x9000 {
		// 9xy0: SNE Vx, Vy
		if cpu.reg[(op&0x0F00)>>8] != cpu.reg[(op&0x00F0)>>4] {
			cpu.pc += 2
		}
	} else if op&0xF000 == 0xA000 {
		// Annn: LD I, addr
		cpu.i = op & 0x0FFF
	} else if op&0xF000 == 0xB000 {
		// Bnnn: JP V0, addr
		cpu.pc = uint16(cpu.reg[0]) + uint16(op&0x0FFF)
	} else if op&0xF000 == 0xC000 {
		// Cxkk: RND Vx, byte
		// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx.
	} else if op&0xF000 == 0xD000 {
		// Dxyn: DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
	} else if op&0xF0FF == 0xE09E {
		// Ex9E: SKP Vx
		if cpu.keyboard.IsKeyPressed(cpu.reg[(op&0x0F00)>>8]) {
			cpu.pc += 2
		}
	} else if op&0xF0FF == 0xE0A1 {
		// ExA1: SKNP Vx
		if !cpu.keyboard.IsKeyPressed(cpu.reg[(op&0x0F00)>>8]) {
			cpu.pc += 2
		}
	} else if op&0xF0FF == 0xF007 {
		// Fx07: LD Vx, DT
		cpu.reg[(op&0x0F00)>>8] = cpu.delay.Get()
	} else if op&0xF0FF == 0xF00A {
		// Fx0A: LD Vx, K
		// Wait for a key press, store the value of the key in Vx.
		cpu.reg[(op&0x0F00)>>8] = cpu.keyboard.WaitForKey()
	} else if op&0xF0FF == 0xF015 {
		// Fx15: LD DT, Vx
		cpu.delay.Set(cpu.reg[(op&0x0F00)>>8])
	} else if op&0xF0FF == 0xF018 {
		// Fx18: LD ST, Vx
		cpu.sound.Set(cpu.reg[(op&0x0F00)>>8])
	} else if op&0xF0FF == 0xF01E {
		// Fx1E: ADD I, Vx
		cpu.i += uint16(cpu.reg[(op&0x0F00)>>8])
	} else if op&0xF0FF == 0xF029 {
		// Fx29: LD F, Vx
		// Set I = location of sprite for digit Vx.
	} else if op&0xF0FF == 0xF033 {
		// Fx33: LD B, Vx
		val := cpu.reg[(op&0x0F00)>>8]
		cpu.memory.WriteByte(cpu.i, val/100)
		cpu.memory.WriteByte(cpu.i+1, (val/10)%10)
		cpu.memory.WriteByte(cpu.i+2, val%10)
	} else if op&0xF0FF == 0xF055 {
		// Fx55: LD [I], Vx
		val := (op & 0x0F00) >> 8

		for i := uint16(0); i <= val; i++ {
			cpu.memory.WriteByte(cpu.i+i, cpu.reg[i])
		}
	} else if op&0xF0FF == 0xF065 {
		// Fx65: LD Vx, [I]
		val := (op & 0x0F00) >> 8

		for i := uint16(0); i <= val; i++ {
			cpu.reg[i] = cpu.memory.ReadByte(cpu.i + i)
		}
	}

	// execute opcode
}
