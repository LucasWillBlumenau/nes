package nes

import (
	"fmt"
	"image"
	"time"

	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/ppu"
)

const cpuCycleDuration int64 = 559

type NES struct {
	Frames           chan image.RGBA
	Paused           bool
	ShowInstructions bool
	ppu              *ppu.PPU
	cpu              *cpu.CPU
	start            time.Time
}

func (n *NES) TogglePrintXCoordinatesOfTopScanlineFlag() {
	n.ppu.TogglePrintXCoordinatesOfTopScanlineFlag()
}

func NewNES(
	frames chan image.RGBA,
	romPath string,
	scaleFactor int,
	joypadOne *joypad.Joypad,
	joypadTwo *joypad.Joypad,
) (*NES, error) {
	cart, err := cartridge.LoadCartridgeFromRom(romPath)
	if err != nil {
		return nil, err
	}

	ppuBus := ppu.NewPPUBus(cart)
	ppu := ppu.NewPPU(ppuBus, frames, scaleFactor)
	bus := cpu.NewBus(ppu, cart, joypadOne, joypadTwo)
	cpu := cpu.NewCPU(bus)

	return &NES{
		Frames: frames,
		Paused: false,
		ppu:    ppu,
		cpu:    cpu,
	}, nil
}

func (n *NES) Run() {
	n.cpu.Reset()
	n.start = time.Now()
	for {
		if n.Paused {
			time.Sleep(100 * time.Millisecond)
		} else {
			n.runStep()
		}
	}
}

func (n *NES) runStep() {
	var cyclesTaken uint16
	var err error
	if n.ShowInstructions {
		fmt.Println(n.cpu.State())
		cyclesTaken, err = n.cpu.Run()
		fmt.Println(n.cpu.GetLastInstruction())
	} else {
		cyclesTaken, err = n.cpu.Run()
	}

	if err != nil {
		panic(err)
	}
	ppuCycles := cyclesTaken * 3
	n.ppu.RunSteps(ppuCycles)

	currentTime := time.Now()
	elapsedTime := currentTime.UnixNano() - n.start.UnixNano()
	expectedElapsedTime := n.cpu.ElapsedCycles() * cpuCycleDuration
	if expectedElapsedTime > elapsedTime {
		diff := time.Duration(expectedElapsedTime - elapsedTime)
		time.Sleep(diff)
	}
}

func (n *NES) ExecuteNextInstruction() {
	if !n.Paused {
		return
	}
	n.runStep()
}

func (n *NES) DumpPPUNametables() image.Image {
	return n.ppu.DumpNametables()
}
