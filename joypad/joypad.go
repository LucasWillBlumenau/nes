package joypad

import (
	"sync"
)

type Joypad struct {
	bufferedState uint8
	status        uint8
	updateStatus  bool
	mu            sync.Mutex
}

func New() *Joypad {
	return &Joypad{
		status:       0x00,
		updateStatus: false,
	}
}

func (j *Joypad) Write(value uint8) uint8 {
	if value > 0 {
		if j.canUpdate() {
			j.status |= value
		} else {
			j.bufferedState |= value
		}
		return value ^ 0xFF
	}

	return 0xFF
}

func (j *Joypad) Read() uint8 {
	value := j.status & 0b01
	j.status >>= 1
	return value
}

func (j *Joypad) SetStrobe(value uint8) {
	f := (value & 0b01) == 1
	j.setUpdateStatus(f)
}

func (j *Joypad) canUpdate() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.updateStatus
}

func (j *Joypad) setUpdateStatus(updateStatus bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if updateStatus {
		j.status = j.bufferedState
		j.bufferedState = 0
	}
	j.updateStatus = updateStatus
}
