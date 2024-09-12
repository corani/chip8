package keyboard

func New() *Keyboard {
	return &Keyboard{
		pressed: nil,
	}
}

type Keyboard struct {
	pressed []uint8
}

func (k *Keyboard) KeyDown(code uint8) {
	// check if key is already pressed
	for _, key := range k.pressed {
		if key == code {
			return
		}
	}

	k.pressed = append(k.pressed, code)
}

func (k *Keyboard) KeyUp(code uint8) {
	// remove key from pressed if found
	for i, key := range k.pressed {
		if key == code {
			k.pressed = append(k.pressed[:i], k.pressed[i+1:]...)
			return
		}
	}
}

func (k *Keyboard) IsKeyPressed(code uint8) bool {
	for _, key := range k.pressed {
		if key == code {
			return true
		}
	}

	return false
}

func (k *Keyboard) GetKeyPress() (uint8, bool) {
	// return the first key pressed
	if len(k.pressed) > 0 {
		return k.pressed[0], true
	}

	return 0, false
}
