package memory

func New() *Memory {
	return &Memory{}
}

type Memory struct {
	RAM [4096]byte
}

func (mem *Memory) Load(addr uint16, data []uint8) {
	for i := 0; i < len(data); i++ {
		mem.RAM[addr+uint16(i)] = data[i]
	}
}

func (mem *Memory) ReadByte(addr uint16) uint8 {
	return mem.RAM[addr]
}

func (mem *Memory) WriteByte(addr uint16, value uint8) {
	mem.RAM[addr] = value
}

func (mem *Memory) ReadWord(addr uint16) uint16 {
	return uint16(mem.RAM[addr])<<8 | uint16(mem.RAM[addr+1])
}

func (mem *Memory) ReadRange(start, length uint16) []byte {
	return mem.RAM[start : start+length]
}
