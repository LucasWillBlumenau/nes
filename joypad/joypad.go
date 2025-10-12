package joypad

import (
	"sync"
)

var nextId int = 1

type Joypad struct {
	id         int
	state      []bool
	savedState []bool
	strobe     bool
	readsCount uint8
	mu         sync.Mutex
}

func New() *Joypad {
	defer func() {
		nextId++
	}()
	return &Joypad{
		id:         nextId,
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

func (j *Joypad) SetControl(index uint8, value bool) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.state[index] = value
}
