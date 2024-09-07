package display

import "github.com/charmbracelet/log"

func New(logger *log.Logger) *Display {
	width := 64
	height := 32

	fb := make([][]uint8, width)
	for x := 0; x < width; x++ {
		fb[x] = make([]uint8, height)
	}

	return &Display{
		logger:      logger,
		width:       64,
		height:      32,
		Framebuffer: fb,
	}
}

type Display struct {
	logger      *log.Logger
	width       int
	height      int
	Framebuffer [][]uint8
}

func (d *Display) Clear() {
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			d.Framebuffer[x][y] = 0
		}
	}
}

func (d *Display) Blit(sx, sy uint16, sprite []uint8) bool {
	// all sprites are 8 pixels (bits) wide
	const width = uint16(8)
	height := uint16(len(sprite))

	res := false

	for y := uint16(0); y < height; y++ {
		for x := uint16(0); x < width; x++ {
			// get the correct bit
			val := (sprite[y] >> (width - x - 1)) & 0x01

			// sprites need to wrap!
			px := (sx + x) % uint16(d.width)
			py := (sy + y) % uint16(d.height)

			if val != d.Framebuffer[px][py] {
				res = true
			}

			d.Framebuffer[px][py] ^= val
		}
	}

	return res
}
