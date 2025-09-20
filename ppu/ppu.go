package ppu

import (
	"image/color"

	"github.com/LucasWillBlumenau/nes/interrupt"
)

const (
	tileSize           uint16 = 16
	highBitPlaneOffset uint16 = 8

	enableStatusVBlank        uint8 = 0b10000000
	disableStatusVBlank       uint8 = 0b01111111
	disableSprite0HitFlag     uint8 = 0b10111111
	disableSpriteOverflowFlag uint8 = 0b11011111
)

type ppuPorts struct {
	control ppuControl
	mask    uint8
	status  uint8
}

type ppuRegisters struct {
	coarseY      uint16
	coarseX      uint16
	fineX        uint16
	fineY        uint16
	currentAddr  uint16
	bufferedData uint8
	writeLatch   bool
	pixelBuffer  pixelsShiftRegister
}

type pixelFetchingState uint8

const (
	stateFetchNametableFirstCycle pixelFetchingState = iota
	stateFetchNametableSecondCycle
	stateFetchAttrtableFirstCycle
	stateFetchAttrtableSecondCycle
	stateFetchPatternTableLowPlaneFirstCycle
	stateFetchPatternTableLowPlaneSecondCycle
	stateFetchPatternTableHighPlaneFirstCycle
	stateFetchPatternTableHighPlaneSecondCycle
	stateFetchesDone
)

type ppuRenderingState struct {
	scanline      uint16
	clock         uint16
	fetchingState pixelFetchingState
	tileIndex     uint8
	lowBitPlane   uint8
	highBitPlane  uint8
	paletteId     uint8
	currentAddr   uint16
}

type PPU struct {
	bus            *bus
	bufferedImage  [256 * 240]color.RGBA
	currentPixel   uint32
	frameDone      bool
	ports          ppuPorts
	renderingState ppuRenderingState
	registers      ppuRegisters
}

func NewPPU(chrRom []uint8) *PPU {
	bus := newBus(chrRom)
	return &PPU{
		bus: bus,
		registers: ppuRegisters{
			writeLatch:   false,
			currentAddr:  0x0000,
			bufferedData: 0x00,
		},
		renderingState: ppuRenderingState{
			scanline:      261,
			clock:         0,
			fetchingState: stateFetchNametableFirstCycle,
		},
	}
}

func (p *PPU) ReadStatusPort() uint8 {
	currentStatus := p.ports.status
	p.ports.status &= enableStatusVBlank ^ 0xFF
	p.registers.writeLatch = false

	return currentStatus
}

func (p *PPU) ReadOAMDataPort() uint8 {
	return 0
}

func (p *PPU) ReadVRamDataPort() uint8 {
	data := p.registers.bufferedData
	p.registers.bufferedData = p.bus.read(p.registers.currentAddr)
	return data
}

func (p *PPU) WritePPUControlPort(value uint8) {
	p.ports.control = newPPUControl(value)
}

func (p *PPU) WritePPUMaskPort(value uint8) {
	p.ports.mask = value
}

func (p *PPU) WriteOAMAddrPort(value uint8) {

}

func (p *PPU) WriteOAMDataPort(value uint8) {

}

func (p *PPU) WritePPUScrollPort(value uint8) {
	if p.registers.writeLatch {
		p.registers.coarseY = uint16(value&0b11111000) >> 3
		p.registers.fineY = uint16(value & 0b00000111)
		p.registers.writeLatch = false
	} else {
		p.registers.coarseX = uint16(value&0b11111000) >> 3
		p.registers.fineX = uint16(value & 0b00000111)
		p.registers.writeLatch = true
	}
}

func (p *PPU) WritePPUAddrPort(value uint8) {
	if p.registers.writeLatch {
		p.registers.currentAddr = (p.registers.currentAddr & 0xFF00) | (uint16(value))
		p.registers.writeLatch = false
	} else {
		p.registers.currentAddr = uint16(value) << 8
		p.registers.writeLatch = true
	}
}

func (p *PPU) WritePPUDataPort(value uint8) {
	p.bus.write(p.registers.currentAddr, value)
	p.registers.currentAddr += p.ports.control.incrementSize
}

func (p *PPU) RunStep() {
	defer p.incrementCycle()

	preRenderScanline := p.renderingState.scanline == 261

	if preRenderScanline {
		if p.renderingState.clock == 1 {
			p.ports.status &= disableStatusVBlank
			p.ports.status &= disableSprite0HitFlag
			p.ports.status &= disableSpriteOverflowFlag
			p.frameDone = false
			p.currentPixel = 0
			p.registers.coarseY = 0
			p.registers.coarseX = 0
		} else if p.renderingState.clock >= 321 && p.renderingState.clock < 337 {
			p.fetchNextTile()
		}
	} else if p.renderingState.scanline < 240 {
		if p.renderingState.clock >= 1 && p.renderingState.clock < 257 {
			nextPixel := p.registers.pixelBuffer.Unbuffer()
			p.appendPixel(nextPixel)
			p.fetchNextTile()
		} else if p.renderingState.clock == 257 {
			if p.registers.fineY == 7 {
				p.registers.fineY = 0
				p.registers.coarseY++
			} else {
				p.registers.fineY++
			}
			p.registers.coarseX = 0
			p.registers.fineX = 0
		} else if p.renderingState.clock >= 321 && p.renderingState.clock < 337 {
			p.fetchNextTile()
		}
	} else if p.renderingState.scanline == 240 && p.renderingState.clock == 0 {
		p.vBlank()
		p.frameDone = true
	} else if p.renderingState.scanline == 241 && p.renderingState.clock == 1 {
		p.ports.status |= enableStatusVBlank
	}

}

func (p *PPU) incrementCycle() {
	p.renderingState.clock++
	if p.renderingState.clock < 341 {
		return
	}

	p.renderingState.clock = 0
	if p.renderingState.scanline == 261 {
		p.renderingState.scanline = 0
	} else {
		p.renderingState.scanline++
	}
}

func (p *PPU) fetchNextTile() {
	switch p.renderingState.fetchingState {
	case stateFetchNametableFirstCycle:
		p.renderingState.currentAddr = p.ports.control.nametableOffset | (p.registers.coarseY << 5) | p.registers.coarseX
		p.renderingState.fetchingState = stateFetchNametableSecondCycle
	case stateFetchNametableSecondCycle:
		p.renderingState.tileIndex = p.bus.read(p.renderingState.currentAddr)
		p.renderingState.fetchingState = stateFetchAttrtableFirstCycle
	case stateFetchAttrtableFirstCycle:
		attrTableX := (p.registers.coarseX & 0b11100) >> 2
		attrTableY := (p.registers.coarseY & 0b11100) >> 2
		p.renderingState.currentAddr = p.ports.control.nametableOffset + 0x3C0 + attrTableX + attrTableY*8
		p.renderingState.fetchingState = stateFetchAttrtableSecondCycle
	case stateFetchAttrtableSecondCycle:
		attrTableByte := p.bus.read(p.renderingState.currentAddr)
		attrTableX := (p.registers.coarseX & 0b10) >> 1
		attrTableY := (p.registers.coarseY & 0b10)
		p.renderingState.paletteId = (attrTableByte >> uint8(attrTableY|attrTableX)) & 0b11
		p.renderingState.fetchingState = stateFetchPatternTableLowPlaneFirstCycle
	case stateFetchPatternTableLowPlaneFirstCycle:
		bitsAddr := (uint16(p.renderingState.tileIndex)*tileSize + uint16(p.registers.fineY)) | p.ports.control.backgroundPatternTableAddr
		p.renderingState.currentAddr = bitsAddr
		p.renderingState.fetchingState = stateFetchPatternTableLowPlaneSecondCycle
	case stateFetchPatternTableLowPlaneSecondCycle:
		p.renderingState.lowBitPlane = p.bus.read(p.renderingState.currentAddr)
		p.renderingState.fetchingState = stateFetchPatternTableHighPlaneFirstCycle
	case stateFetchPatternTableHighPlaneFirstCycle:
		bitsAddr := (uint16(p.renderingState.tileIndex)*tileSize + uint16(p.registers.fineY) + highBitPlaneOffset) | p.ports.control.backgroundPatternTableAddr
		p.renderingState.currentAddr = bitsAddr
		p.renderingState.fetchingState = stateFetchPatternTableHighPlaneSecondCycle
	case stateFetchPatternTableHighPlaneSecondCycle:
		p.renderingState.highBitPlane = p.bus.read(p.renderingState.currentAddr)
		p.updateShiftRegisterWithFetchedData()
		p.registers.coarseX++
		p.renderingState.fetchingState = stateFetchNametableFirstCycle
	}
}

func (p *PPU) vBlank() {
	if p.ports.control.nmiEnabled {
		interrupt.InterruptSignal.Send(interrupt.NonMaskableInterrupt)
	}
}

func (p *PPU) updateShiftRegisterWithFetchedData() {
	for i := range 8 {
		shift := 7 - i
		hi := (p.renderingState.highBitPlane >> shift) & 0b01
		lo := (p.renderingState.lowBitPlane >> shift) & 0b01
		colorIdx := (hi << 1) | lo
		addr := 0x3F00 + uint16(p.renderingState.paletteId*4+colorIdx)

		color := nesPalette[p.bus.read(addr)]
		p.registers.pixelBuffer.Buffer(color)
	}
}

func (p *PPU) appendPixel(color color.RGBA) {
	p.bufferedImage[p.currentPixel] = color
	p.currentPixel++
}

func (p *PPU) FrameDone() bool {
	return p.frameDone
}

func (p *PPU) GetGeneratedImage() []color.RGBA {
	p.frameDone = false
	return p.bufferedImage[:]
}
