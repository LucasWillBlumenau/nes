package ppu

import (
	"time"

	"github.com/LucasWillBlumenau/nes/interruption"
)

const kb = 1024
const statusVBlank uint8 = 0b10000000

type ppuPorts struct {
	Control uint8
	Mask    uint8
	Status  uint8
}

type PPU struct {
	oddFrame       bool
	vram           []uint8
	ppuPorts       ppuPorts
	writeLatch     bool
	currentAddress uint16
	bus            bus
}

func NewPPU(chrRom []uint8) *PPU {
	vram := make([]uint8, 2*kb)
	return &PPU{oddFrame: false, vram: vram, writeLatch: false, bus: newBus(chrRom)}
}

func (p *PPU) Run() {
	for {
		p.generateFrame()
	}
}

func (p *PPU) generateFrame() {
	for scanlineNumber := range 240 {
		// perform scanline
		p.performScanline(scanlineNumber)
	}
	p.vBlank()
}

func (p *PPU) performScanline(number int) {
}

func (p *PPU) vBlank() {
	p.ppuPorts.Status |= statusVBlank
	interruption.InterruptionHandler <- interruption.NonMaskableInterrupt
	time.Sleep(time.Millisecond * 50)
}

func (p *PPU) ReadStatusPort() uint8 {
	currentStatus := p.ppuPorts.Status
	p.ppuPorts.Status &= statusVBlank ^ 0xFF
	p.writeLatch = false

	return currentStatus
}

func (p *PPU) ReadOAMDataPort() uint8 {
	return 0
}

func (p *PPU) ReadVRamDataPort() uint8 {
	return 0
}

func (p *PPU) WritePPUControlPort(value uint8) {
	p.ppuPorts.Control = value
}

func (p *PPU) WritePPUMaskPort(value uint8) {
	p.ppuPorts.Mask = value
}

func (p *PPU) WriteOAMAddrPort(value uint8) {

}

func (p *PPU) WriteOAMDataPort(value uint8) {

}

func (p *PPU) WritePPUScrollPort(value uint8) {

}

func (p *PPU) WritePPUAddrPort(value uint8) {
	if p.writeLatch {
		p.currentAddress += uint16(value)
	} else {
		p.currentAddress = uint16(value) << 8
	}
	p.writeLatch = !p.writeLatch
}

func (p *PPU) WritePPUDataPort(value uint8) {
	incrementSizeFlag := (p.ppuPorts.Control & 00000100) > 0
	var incrementSize uint16
	if incrementSizeFlag {
		incrementSize = 32
	} else {
		incrementSize = 1
	}

	p.bus.write(p.currentAddress, value)
	p.currentAddress += incrementSize
}
