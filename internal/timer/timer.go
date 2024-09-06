package timer

import "time"

func New() *Timer {
	return &Timer{
		count: 0,
		fps:   60,
		dt:    0,
	}
}

type Timer struct {
	count uint8
	fps   uint8
	dt    time.Duration
}

func (t *Timer) Tick(dt time.Duration) {
	if t.count == 0 {
		return
	}

	// accumulate delta time
	t.dt += dt

	// decrement count at one frame per second
	for t.dt >= time.Second/time.Duration(t.fps) {
		t.count--
		t.dt -= time.Second / time.Duration(t.fps)
	}
}

func (t *Timer) Set(c uint8) {
	t.count = c
}

func (t *Timer) Get() uint8 {
	return t.count
}

func (t *Timer) Active() bool {
	return t.count > 0
}
