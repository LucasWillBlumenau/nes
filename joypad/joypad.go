package joypad

import (
	"sync"
)

type Button uint8

const (
	ButtonA      Button = 0
	ButtonB      Button = 1
	ButtonSelect Button = 2
	ButtonStart  Button = 3
	ButtonUp     Button = 4
	ButtonDown   Button = 5
	ButtonLeft   Button = 6
	ButtonRight  Button = 7
)

type Joypad struct {
	state      []bool
	savedState []bool
	strobe     bool
	readsCount uint8
	mu         sync.Mutex
}

func New() *Joypad {
	return &Joypad{
		state:      make([]bool, 8),
		savedState: make([]bool, 8),
		strobe:     false,
	}
}

func (j *Joypad) Write(value uint8) {
	j.mu.Lock()
	defer j.mu.Unlock()

	newStrobeValue := (value & 0b01) == 1
	if j.strobe && !newStrobeValue {
		copy(j.savedState, j.state)
	}
	j.strobe = newStrobeValue
	j.readsCount = 0
}

func (j *Joypad) Read() uint8 {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.readsCount > 7 || j.strobe {
		return 1
	}

	buttonPressed := j.savedState[j.readsCount]
	j.readsCount++

	if buttonPressed {
		return 1
	}
	return 0
}

func (j *Joypad) SetControl(button Button, value bool) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.state[button] = value
}
