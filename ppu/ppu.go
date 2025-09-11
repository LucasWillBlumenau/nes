package ppu

import (
	"time"

	"github.com/LucasWillBlumenau/nes/interruption"
)

const kb = 1024
const statusVBlank uint8 = 0b10000000

const (
	controlBaseNameTableAddr          = 0b00000011
	controlIncrementSize              = 0b00000100 // 0 = 1, 1 = 32
	controlSpritePatternTableAddr     = 0b00001000
	controlBackgroundPatternTableAddr = 0b00010000 // 0 = 0x0000, 1 = 0x1000
	controlSpriteSize                 = 0b00100000 // 0 = 8x8, 1 = 8x16
	controlMasterSlaveSelect          = 0b01000000
	controlPortNmiEnabled             = 0b10000000
)

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
	if (p.ppuPorts.Control & controlPortNmiEnabled) > 0 {
		interruption.InterruptionHandler <- interruption.NonMaskableInterrupt
	}
	p.ppuPorts.Status |= statusVBlank
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
	var incrementSize uint16
	if (p.ppuPorts.Control & controlIncrementSize) > 0 {
		incrementSize = 32
	} else {
		incrementSize = 1
	}

	p.bus.write(p.currentAddress, value)
	p.currentAddress += incrementSize
}
