package joypad

import (
	"sync"
)

type Joypad struct {
	status         uint8
	bufferedStatus uint8
	updateStatus   bool
	mu             sync.Mutex
}

func New() *Joypad {
	return &Joypad{
		status:         0xFF,
		bufferedStatus: 0xFF,
		updateStatus:   false,
	}
}

func (j *Joypad) Write(value uint8) {
	j.status = value
	if j.canUpdate() {
		j.bufferedStatus = value
		j.status = 0xFF
	}
}

func (j *Joypad) Read() uint8 {
	value := (j.bufferedStatus ^ 0xFF) & 0b01
	j.bufferedStatus >>= 1
	return value
}

func (j *Joypad) SetStrobe(value uint8) {
	j.setUpdateStatus((value & 0b01) == 1)
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
		j.bufferedStatus = j.status
	}
	j.updateStatus = updateStatus
}
