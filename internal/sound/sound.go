package sound

import (
	"time"

	"github.com/corani/chip-8/internal/timer"
)

func New(timer *timer.Timer) *Sound {
	return &Sound{timer: timer}
}

type Sound struct {
	timer *timer.Timer
}

func (s *Sound) Tick(dt time.Duration) {
	s.timer.Tick(dt)
}

func (s *Sound) SetActive(duration uint8) {
	s.timer.Set(duration)
}

func (s *Sound) Active() bool {
	return s.timer.Active()
}
