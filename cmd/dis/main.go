package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type ROM []uint16

func load(filename string) (ROM, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(bs)%2 != 0 {
		return nil, err
	}

	rom := ROM{}

	for i := 0; i < len(bs); i += 2 {
		rom = append(rom, uint16(bs[i])<<8|uint16(bs[i+1]))
	}

	return rom, nil
}

func generateLabels(rom ROM) map[uint16]string {
	labels := map[uint16]string{}

	count := 0

	for _, instr := range rom {
		m := instr >> 12
		addr := instr & 0x0FFF

		switch m {
		case 0x1:
			// 1nnn: JP addr
			labels[addr] = fmt.Sprintf("label%02d", count)
			count++
		case 0x2:
			// 2nnn: CALL addr
			labels[addr] = fmt.Sprintf("routine%02d", count)
			count++
		case 0xB:
			// Bnnn: JP V0, addr
			labels[addr] = fmt.Sprintf("table%02d", count)
			count++
		}
	}

	return labels
}

func disassemble(rom ROM, labels map[uint16]string) {
	for i, instr := range rom {
		m := instr >> 12
		x := (instr & 0x0F00) >> 8
		y := (instr & 0x00F0) >> 4
		n := instr & 0x000F
		kk := instr & 0x00FF
		addr := instr & 0x0FFF

		if label, ok := labels[uint16(i+0x200)]; ok {
			fmt.Printf("%s:\n", label)
		}

		dis := fmt.Sprintf("%04x\t%04x\t", i+0x200, instr)

		switch m {
		case 0x0:
			switch addr {
			case 0x0E0:
				// 00E0: CLS
				dis += "CLS"
			case 0x0EE:
				// 00EE: RET
				dis += "RET"
			}
		case 0x1:
			// 1nnn: JP addr
			dis += fmt.Sprintf("JP   %04x", addr)

			if label, ok := labels[addr]; ok {
				dis += fmt.Sprintf(" ; %s", label)
			}
		case 0x2:
			// 2nnn: CALL addr
			dis += fmt.Sprintf("CALL %04x", addr)

			if label, ok := labels[addr]; ok {
				dis += fmt.Sprintf(" ; %s", label)
			}
		case 0x3:
			// 3xkk: SE Vx, byte
			dis += fmt.Sprintf("SE   V%01x, %02x", x, kk)
		case 0x4:
			// 4xkk: SNE Vx, byte
			dis += fmt.Sprintf("SNE  V%01x, %02x", x, kk)
		case 0x5:
			if n == 0x0 {
				// 5xy0: SE Vx, Vy
				dis += fmt.Sprintf("SE   V%01x, V%01x", x, y)
			}
		case 0x6:
			// 6xkk: LD Vx, byte
			dis += fmt.Sprintf("LD   V%01x, %02x", x, kk)
		case 0x7:
			// 7xkk: ADD Vx, byte
			dis += fmt.Sprintf("ADD  V%01x, %02x", x, kk)
		case 0x8:
			switch n {
			case 0x0:
				// 8xy0: LD Vx, Vy
				dis += fmt.Sprintf("LD   V%x, V%x", x, y)
			case 0x1:
				// 8xy1: OR Vx, Vy
				dis += fmt.Sprintf("OR   V%x, V%x", x, y)
			case 0x2:
				// 8xy2: AND Vx, Vy
				dis += fmt.Sprintf("AND  V%x, V%x", x, y)
			case 0x3:
				// 8xy3: XOR Vx, Vy
				dis += fmt.Sprintf("XOR  V%x, V%x", x, y)
			case 0x4:
				// 8xy4: ADD Vx, Vy
				dis += fmt.Sprintf("ADD  V%x, V%x", x, y)
			case 0x5:
				// 8xy5: SUB Vx, Vy
				dis += fmt.Sprintf("SUB  V%x, V%x", x, y)
			case 0x6:
				// 8xy6: SHR Vx {, Vy}
				dis += fmt.Sprintf("SHR  V%x {, V%x}", x, y)
			case 0x7:
				// 8xy7: SUBN Vx, Vy
				dis += fmt.Sprintf("SUBN V%x, V%x", x, y)
			case 0xE:
				// 8xyE: SHL Vx {, Vy}
				dis += fmt.Sprintf("SHL  V%x {, V%x}", x, y)
			}
		case 0x9:
			if n == 0x0 {
				// 9xy0: SNE Vx, Vy
				dis += fmt.Sprintf("SNE  V%x, V%x", x, y)
			}
		case 0xA:
			// Annn: LD I, addr
			dis += fmt.Sprintf("LD   I, %04x", addr)
		case 0xB:
			// Bnnn: JP V0, addr
			dis += fmt.Sprintf("JP   V0, %04x", addr)

			if label, ok := labels[addr]; ok {
				dis += fmt.Sprintf(" ; %s", label)
			}
		case 0xC:
			// Cxkk: RND Vx, byte
			dis += fmt.Sprintf("RND  V%x, %02x", x, kk)
		case 0xD:
			// Dxyn: DRW Vx, Vy, nibble
			dis += fmt.Sprintf("DRW  V%x, V%x, %x", x, y, n)
		case 0xE:
			switch kk {
			case 0x9E:
				// Ex9E: SKP Vx
				dis += fmt.Sprintf("SKP  V%x", x)
			case 0xA1:
				// ExA1: SKNP Vx
				dis += fmt.Sprintf("SKNP V%x", x)
			}
		case 0xF:
			switch kk {
			case 0x07:
				// Fx07: LD Vx, DT
				dis += fmt.Sprintf("LD   V%x, DT", x)
			case 0x0A:
				// Fx0A: LD Vx, K
				dis += fmt.Sprintf("LD   V%x, K", x)
			case 0x15:
				// Fx15: LD DT, Vx
				dis += fmt.Sprintf("LD   DT, V%x", x)
			case 0x18:
				// Fx18: LD ST, Vx
				dis += fmt.Sprintf("LD   ST, V%x", x)
			case 0x1E:
				// Fx1E: ADD I, Vx
				dis += fmt.Sprintf("ADD  I, V%x", x)
			case 0x29:
				// Fx29: LD F, Vx
				dis += fmt.Sprintf("LD   F, V%x", x)
			case 0x33:
				// Fx33: LD B, Vx
				dis += fmt.Sprintf("LD   B, V%x", x)
			case 0x55:
				// Fx55: LD [I], Vx
				dis += fmt.Sprintf("LD   [I], V%x", x)
			case 0x65:
				// Fx65: LD Vx, [I]
				dis += fmt.Sprintf("LD   V%x, [I]", x)
			}
		}

		fmt.Println(dis)
	}
}

func main() {
	logger := log.New(os.Stdout)
	logger.SetReportTimestamp(true)

	if len(os.Args) < 2 {
		logger.Errorf("Usage: %s <source.ch8>", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	rom, err := load(filename)
	if err != nil {
		logger.Errorf("failed to load rom: %v", err)
		os.Exit(1)
	}

	labels := generateLabels(rom)

	disassemble(rom, labels)
}
