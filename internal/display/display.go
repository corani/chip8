package display

func New() *Display {
	var fb [64][32]byte

	// checkerboard pattern
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			fb[x][y] = uint8((x + y) % 2)
		}
	}

	return &Display{
		Framebuffer: fb,
	}
}

type Display struct {
	Framebuffer [64][32]byte
}

func (d *Display) Clear() {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			d.Framebuffer[x][y] = 0
		}
	}
}
