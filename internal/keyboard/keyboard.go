package keyboard

import "time"

func New() *Keyboard {
	// TODO(daniel): `hold` needs to be tuned.
	// TODO(daniel): mutex for pressed.
	return &Keyboard{
		pressed: nil,
		hold:    30 * time.Millisecond,
		dt:      0,
	}
}

type Keyboard struct {
	pressed *uint8
	hold    time.Duration
	dt      time.Duration
}

func (k *Keyboard) KeyPress(code uint8) {
	k.pressed = &code
}

func (k *Keyboard) Tick(dt time.Duration) {
	if k.pressed == nil {
		return
	}

	// accumulate delta time
	k.dt += dt

	// if key is held down for `hold` duration, reset the keypress
	if k.dt >= k.hold {
		k.pressed = nil
		k.dt = 0
	}
}

func (k *Keyboard) IsKeyPressed(code uint8) bool {
	if k.pressed == nil {
		return false
	}

	return *k.pressed == code
}

func (k *Keyboard) GetKeyPress() (uint8, bool) {
	if k.pressed == nil {
		return 0, false
	}

	return *k.pressed, true
}
