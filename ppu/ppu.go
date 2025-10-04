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
	mask    ppuMask
	status  uint8
}

type ppuRegisters struct {
	bufferedData uint8
	writeLatch   bool
	pixelBuffer  pixelsShiftRegister
	nametable    uint16
	oamAddr      uint8
}

type backgroundFetchingState uint8

const (
	bgFetchingStateIdle backgroundFetchingState = iota
	bgFetchingStateFetchNametable
	bgFetchingStateFetchAttrtable
	bgFetchingStateFetchLowBitplane
	bgFetchingStateFetchHighBitplane
)

var backgroundFetchingStateByClock = loadBgFetchingStates()

type spriteFetchingState uint8

const (
	spriteFetchingStateIdle spriteFetchingState = iota
	spriteFetchingStateFetchLowBitPlane
	spriteFetchingStateFetchHighBitPlane
)

var spriteFetchingStateByClock = loadSpriteFetchingStates()

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

type ppuRenderingState struct {
	scanline                  uint16
	clock                     uint16
	tileIndex                 uint8
	lowBitPlane               uint8
	highBitPlane              uint8
	paletteId                 uint8
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
	bus                *PPUBus
	bufferedImage      [256 * 240]color.RGBA
	oam                [256]uint8
	secondaryOAM       [32]uint8
	foreground         [256]color.RGBA
	foregroundPriority [256]bool
	foregroundFilled   [256]bool
	secondaryOAMIndex  int
	currentPixel       uint32
	rendering          bool
	ports              ppuPorts
	renderingState     ppuRenderingState
	registers          ppuRegisters
	imageChannel       chan []color.RGBA
	currentAddr        vRegister
	tempAddr           uint16
	fineX              uint16
	frameCount         uint64
}

func NewPPU(bus *PPUBus, imageChannel chan []color.RGBA) *PPU {
	return &PPU{
		bus: bus,
		registers: ppuRegisters{
			writeLatch:   false,
			bufferedData: 0x00,
		},
		renderingState: ppuRenderingState{
			scanline: 261,
			clock:    0,
		},
		imageChannel: imageChannel,
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
	p.registers.bufferedData = p.bus.read(p.currentAddr.Value())
	return data
}

func (p *PPU) WritePPUControlPort(value uint8) {
	p.ports.control = newPPUControl(value)
	p.registers.nametable = p.ports.control.nametable
	p.tempAddr |= uint16(p.ports.control.nametable) << 10
}

func (p *PPU) WritePPUMaskPort(value uint8) {
	p.ports.mask.value = value
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
		p.tempAddr |= uint16(value&0b11111000) << 2
		p.tempAddr |= uint16(value&0b00000111) << 12
	} else {
		p.tempAddr = uint16(value&0b11111000) >> 3
		p.fineX = uint16(value & 0b00000111)
	}
	p.registers.writeLatch = !p.registers.writeLatch
}

func (p *PPU) WritePPUAddrPort(value uint8) {
	if p.registers.writeLatch {
		p.tempAddr = (p.tempAddr & 0xFF00) | (uint16(value))
		p.currentAddr.SetValue(p.tempAddr)
	} else {
		p.tempAddr = uint16(value) << 8
	}
	p.registers.writeLatch = !p.registers.writeLatch
}

func (p *PPU) WritePPUDataPort(value uint8) {
	p.bus.write(p.currentAddr.Value(), value)
	p.currentAddr.SetValue(p.currentAddr.Value() + p.ports.control.incrementSize)
}

func (p *PPU) RunStep() {
	defer p.incrementCycle()

	preRenderScanline := p.renderingState.scanline == 261

	if preRenderScanline {
		p.handlePreRenderScanline()
	} else if p.renderingState.scanline < 240 {
		p.handleVisibleScanline()
	} else if p.renderingState.scanline == 240 && p.renderingState.clock == 0 {
		p.frameCount++
		p.imageChannel <- p.bufferedImage[:]
		p.vBlank()
		p.rendering = false
	} else if p.renderingState.scanline == 241 && p.renderingState.clock == 1 {
		p.ports.status |= setStatusVBlank
	}

}

func (p *PPU) handlePreRenderScanline() {
	if p.renderingState.clock == 1 {
		p.ports.status &= resetStatusVBlank
		p.ports.status &= resetSprite0HitFlag
		p.ports.status &= resetSpriteOverflowFlag
		p.rendering = true
		p.currentPixel = 0
	} else if p.renderingState.clock == 257 {
		p.currentAddr.SetHorizontalBits(p.tempAddr)
	} else if p.renderingState.clock >= 280 && p.renderingState.clock < 305 {
		p.currentAddr.SetVerticalBits(p.tempAddr)
	} else if p.renderingState.clock >= 321 && p.renderingState.clock < 337 {
		p.fetchBackgroundTile()
	} else if p.renderingState.clock == 337 {
		p.registers.pixelBuffer.DishMany(p.fineX)
	}
}

func (p *PPU) handleVisibleScanline() {
	if p.renderingState.clock == 0 {
		p.renderingState.oamSprite = 0
		p.renderingState.oamSpriteByte = 0
		p.secondaryOAMIndex = 0
		p.renderingState.currentSpriteBeingFetched = 0
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
	} else if p.renderingState.clock < 257 {
		nextPixel := p.registers.pixelBuffer.Unbuffer()
		p.appendPixel(nextPixel)
		p.fetchBackgroundTile()

		if p.renderingState.clock < 65 {
			index := (p.renderingState.clock - 1) >> 1
			p.secondaryOAM[index] = 0xFF
		} else {
			if p.renderingState.clock == 256 {
				p.currentAddr.IncrementY()
			}
			p.evaluateSprite()
		}
	} else if p.renderingState.clock < 321 {
		if p.renderingState.clock == 257 {
			p.currentAddr.SetHorizontalBits(p.tempAddr)
		}
		isDoneFetchingSprites := p.secondaryOAMIndex < p.renderingState.currentSpriteBeingFetched
		if !isDoneFetchingSprites {
			p.fetchSprite()
		}
	} else if p.renderingState.clock < 337 {
		p.fetchBackgroundTile()
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

func (p *PPU) evaluateSprite() {
	if !p.ports.mask.RenderingEnabled() {
		return
	}

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
	if !p.ports.mask.RenderingEnabled() {
		return
	}

	state := spriteFetchingStateByClock[p.renderingState.clock]
	switch state {
	case spriteFetchingStateFetchLowBitPlane:
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
		patternTableIndex := p.renderingState.oamTileIndex*tileSize + p.renderingState.oamSpriteY
		p.renderingState.oamSpriteAddress = p.ports.control.spritePatternTableAddr | patternTableIndex
		p.renderingState.oamLowBitPlane = p.bus.read(p.renderingState.oamSpriteAddress)
	case spriteFetchingStateFetchHighBitPlane:
		patternTableIndex := p.renderingState.oamTileIndex*tileSize + p.renderingState.oamSpriteY + 8
		p.renderingState.oamSpriteAddress = p.ports.control.spritePatternTableAddr | patternTableIndex
		p.renderingState.oamHighBitPlane = p.bus.read(p.renderingState.oamSpriteAddress)
		p.addForegroundPixels()
		p.renderingState.currentSpriteBeingFetched++
	}
}

func (p *PPU) addForegroundPixels() {
	attrByte := p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4+2]
	attr := newSpriteAttributesFromByte(attrByte)
	xPosition := int(p.secondaryOAM[p.renderingState.currentSpriteBeingFetched*4+3])
	for i := range 8 {
		index := xPosition + i
		isForegroundSpaceFilled := index >= len(p.foreground) ||
			(p.foregroundFilled[index] &&
				p.foreground[index] != p.bus.GetBackdropColor())
		if isForegroundSpaceFilled {
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
		p.foregroundFilled[index] = true
	}
}

func (p *PPU) appendSecondaryOAMData(data uint8) {
	p.secondaryOAM[p.secondaryOAMIndex] = data
	p.secondaryOAMIndex++
}

func (p *PPU) fetchBackgroundTile() {
	if !p.ports.mask.RenderingEnabled() {
		return
	}

	state := backgroundFetchingStateByClock[p.renderingState.clock]
	switch state {
	case bgFetchingStateFetchNametable:
		tileAddr := p.currentAddr.NametableAddress()
		p.renderingState.tileIndex = p.bus.read(tileAddr)
	case bgFetchingStateFetchAttrtable:
		attrTableAddr := p.currentAddr.AttrTableAddress()
		attrTableByte := p.bus.read(attrTableAddr)
		p.renderingState.paletteId = (attrTableByte >> (p.currentAddr.AttrTableBytePart() * 2)) & 0b11
	case bgFetchingStateFetchLowBitplane:
		bitsAddr := (uint16(p.renderingState.tileIndex)*tileSize + p.currentAddr.FineY()) | p.ports.control.backgroundPatternTableAddr
		p.renderingState.lowBitPlane = p.bus.read(bitsAddr)
	case bgFetchingStateFetchHighBitplane:
		bitsAddr := (uint16(p.renderingState.tileIndex)*tileSize + p.currentAddr.FineY() + highBitPlaneOffset) | p.ports.control.backgroundPatternTableAddr
		p.renderingState.highBitPlane = p.bus.read(bitsAddr)
		p.updateShiftRegisterWithFetchedData()
		p.currentAddr.IncrementX()
	}
}

func (p *PPU) updateShiftRegisterWithFetchedData() {
	if !p.ports.mask.BackgroundRenderingEnabled() {
		return
	}

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
	p.foregroundFilled[pixelIndex] = false

	if !p.ports.mask.RenderingEnabled() {
		p.bufferedImage[p.currentPixel] = p.bus.GetBackdropColor()
		p.currentPixel++
		return
	}

	bgIsTransparent := pixel == p.bus.GetBackdropColor()
	spriteIsTransparent := foregroundPixel == p.bus.GetBackdropColor()
	spriteRenderingEnabled := p.ports.mask.SpriteRenderingEnabled()
	spriteHasPriority := p.foregroundPriority[pixelIndex]

	drawSprite := spriteRenderingEnabled &&
		!spriteIsTransparent &&
		(bgIsTransparent || spriteHasPriority)

	if drawSprite {
		p.bufferedImage[p.currentPixel] = foregroundPixel
		p.foreground[pixelIndex] = defaultColor
	} else {
		p.bufferedImage[p.currentPixel] = pixel
	}
	p.currentPixel++
}

func (p *PPU) vBlank() {
	if p.ports.control.nmiEnabled {
		interrupt.InterruptSignal.Send(interrupt.NonMaskableInterrupt)
	}
}

func loadBgFetchingStates() [341]backgroundFetchingState {
	states := [341]backgroundFetchingState{}
	for cycle := range states {
		states[cycle] = determineBackgroundFetchingState(cycle)
	}
	return states
}

func determineBackgroundFetchingState(clock int) backgroundFetchingState {
	shouldFetchTile := (clock >= 1 && clock < 257) ||
		(clock >= 321 && clock < 337)

	if !shouldFetchTile {
		return bgFetchingStateIdle
	}

	switch (clock - 1) % 8 {
	case 1:
		return bgFetchingStateFetchNametable
	case 3:
		return bgFetchingStateFetchAttrtable
	case 5:
		return bgFetchingStateFetchLowBitplane
	case 7:
		return bgFetchingStateFetchHighBitplane
	default:
		return bgFetchingStateIdle
	}
}

func loadSpriteFetchingStates() [341]spriteFetchingState {
	states := [341]spriteFetchingState{}

	for cycle := range states {
		states[cycle] = determineSpriteFetchingState(cycle)
	}
	return states
}

func determineSpriteFetchingState(cycle int) spriteFetchingState {
	shouldFetchTile := cycle >= 257 && cycle < 321
	if !shouldFetchTile {
		return spriteFetchingStateIdle
	}

	switch (cycle - 1) % 8 {
	case 5:
		return spriteFetchingStateFetchLowBitPlane
	case 7:
		return spriteFetchingStateFetchHighBitPlane
	}
	return spriteFetchingStateIdle
}
