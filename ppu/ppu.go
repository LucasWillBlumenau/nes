package ppu

import (
	"image/color"

	"github.com/LucasWillBlumenau/nes/interrupt"
	"github.com/LucasWillBlumenau/nes/window"
)

const kb = 1024
const statusVBlank uint8 = 0b10000000
const visibleScalines = 240
const totalScanlines = 261
const dotsPerScanline = 341

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
	window *window.Window
	bus    bus

	tiles []Tile

	ppuPorts       ppuPorts
	writeLatch     bool
	oddFrame       bool
	currentAddress uint16
	bufferedData   uint8

	coarseX uint8
	coarseY uint8
	fineX   uint8
	fineY   uint8

	currentScanline uint16
	currentDot      uint16
	VBlankCount     uint
	ElapseCycles    uint
}

func NewPPU(window *window.Window, chrRom []uint8) *PPU {
	bus := newBus(chrRom)
	tiles := generateTilesFromChrRom(chrRom)
	return &PPU{
		oddFrame:   false,
		writeLatch: false,
		bus:        bus,
		tiles:      tiles,
		window:     window,
	}
}

func (p *PPU) ElapseCPUCycles(cpuCycles uint8) {
	remaingPPUCycles := cpuCycles * 3
	for remaingPPUCycles > 0 {
		p.performScanlineStep()
		remaingPPUCycles--
		p.ElapseCycles++
	}

}

func (p *PPU) performScanlineStep() {
	if p.currentScanline == 240 && p.currentDot == 0 {
		p.outputCurrentNameTable()
		p.vBlank()
		p.currentDot = 1
		return
	}

	if p.currentDot == (dotsPerScanline - 1) {
		if p.currentScanline == (totalScanlines - 1) {
			p.currentScanline = 0
		} else {
			p.currentScanline++
		}
		p.currentDot = 0
		return
	}
	p.currentDot++
}

func (p *PPU) vBlank() {
	if (p.ppuPorts.Control & controlPortNmiEnabled) > 0 {
		interrupt.InterruptSignal.Send(interrupt.NonMaskableInterrupt)
	}
	p.ppuPorts.Status |= statusVBlank
	p.VBlankCount++
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
	data := p.bufferedData
	p.bufferedData = p.bus.read(p.currentAddress)
	return data
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
	if p.writeLatch {
		p.coarseY = (value & 0b11111000) >> 3
		p.fineY = (value & 0b00000111)
	} else {
		p.coarseX = (value & 0b11111000) >> 3
		p.fineX = (value & 0b00000111)
	}
	p.writeLatch = !p.writeLatch
}

func (p *PPU) WritePPUAddrPort(value uint8) {
	if p.writeLatch {
		p.currentAddress = (p.currentAddress & 0xFF00) | (uint16(value))
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

type Tile [8][8]uint8

func generateTilesFromChrRom(rom []uint8) []Tile {

	var tiles []Tile
	offset := 0
	length := 8
	for offset < len(rom) {
		leastSignificantBits := rom[offset : offset+length]
		offset += length
		mostSignificantBits := rom[offset : offset+length]
		offset += length

		var tile Tile
		for i := range length {
			loBits := leastSignificantBits[i]
			hiBits := mostSignificantBits[i]

			for j := range 8 {
				shiftSize := 7 - uint8(j)
				lo := (loBits >> shiftSize) & 0b00000001
				hi := (hiBits >> shiftSize) & 0b00000001
				tile[i][j] = (hi << 1) + lo
			}
		}
		tiles = append(tiles, tile)
	}

	return tiles
}

func (p *PPU) outputCurrentNameTable() {
	image := make([]color.RGBA, p.window.Width()*p.window.Height())
	nameTable := p.bus.ram[0:960]
	attrTable := p.bus.ram[960:1024]
	palettes := loadPalettes(p.bus.backgroundPalette)
	tiles := p.tiles[:256]
	if (p.ppuPorts.Control & controlBackgroundPatternTableAddr) > 0 {
		tiles = p.tiles[256:]
	}

	for i, tileIdx := range nameTable {
		tile := tiles[tileIdx]
		tileX := (i % 32) * 8
		tileY := (i / 32) * 8

		attrTableX := tileX / 32
		attrTableY := tileY / 32
		attrTableByte := attrTable[attrTableX+attrTableY*8]
		attrTableByteX := attrTableX % 2
		attrTableByteY := attrTableY % 2
		paletteId := (attrTableByte >> uint8(attrTableByteY<<1|attrTableByteX)) & 0b11

		for y := range 8 {
			for x := range 8 {
				index := tile[y][x]
				color := palettes.Read(paletteId, index)
				image[(tileX+x)+(tileY+y)*p.window.Width()] = color
			}
		}
	}
	p.window.UpdateImageBuffer(image)
}
