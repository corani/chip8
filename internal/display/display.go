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

func (d *Display) Blit(sx, sy uint16, sprite []uint8) bool {
	// all sprites are 8 pixels (bits) wide
	const width = uint16(8)
	height := uint16(len(sprite))

	for y := uint16(0); y < height; y++ {
		for x := uint16(0); x < width; x++ {
			// get the correct bit
			val := (sprite[y*height] >> (width - x - 1)) & 0x01
			// sprites need to wrap!
			d.Framebuffer[(sx+x)%width][(sy+y)%height] = val
		}
	}

	return false
}
