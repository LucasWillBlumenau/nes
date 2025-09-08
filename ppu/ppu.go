package ppu

import (
	"time"

	"github.com/LucasWillBlumenau/nes/interruption"
)

const kb = 1024
const scanlinesQuantity = 261

const statusVBlank uint8 = 0b10000000

type ppuPorts struct {
	Control1 uint8
	Control2 uint8
	Status   uint8
}

type PPU struct {
	oddFrame bool
	vram     []uint8
	chrRom   []uint8
	ppuPorts ppuPorts
}

func NewPPU(chrRom []uint8) *PPU {
	vram := make([]uint8, 2*kb)
	return &PPU{oddFrame: false, vram: vram, chrRom: chrRom}
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
	return 0
}

func (p *PPU) ReadOAMDataPort() uint8 {
	return 0
}

func (p *PPU) ReadVRamDataPort() uint8 {
	return 0
}

func (p *PPU) WritePPUControlPort(value uint8) {

}

func (p *PPU) WritePPUMaskPort(value uint8) {

}

func (p *PPU) WriteOAMAddrPort(value uint8) {

}

func (p *PPU) WriteOAMDataPort(value uint8) {

}

func (p *PPU) WritePPUScrollPort(value uint8) {

}

func (p *PPU) WritePPUAddrPort(value uint8) {

}

func (p *PPU) WritePPUDataPort(value uint8) {

}
