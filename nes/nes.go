package nes

import (
	"image"
	"os"
	"time"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/ppu"
)

const cpuCycleDuration int64 = 559

type NES struct {
	Frames chan image.RGBA
	ppu    *ppu.PPU
	cpu    *cpu.CPU
}

func NewNES(
	frames chan image.RGBA,
	path string,
	scaleFactor int,
	joypadOne *joypad.Joypad,
	joypadTwo *joypad.Joypad,
) (*NES, error) {
	f, err := os.Create("output.log")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cart, err := cartridge.LoadCartridge(path)
	if err != nil {
		return nil, err
	}

	ppuBus := ppu.NewPPUBus(cart)
	ppu := ppu.NewPPU(ppuBus, frames, scaleFactor)
	bus := bus.NewBus(ppu, cart, joypadOne, joypadTwo)
	cpu := cpu.NewCPU(bus)

	return &NES{
		Frames: frames,
		ppu:    ppu,
		cpu:    cpu,
	}, nil
}

func (n *NES) Run() {
	n.cpu.Reset()
	start := time.Now()
	for {
		cyclesTaken, err := n.cpu.Run()
		if err != nil {
			panic(err)
		}
		ppuCycles := cyclesTaken * 3
		n.ppu.RunSteps(ppuCycles)

		currentTime := time.Now()
		elapsedTime := currentTime.UnixNano() - start.UnixNano()
		expectedElapsedTime := n.cpu.ElapsedCycles() * cpuCycleDuration
		if expectedElapsedTime > elapsedTime {
			diff := time.Duration(expectedElapsedTime - elapsedTime)
			time.Sleep(diff)
		}
	}
}
