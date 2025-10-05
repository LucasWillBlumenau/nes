package ppu

import (
	"fmt"
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

	spriteYPosition uint8 = 0
	spriteTileIndex uint8 = 1
	spriteAttrByte  uint8 = 2
	spriteXPosition uint8 = 3
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
	spriteFetchingStateFetchBitsPlanes
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
	scanline           uint16
	clock              uint16
	tileIndex          uint8
	lowBitPlane        uint8
	highBitPlane       uint8
	paletteId          uint8
	oamEvaluationState oamEvaluationState
	oamData            [4]uint8
	oamSprite          uint8
	currentSpriteIndex int
}

type PPU struct {
	bus                *PPUBus
	bufferedImage      [256 * 240]color.RGBA
	oam                [64][4]uint8
	secondaryOAM       [8][4]uint8
	foreground         [256]color.RGBA
	foregroundPriority [256]bool
	foregroundFilled   [256]bool
	secondaryOAMIndex  int
	rendering          bool
	ports              ppuPorts
	renderingState     ppuRenderingState
	registers          ppuRegisters
	imageChannel       chan []color.RGBA
	currentAddr        vRegister
	tempAddr           uint16
	fineX              uint8
	frameCount         uint64
	cycles             uint64
	cleanVBlank        bool
	oddFrame           bool
}

func NewPPU(bus *PPUBus, imageChannel chan []color.RGBA) *PPU {
	return &PPU{
		bus: bus,
		registers: ppuRegisters{
			writeLatch:   false,
			bufferedData: 0x00,
		},
		imageChannel: imageChannel,
	}
}

func (p *PPU) ReadStatusPort() uint8 {
	currentStatus := p.ports.status
	p.cleanVBlank = true
	p.ports.status &= resetStatusVBlank
	p.registers.writeLatch = false
	return currentStatus | 0x10
}

func (p *PPU) ReadOAMDataPort() uint8 {
	spriteIndex := p.registers.oamAddr >> 2
	spriteByte := p.registers.oamAddr & 0b11
	return p.oam[spriteIndex][spriteByte]
}

func (p *PPU) ReadVRamDataPort() uint8 {
	data := p.registers.bufferedData
	p.registers.bufferedData = p.bus.read(p.currentAddr.Value())
	if p.currentAddr.Value() >= 0x3F00 && p.currentAddr.Value() < 0x4000 {
		data = p.registers.bufferedData
	}
	return data
}

func (p *PPU) WritePPUControlPort(value uint8) {
	p.ports.control = newPPUControl(value)
	p.registers.nametable = p.ports.control.nametable
	p.tempAddr |= uint16(p.ports.control.nametable) << 10

	vBlankStatus := (p.ports.status & setStatusVBlank) > 0
	if p.ports.control.nmiEnabled && vBlankStatus {
		interrupt.InterruptSignal.Send(interrupt.NonMaskableInterrupt)
	}
}

func (p *PPU) WritePPUMaskPort(value uint8) {
	p.ports.mask.value = value
}

func (p *PPU) WriteOAMAddrPort(value uint8) {
	p.registers.oamAddr = value
}

func (p *PPU) WriteOAMDataPort(value uint8) {
	spriteIndex := p.registers.oamAddr >> 2
	spriteByte := p.registers.oamAddr & 0b11
	p.oam[spriteIndex][spriteByte] = value
	p.registers.oamAddr++
}

func (p *PPU) WritePPUScrollPort(value uint8) {
	if p.registers.writeLatch {
		p.tempAddr |= uint16(value&0b11111000) << 2
		p.tempAddr |= uint16(value&0b00000111) << 12
	} else {
		p.tempAddr = uint16(value&0b11111000) >> 3
		p.fineX = value & 0b00000111
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

func (p *PPU) RunSteps(cycles uint16) {
	for range cycles {
		p.runStep()
	}
	p.cleanVBlank = false
}

func (p *PPU) runStep() {
	defer p.incrementCycle()

	preRenderScanline := p.renderingState.scanline == 261

	if preRenderScanline {
		p.handlePreRenderScanline()
	} else if p.renderingState.scanline < 240 {
		p.handleVisibleScanline()
	} else if p.renderingState.scanline == 240 && p.renderingState.clock == 0 {
		p.rendering = false
	} else if p.renderingState.scanline == 241 && p.renderingState.clock == 1 {
		if !p.cleanVBlank {
			p.ports.status |= setStatusVBlank
		}
		p.vBlank()
	}
}

func (p *PPU) handlePreRenderScanline() {
	if p.renderingState.clock == 1 {
		p.ports.status &= resetStatusVBlank
		p.ports.status &= resetSprite0HitFlag
		p.ports.status &= resetSpriteOverflowFlag
		p.frameCount++
		p.imageChannel <- p.bufferedImage[:]
		p.rendering = true
		p.printOAM()
	} else if p.renderingState.clock == 257 && p.ports.mask.RenderingEnabled() {
		p.currentAddr.SetHorizontalBits(p.tempAddr)
	} else if p.renderingState.clock >= 280 && p.renderingState.clock < 305 && p.ports.mask.RenderingEnabled() {
		p.currentAddr.SetVerticalBits(p.tempAddr)
	} else if p.renderingState.clock >= 321 && p.renderingState.clock < 337 {
		p.fetchBackgroundTile()
	}
}

func (p *PPU) printOAM() {
	for _, oam := range p.oam {
		fmt.Printf("%02x %02x %02x %02x ", oam[0], oam[1], oam[2], oam[3])
	}
	fmt.Println()
}

func (p *PPU) handleVisibleScanline() {
	if p.renderingState.clock == 0 {
		p.renderingState.oamSprite = 0
		p.secondaryOAMIndex = 0
		p.renderingState.currentSpriteIndex = 0
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
	} else if p.renderingState.clock < 257 {
		nextPixel := p.registers.pixelBuffer.Unbuffer(p.fineX)
		p.appendPixel(nextPixel)
		p.fetchBackgroundTile()

		if p.renderingState.clock < 65 {
			index := (p.renderingState.clock - 1) >> 1
			spriteIndex := index >> 2
			spriteByte := index & 0b11
			p.secondaryOAM[spriteIndex][spriteByte] = 0xFF
		} else {
			if p.renderingState.clock == 256 && p.ports.mask.RenderingEnabled() {
				p.currentAddr.IncrementY()
			}
			p.evaluateSprite()
		}
	} else if p.renderingState.clock < 321 {
		if p.renderingState.clock == 257 && p.ports.mask.RenderingEnabled() {
			p.currentAddr.SetHorizontalBits(p.tempAddr)
		}
		isDoneFetchingSprites := p.secondaryOAMIndex < p.renderingState.currentSpriteIndex
		if !isDoneFetchingSprites {
			p.fetchSprite()
		}
	} else if p.renderingState.clock < 337 {
		p.fetchBackgroundTile()
	}
}

func (p *PPU) incrementCycle() {
	p.cycles++
	p.renderingState.clock++
	shouldSkipNextDot := p.ports.mask.RenderingEnabled() &&
		p.oddFrame &&
		p.renderingState.scanline == 261 &&
		p.renderingState.clock == 340
	if shouldSkipNextDot {
		p.renderingState.clock++
	}

	if p.renderingState.clock < 341 {
		return

	}

	p.renderingState.clock = 0
	if p.renderingState.scanline == 261 {
		p.renderingState.scanline = 0
		p.oddFrame = !p.oddFrame
	} else {
		p.renderingState.scanline++
	}
}

func (p *PPU) evaluateSprite() {
	switch p.renderingState.oamEvaluationState {
	case evaluateOAMYPositionByteFirstCycle:
		p.renderingState.oamData = p.oam[p.renderingState.oamSprite]
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteSecondCycle
	case evaluateOAMYPositionByteSecondCycle:
		if p.secondaryOAMIndex == len(p.secondaryOAM) {
			p.ports.status |= setSpriteOverflowFlag
			return
		}
		start := p.renderingState.oamData[spriteYPosition]
		end := start + 8
		nextScanline := uint8(p.renderingState.scanline) + 1
		spriteInScanline := nextScanline >= start && nextScanline < end
		if spriteInScanline {
			p.appendSecondaryOAMData(p.renderingState.oamData)
			p.renderingState.oamEvaluationState = evaluateOAMTileIndexFirstCycle
		} else {
			p.renderingState.oamSprite = (p.renderingState.oamSprite + 1) & 0b111111
			p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
		}
	case evaluateOAMTileIndexFirstCycle:
		p.renderingState.oamEvaluationState = evaluateOAMTileIndexSecondCycle
	case evaluateOAMTileIndexSecondCycle:
		p.renderingState.oamEvaluationState = evaluateOAMAttrByteFirstCycle
	case evaluateOAMAttrByteFirstCycle:
		p.renderingState.oamEvaluationState = evaluateOAMAttrByteSecondCycle
	case evaluateOAMAttrByteSecondCycle:
		p.renderingState.oamEvaluationState = evaluateOAMXPositionFirstCycle
	case evaluateOAMXPositionFirstCycle:
		p.renderingState.oamEvaluationState = evaluateOAMXPositionSecondCycle
	case evaluateOAMXPositionSecondCycle:
		p.renderingState.oamSprite = (p.renderingState.oamSprite + 1) & 0b111111
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
	}
}

func (p *PPU) fetchSprite() {
	state := spriteFetchingStateByClock[p.renderingState.clock]
	switch state {
	case spriteFetchingStateFetchBitsPlanes:
		sprite := p.secondaryOAM[p.renderingState.currentSpriteIndex]
		spriteY := sprite[spriteYPosition]
		tileIndex := sprite[spriteTileIndex]
		spriteAttr := sprite[spriteAttrByte]
		oamTileIndex := uint16(tileIndex)
		oamSpriteAttr := newSpriteAttributesFromByte(spriteAttr)
		deltaY := (p.renderingState.scanline + 1) - uint16(spriteY)
		oamSpriteY := deltaY
		if oamSpriteAttr.FlipVertically {
			oamSpriteY = 7 - deltaY
		}
		patternTableIndex := oamTileIndex*tileSize + oamSpriteY
		oamSpriteAddress := p.ports.control.spritePatternTableAddr | patternTableIndex
		oamLowBitPlane := p.bus.read(oamSpriteAddress)
		patternTableIndex = oamTileIndex*tileSize + oamSpriteY + 8
		oamSpriteAddress = p.ports.control.spritePatternTableAddr | patternTableIndex
		oamHighBitPlane := p.bus.read(oamSpriteAddress)
		p.addForegroundPixels(oamHighBitPlane, oamLowBitPlane)
		p.renderingState.currentSpriteIndex++
	}
}

func (p *PPU) addForegroundPixels(highBitPlane uint8, lowBitPlane uint8) {
	sprite := p.secondaryOAM[p.renderingState.currentSpriteIndex]
	attrByte := sprite[spriteAttrByte]
	attr := newSpriteAttributesFromByte(attrByte)
	xPosition := int(sprite[spriteXPosition])
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
		hi := (highBitPlane >> shiftSize) & 0b01
		lo := (lowBitPlane >> shiftSize) & 0b01
		colorIdx := hi<<1 | lo
		p.foreground[index] = p.bus.GetSpriteColor(attr.Palette, colorIdx)
		p.foregroundPriority[index] = attr.HasPriority
		p.foregroundFilled[index] = true
	}
}

func (p *PPU) appendSecondaryOAMData(data [4]uint8) {
	p.secondaryOAM[p.secondaryOAMIndex] = data
	p.secondaryOAMIndex++
}

func (p *PPU) fetchBackgroundTile() {
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
		p.fillShiftRegisters()
		if p.ports.mask.RenderingEnabled() {
			p.currentAddr.IncrementX()
		}
	}
}

func (p *PPU) fillShiftRegisters() {
	if p.ports.mask.RenderingEnabled() {
		for i := range 8 {
			shift := 7 - i
			hi := (p.renderingState.highBitPlane >> shift) & 0b01
			lo := (p.renderingState.lowBitPlane >> shift) & 0b01
			colorIdx := (hi << 1) | lo
			color := p.bus.GetBackgroundColor(p.renderingState.paletteId, colorIdx)
			p.registers.pixelBuffer.Buffer(color)
		}
	}
}

func (p *PPU) appendPixel(pixel color.RGBA) {
	x := p.renderingState.clock - 1
	y := p.renderingState.scanline
	index := y*256 + x

	foregroundFilled := p.foregroundFilled[x]
	spriteHasPriority := p.foregroundPriority[x]
	bgIsTransparent := pixel == p.bus.GetBackdropColor()
	spriteIsTransparent := p.foreground[x] == p.bus.GetBackdropColor()
	spriteRenderingEnabled := p.ports.mask.SpriteRenderingEnabled()

	p.foregroundFilled[x] = false
	p.foregroundPriority[x] = false

	if !p.ports.mask.RenderingEnabled() {
		p.bufferedImage[index] = p.bus.GetBackdropColor()
		return
	}

	drawSprite := spriteRenderingEnabled &&
		!spriteIsTransparent &&
		(bgIsTransparent || spriteHasPriority) &&
		foregroundFilled

	if drawSprite {
		p.bufferedImage[index] = p.foreground[x]
	} else {
		p.bufferedImage[index] = pixel
	}
}

func (p *PPU) vBlank() {
	vBlankStatus := (p.ports.status & setStatusVBlank) > 0
	if p.ports.control.nmiEnabled && vBlankStatus {
		interrupt.InterruptSignal.Send(interrupt.NonMaskableInterrupt)
	}
}

func (p *PPU) Scanline() uint16 {
	return p.renderingState.scanline
}

func (p *PPU) Clock() uint16 {
	return p.renderingState.clock
}

func (p *PPU) Frame() uint64 {
	return p.frameCount
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
		return spriteFetchingStateFetchBitsPlanes
	}
	return spriteFetchingStateIdle
}
