package ppu

import (
	"image/color"

	"github.com/LucasWillBlumenau/nes/interrupt"
)

const (
	tileSize           uint16 = 16
	highBitPlaneOffset uint16 = 8

	setStatusVBlank         uint8 = 0b10000000
	setSpriteOverflowFlag   uint8 = 0b00100000
	resetStatusVBlank       uint8 = 0b01111111
	resetSprite0HitFlag     uint8 = 0b10111111
	resetSpriteOverflowFlag uint8 = 0b11011111
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
	nametable    uint16
	oamAddr      uint8
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
)

type oamEvaluationState uint8

const (
	evaluateOAMYPositionByteFirstCycle oamEvaluationState = iota
	evaluateOAMYPositionByteSecondCycle
	evaluateOAMTileIndexFirstCycle
	evaluateOAMTileIndexSecondCycle
	evaluateOAMAttrByteFirstCycle
	evaluateOAMAttrByteSecondCycle
	evaluateOAMXPositionFirstCycle
	evaluateOAMXPositionSecondCycle
)

type oamFetchingState uint8

const (
	fetchOAMGarbageReadFirstCycle oamFetchingState = iota
	fetchOAMGarbageReadSecondCycle
	fetchOAMGarbageRead2FirstCycle
	fetchOAMGarbageRead2SecondCycle
	fetchOAMLowBitPlaneFirstCycle
	fetchOAMLowBitPlaneSecondCycle
	fetchOAMHighBitPlaneFirstCycle
	fetchOAMHighBitPlaneSecondCycle
)

type ppuRenderingState struct {
	scanline                  uint16
	clock                     uint16
	fetchingState             pixelFetchingState
	tileIndex                 uint8
	lowBitPlane               uint8
	highBitPlane              uint8
	paletteId                 uint8
	currentAddr               uint16
	coarseX                   uint16
	coarseY                   uint16
	fineX                     uint16
	fineY                     uint16
	oamFetchingState          oamFetchingState
	oamEvaluationState        oamEvaluationState
	oamData                   uint8
	oamSprite                 uint8
	oamSpriteByte             uint8
	oamSpriteAddress          uint16
	oamLowBitPlane            uint8
	oamHighBitPlane           uint8
	oamSpriteY                uint16
	oamTileIndex              uint16
	oamSpriteAttr             spriteAttributes
	currentSpriteBeingFetched int
}

type PPU struct {
	bus                *bus
	bufferedImage      [256 * 240]color.RGBA
	oam                [256]uint8
	secondaryOAM       [32]uint8
	foreground         [256]color.RGBA
	foregroundPriority [256]bool
	secondaryOAMIndex  int
	currentPixel       uint32
	frameDone          bool
	ports              ppuPorts
	renderingState     ppuRenderingState
	registers          ppuRegisters
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
	p.ports.status &= setStatusVBlank ^ 0xFF
	p.registers.writeLatch = false

	return currentStatus
}

func (p *PPU) ReadOAMDataPort() uint8 {
	return p.oam[p.registers.oamAddr]
}

func (p *PPU) ReadVRamDataPort() uint8 {
	data := p.registers.bufferedData
	p.registers.bufferedData = p.bus.read(p.registers.currentAddr)
	return data
}

func (p *PPU) WritePPUControlPort(value uint8) {
	p.ports.control = newPPUControl(value)
	p.registers.nametable = p.ports.control.nametable
}

func (p *PPU) WritePPUMaskPort(value uint8) {
	p.ports.mask = value
}

func (p *PPU) WriteOAMAddrPort(value uint8) {
	p.registers.oamAddr = value
}

func (p *PPU) WriteOAMDataPort(value uint8) {
	p.oam[p.registers.oamAddr] = value
	p.registers.oamAddr++
}

func (p *PPU) WritePPUScrollPort(value uint8) {
	if p.registers.writeLatch {
		p.registers.coarseY = uint16(value&0b11111000) >> 3
		p.registers.fineY = uint16(value & 0b00000111)
	} else {
		p.registers.coarseX = uint16(value&0b11111000) >> 3
		p.registers.fineX = uint16(value & 0b00000111)
	}
	p.registers.writeLatch = !p.registers.writeLatch
}

func (p *PPU) WritePPUAddrPort(value uint8) {
	if p.registers.writeLatch {
		p.registers.currentAddr = (p.registers.currentAddr & 0xFF00) | (uint16(value))
	} else {
		p.registers.currentAddr = uint16(value) << 8
	}
	p.registers.writeLatch = !p.registers.writeLatch
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
			p.ports.status &= resetStatusVBlank
			p.ports.status &= resetSprite0HitFlag
			p.ports.status &= resetSpriteOverflowFlag
			p.frameDone = false
			p.currentPixel = 0
			p.renderingState.coarseY = p.registers.coarseY
			p.renderingState.coarseX = p.registers.coarseX
			p.renderingState.fineY = p.registers.fineY
			p.renderingState.fineX = p.registers.fineX
		} else if p.renderingState.clock >= 321 && p.renderingState.clock < 337 {
			p.fetchNextTile()
		} else if p.renderingState.clock == 337 {
			for range p.renderingState.fineX {
				p.registers.pixelBuffer.Unbuffer()
			}
		}
	} else if p.renderingState.scanline < 240 {
		if p.renderingState.clock == 0 {
			p.renderingState.oamSprite = 0
			p.renderingState.oamSpriteByte = 0
			p.secondaryOAMIndex = 0
			p.renderingState.currentSpriteBeingFetched = 0
			p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
		} else if p.renderingState.clock < 257 {
			nextPixel := p.registers.pixelBuffer.Unbuffer()
			p.appendPixel(nextPixel)
			p.fetchNextTile()

			if p.renderingState.clock < 65 {
				index := (p.renderingState.clock - 1) >> 1
				p.secondaryOAM[index] = 0xFF
			} else {
				p.evaluateSprite()
			}
		} else if p.renderingState.clock < 321 {
			if p.renderingState.clock == 257 {
				p.moveToNextScanline()
			}
			isDoneFetchingSprites := p.secondaryOAMIndex < p.renderingState.currentSpriteBeingFetched
			if !isDoneFetchingSprites {
				p.fetchSprite()
			}
		} else if p.renderingState.clock < 337 {
			p.fetchNextTile()
		}
	} else if p.renderingState.scanline == 240 && p.renderingState.clock == 0 {
		p.vBlank()
		p.frameDone = true
	} else if p.renderingState.scanline == 241 && p.renderingState.clock == 1 {
		p.ports.status |= setStatusVBlank
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

func (p *PPU) moveToNextScanline() {
	if p.registers.fineY == 7 {
		p.registers.fineY = 0
		if p.renderingState.coarseY == 31 {
			p.renderingState.coarseY = 0
			p.registers.nametable ^= 0b10
		} else {
			p.renderingState.coarseY++
		}
	} else {
		p.registers.fineY++
	}
	p.renderingState.coarseX = p.registers.coarseX
	p.renderingState.fineX = p.registers.fineX
}

func (p *PPU) evaluateSprite() {
	switch p.renderingState.oamEvaluationState {
	case evaluateOAMYPositionByteFirstCycle:
		index := p.renderingState.oamSprite*4 + p.renderingState.oamSpriteByte
		p.renderingState.oamData = p.oam[index]
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteSecondCycle
	case evaluateOAMYPositionByteSecondCycle:
		if p.secondaryOAMIndex == 32 {
			p.ports.status |= setSpriteOverflowFlag
			return
		}
		start := p.renderingState.oamData
		end := start + 8
		nextScanline := uint8(p.renderingState.scanline) + 1
		spriteInScanline := nextScanline >= start && nextScanline < end
		if spriteInScanline {
			p.appendSecondaryOAMData(p.renderingState.oamData)
			p.renderingState.oamSpriteByte = (p.renderingState.oamSpriteByte + 1) & 0b11
			p.renderingState.oamEvaluationState = evaluateOAMTileIndexFirstCycle
		} else {
			p.renderingState.oamSprite = (p.renderingState.oamSprite + 1) & 0b111111
			p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
		}
	case evaluateOAMTileIndexFirstCycle:
		index := p.renderingState.oamSprite*4 + p.renderingState.oamSpriteByte
		p.renderingState.oamData = p.oam[index]
		p.renderingState.oamEvaluationState = evaluateOAMTileIndexSecondCycle
	case evaluateOAMTileIndexSecondCycle:
		p.appendSecondaryOAMData(p.renderingState.oamData)
		p.renderingState.oamSpriteByte = (p.renderingState.oamSpriteByte + 1) & 0b11
		p.renderingState.oamEvaluationState = evaluateOAMAttrByteFirstCycle
	case evaluateOAMAttrByteFirstCycle:
		index := p.renderingState.oamSprite*4 + p.renderingState.oamSpriteByte
		p.renderingState.oamData = p.oam[index]
		p.renderingState.oamEvaluationState = evaluateOAMAttrByteSecondCycle
	case evaluateOAMAttrByteSecondCycle:
		p.appendSecondaryOAMData(p.renderingState.oamData)
		p.renderingState.oamSpriteByte = (p.renderingState.oamSpriteByte + 1) & 0b11
		p.renderingState.oamEvaluationState = evaluateOAMXPositionFirstCycle
	case evaluateOAMXPositionFirstCycle:
		index := p.renderingState.oamSprite*4 + p.renderingState.oamSpriteByte
		p.renderingState.oamData = p.oam[index]
		p.renderingState.oamEvaluationState = evaluateOAMXPositionSecondCycle
	case evaluateOAMXPositionSecondCycle:
		p.appendSecondaryOAMData(p.renderingState.oamData)
		p.renderingState.oamSpriteByte = (p.renderingState.oamSpriteByte + 1) & 0b11
		p.renderingState.oamSprite = (p.renderingState.oamSprite + 1) & 0b111111
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
	}
}

func (p *PPU) fetchSprite() {
	switch p.renderingState.oamFetchingState {
	case fetchOAMGarbageReadFirstCycle:
		p.renderingState.oamFetchingState = fetchOAMGarbageReadSecondCycle
	case fetchOAMGarbageReadSecondCycle:
		p.renderingState.oamFetchingState = fetchOAMGarbageRead2FirstCycle
	case fetchOAMGarbageRead2FirstCycle:
		p.renderingState.oamFetchingState = fetchOAMGarbageRead2SecondCycle
	case fetchOAMGarbageRead2SecondCycle:
		spriteY := p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4]
		tileIndex := p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4+1]
		spriteAttr := p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4+2]
		p.renderingState.oamTileIndex = uint16(tileIndex)
		p.renderingState.oamSpriteAttr = newSpriteAttributesFromByte(spriteAttr)

		deltaY := (p.renderingState.scanline + 1) - uint16(spriteY)
		p.renderingState.oamSpriteY = deltaY
		if p.renderingState.oamSpriteAttr.FlipVertically {
			p.renderingState.oamSpriteY = 7 - deltaY
		}
		p.renderingState.oamFetchingState = fetchOAMLowBitPlaneFirstCycle
	case fetchOAMLowBitPlaneFirstCycle:
		patternTableIndex := p.renderingState.oamTileIndex*tileSize + p.renderingState.oamSpriteY
		p.renderingState.oamSpriteAddress = p.ports.control.spritePatternTableAddr | patternTableIndex
		p.renderingState.oamFetchingState = fetchOAMLowBitPlaneSecondCycle
	case fetchOAMLowBitPlaneSecondCycle:
		p.renderingState.oamLowBitPlane = p.bus.read(p.renderingState.oamSpriteAddress)
		p.renderingState.oamFetchingState = fetchOAMHighBitPlaneFirstCycle
	case fetchOAMHighBitPlaneFirstCycle:
		patternTableIndex := p.renderingState.oamTileIndex*tileSize + p.renderingState.oamSpriteY + 8
		p.renderingState.oamSpriteAddress = p.ports.control.spritePatternTableAddr | patternTableIndex
		p.renderingState.oamFetchingState = fetchOAMHighBitPlaneSecondCycle
	case fetchOAMHighBitPlaneSecondCycle:
		p.renderingState.oamHighBitPlane = p.bus.read(p.renderingState.oamSpriteAddress)
		p.addForegroundPixels()
		p.renderingState.currentSpriteBeingFetched++
		p.renderingState.oamFetchingState = fetchOAMGarbageReadFirstCycle
	}
}

func (p *PPU) addForegroundPixels() {
	attrByte := p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4+2]
	attr := newSpriteAttributesFromByte(attrByte)
	xPosition := int(p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4+3])
	for i := range 8 {
		index := xPosition + i
		if index >= len(p.foreground) {
			break
		}
		shiftSize := i
		if !attr.FlipHorizontally {
			shiftSize = 7 - i
		}
		hi := (p.renderingState.oamHighBitPlane >> shiftSize) & 0b01
		lo := (p.renderingState.oamLowBitPlane >> shiftSize) & 0b01
		colorIdx := hi<<1 | lo
		p.foreground[index] = p.bus.GetSpriteColor(attr.Palette, colorIdx)
		p.foregroundPriority[index] = attr.HasPriority
	}
}

func (p *PPU) fetchNextTile() {
	switch p.renderingState.fetchingState {
	case stateFetchNametableFirstCycle:
		p.renderingState.currentAddr = (0b1000|p.registers.nametable)<<10 | (p.renderingState.coarseY << 5) | p.renderingState.coarseX
		p.renderingState.fetchingState = stateFetchNametableSecondCycle
	case stateFetchNametableSecondCycle:
		p.renderingState.tileIndex = p.bus.read(p.renderingState.currentAddr)
		p.renderingState.fetchingState = stateFetchAttrtableFirstCycle
	case stateFetchAttrtableFirstCycle:
		attrTableX := (p.renderingState.coarseX & 0b11100) >> 2
		attrTableY := (p.renderingState.coarseY & 0b11100) >> 2
		p.renderingState.currentAddr = (0b1000|p.registers.nametable)<<10 + 0x3C0 + attrTableX + attrTableY*8
		p.renderingState.fetchingState = stateFetchAttrtableSecondCycle
	case stateFetchAttrtableSecondCycle:
		attrTableByte := p.bus.read(p.renderingState.currentAddr)
		attrTableX := (p.renderingState.coarseX & 0b10) >> 1
		attrTableY := (p.renderingState.coarseY & 0b10)
		p.renderingState.paletteId = (attrTableByte >> (uint8(attrTableY|attrTableX) * 2)) & 0b11
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
		if p.renderingState.coarseX == 31 {
			p.renderingState.coarseX = 0
			p.registers.nametable ^= 0b01
		} else {
			p.renderingState.coarseX++
		}
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
		color := p.bus.GetBackgroundColor(p.renderingState.paletteId, colorIdx)
		p.registers.pixelBuffer.Buffer(color)
	}
}

func (p *PPU) appendPixel(pixel color.RGBA) {
	var defaultColor color.RGBA
	pixelIndex := p.renderingState.clock - 1
	foregroundPixel := p.foreground[pixelIndex]

	bgIsTransparent := pixel == p.bus.GetBackgroundTransparencyColor()
	spriteIsTransparent := foregroundPixel == p.bus.GetSpriteTransparencyColor()
	spriteHasPriority := p.foregroundPriority[pixelIndex]

	if !spriteIsTransparent && (bgIsTransparent || spriteHasPriority) {
		p.bufferedImage[p.currentPixel] = foregroundPixel
		p.foreground[pixelIndex] = defaultColor
	} else {
		p.bufferedImage[p.currentPixel] = pixel
	}
	p.currentPixel++
}

func (p *PPU) appendSecondaryOAMData(data uint8) {
	p.secondaryOAM[p.secondaryOAMIndex] = data
	p.secondaryOAMIndex++
}

func (p *PPU) FrameDone() bool {
	return p.frameDone
}

func (p *PPU) GetGeneratedImage() []color.RGBA {
	p.frameDone = false
	return p.bufferedImage[:]
}
