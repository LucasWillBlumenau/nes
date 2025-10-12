package ppu

import (
	"image"
	"image/color"

	"github.com/LucasWillBlumenau/nes/interrupt"
)

const (
	tileSize           uint16 = 16
	highBitPlaneOffset uint16 = 8

	setStatusVBlank         uint8 = 0b10000000
	setSprite0HitFlag       uint8 = 0b01000000
	setSpriteOverflowFlag   uint8 = 0b00100000
	resetStatusVBlank       uint8 = 0b01111111
	resetSprite0HitFlag     uint8 = 0b10111111
	resetSpriteOverflowFlag uint8 = 0b11011111

	spriteYPosition uint8 = 0
	spriteTileIndex uint8 = 1
	spriteAttrByte  uint8 = 2
	spriteXPosition uint8 = 3

	originalWidth  = 256
	originalHeigth = 240
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

type foreground struct {
	pixels   [256]pixel
	priority [256]bool
	filled   [256]bool
	tileIds  [256]uint8
}

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
	bus               *PPUBus
	currentFrame      image.RGBA
	oam               [64][4]uint8
	secondaryOAM      [8][4]uint8
	secondaryTileIds  [8]uint8
	foreground        foreground
	nextForeground    foreground
	secondaryOAMIndex int
	rendering         bool
	ports             ppuPorts
	renderingState    ppuRenderingState
	registers         ppuRegisters
	frameChannel      chan image.RGBA
	currentAddr       vRegister
	tempAddr          uint16
	fineX             uint8
	frameCount        uint64
	cycles            uint64
	scaleFactor       int
	cleanVBlank       bool
	oddFrame          bool
}

func NewPPU(bus *PPUBus, frameChannel chan image.RGBA, scaleFactor int) *PPU {
	ppu := &PPU{
		bus:          bus,
		frameChannel: frameChannel,
		scaleFactor:  scaleFactor,
		currentFrame: *image.NewRGBA(image.Rect(0, 0, originalWidth*scaleFactor, originalHeigth*scaleFactor)),
	}
	return ppu
}

func (p *PPU) ReadStatusPort() uint8 {
	currentStatus := p.ports.status
	p.cleanVBlank = true
	p.ports.status &= resetStatusVBlank
	p.registers.writeLatch = false
	return currentStatus | (p.registers.bufferedData & 0b00011111)
}

func (p *PPU) ReadOAMDataPort() uint8 {
	spriteIndex := p.registers.oamAddr >> 2
	spriteByte := p.registers.oamAddr & 0b11
	return p.oam[spriteIndex][spriteByte]
}

func (p *PPU) ReadVRamDataPort() uint8 {
	data := p.registers.bufferedData
	addr := p.currentAddr.Value
	p.registers.bufferedData = p.bus.read(addr)
	mirroredAddr := addr & 0x3FFF
	currentAddrPointsToPaletteData := mirroredAddr >= 0x3F00 && mirroredAddr < 0x4000
	if currentAddrPointsToPaletteData {
		data = p.registers.bufferedData
	}
	p.currentAddr.Value += p.ports.control.incrementSize
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
		p.currentAddr.Value = p.tempAddr
	} else {
		p.tempAddr = uint16(value) << 8
	}
	p.registers.writeLatch = !p.registers.writeLatch
}

func (p *PPU) WritePPUDataPort(value uint8) {
	p.bus.write(p.currentAddr.Value, value)
	p.currentAddr.Value += p.ports.control.incrementSize
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
	} else if p.renderingState.scanline == 241 && p.renderingState.clock == 1 {
		if !p.cleanVBlank {
			p.ports.status |= setStatusVBlank
		}
		p.vBlank()
		p.rendering = false
	}
}

func (p *PPU) handlePreRenderScanline() {
	if p.renderingState.clock == 1 {
		p.ports.status &= resetStatusVBlank
		p.ports.status &= resetSprite0HitFlag
		p.ports.status &= resetSpriteOverflowFlag
		p.frameCount++
		p.frameChannel <- p.currentFrame
		p.rendering = true
	} else if p.renderingState.clock == 257 && p.ports.mask.RenderingEnabled() {
		p.currentAddr.SetHorizontalBits(p.tempAddr)
	} else if p.renderingState.clock >= 280 && p.renderingState.clock < 305 && p.ports.mask.RenderingEnabled() {
		p.currentAddr.SetVerticalBits(p.tempAddr)
	} else if p.renderingState.clock >= 321 && p.renderingState.clock < 337 {
		p.fetchBackgroundTile()
	}
}

func (p *PPU) handleVisibleScanline() {
	if p.renderingState.clock == 0 {
		p.renderingState.oamSprite = 0
		p.secondaryOAMIndex = 0
		p.renderingState.currentSpriteIndex = 0
		p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
		p.foreground = p.nextForeground
		p.nextForeground = foreground{}
	} else if p.renderingState.clock < 257 {
		p.appendPixel()
		p.fetchBackgroundTile()

		if p.renderingState.clock < 65 {
			index := (p.renderingState.clock - 1) >> 1
			spriteIndex := index >> 2
			spriteByte := index & 0b11
			p.secondaryOAM[spriteIndex][spriteByte] = 0xFF
		} else {
			if p.ports.mask.RenderingEnabled() {
				if p.renderingState.clock == 256 {
					p.currentAddr.IncrementY()
				}
				p.evaluateSprite()
			}
		}
	} else if p.renderingState.clock < 321 {
		if p.renderingState.clock == 257 && p.ports.mask.RenderingEnabled() {
			p.currentAddr.SetHorizontalBits(p.tempAddr)
		}
		isDoneFetchingSprites := p.secondaryOAMIndex < p.renderingState.currentSpriteIndex
		if !isDoneFetchingSprites && p.ports.mask.RenderingEnabled() {
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
		start := p.renderingState.oamData[spriteYPosition]
		end := start + 8
		nextScanline := uint8(p.renderingState.scanline) + 1
		spriteInScanline := nextScanline >= start && nextScanline < end

		if !spriteInScanline {
			p.renderingState.oamSprite = (p.renderingState.oamSprite + 1) & 0b111111
			p.renderingState.oamEvaluationState = evaluateOAMYPositionByteFirstCycle
			return
		}

		p.ports.status |= setSpriteOverflowFlag
		if p.secondaryOAMIndex < len(p.secondaryOAM) {
			p.appendSecondaryOAMData(p.renderingState.oamData)
			p.renderingState.oamEvaluationState = evaluateOAMTileIndexFirstCycle
		} else {
			p.ports.status |= setSpriteOverflowFlag
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
	}
}

func (p *PPU) addForegroundPixels(highBitPlane uint8, lowBitPlane uint8) {
	sprite := p.secondaryOAM[p.renderingState.currentSpriteIndex]
	tileId := p.secondaryTileIds[p.renderingState.currentSpriteIndex]

	attrByte := sprite[spriteAttrByte]
	attr := newSpriteAttributesFromByte(attrByte)
	xPosition := int(sprite[spriteXPosition])
	for i := range 8 {
		index := xPosition + i
		isForegroundSpaceFilled := index >= len(p.nextForeground.pixels) ||
			(p.nextForeground.filled[index] &&
				p.nextForeground.pixels[index].Color > 0)
		if isForegroundSpaceFilled {
			break
		}
		shiftSize := i
		if !attr.FlipHorizontally {
			shiftSize = 7 - i
		}
		hi := (highBitPlane >> shiftSize) & 0b01
		lo := (lowBitPlane >> shiftSize) & 0b01
		color := hi<<1 | lo
		p.nextForeground.pixels[index] = pixel{Palette: attr.Palette, Color: color}
		p.nextForeground.priority[index] = attr.HasPriority
		p.nextForeground.filled[index] = true
		p.nextForeground.tileIds[index] = tileId
	}
	p.renderingState.currentSpriteIndex++
}

func (p *PPU) appendSecondaryOAMData(data [4]uint8) {
	p.secondaryOAM[p.secondaryOAMIndex] = data
	p.secondaryTileIds[p.secondaryOAMIndex] = p.renderingState.oamSprite + 1
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
		shift := (p.currentAddr.AttrTableBytePart()) << 1
		p.renderingState.paletteId = (attrTableByte >> shift) & 0b11
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
			color := (hi << 1) | lo
			pixel := pixel{Palette: p.renderingState.paletteId, Color: color}
			p.registers.pixelBuffer.Buffer(pixel)
		}
	}
}

func (p *PPU) appendPixel() {
	x := int(p.renderingState.clock - 1)
	y := int(p.renderingState.scanline)

	pixelColor := p.getCurrentPixelColor(x)
	for i := range p.scaleFactor {
		x := x*p.scaleFactor + i
		y := y*p.scaleFactor + i
		p.currentFrame.Set(x, y, pixelColor)
	}
}

func (p *PPU) getCurrentPixelColor(x int) color.RGBA {
	foregroundFilled := p.foreground.filled[x]
	spriteHasPriority := p.foreground.priority[x]
	tileId := p.foreground.tileIds[x]
	fgPixel := p.foreground.pixels[x]
	bgPixel := p.registers.pixelBuffer.Unbuffer(p.fineX)
	bgIsTransparent := bgPixel.Color == 0
	spriteIsTransparent := fgPixel.Color == 0

	if !p.ports.mask.RenderingEnabled() {
		return p.bus.GetBackdropColor()
	}

	if !spriteIsTransparent && !bgIsTransparent && tileId == 1 {
		p.ports.status |= setSprite0HitFlag
	}

	if spriteIsTransparent && bgIsTransparent {
		return p.bus.GetBackdropColor()
	} else if !bgIsTransparent &&
		(spriteIsTransparent ||
			!spriteHasPriority ||
			!foregroundFilled ||
			!p.ports.mask.SpriteRenderingEnabled()) {
		return p.bus.GetBackgroundColor(bgPixel.Palette, bgPixel.Color)
	}

	return p.bus.GetSpriteColor(fgPixel.Palette, fgPixel.Color)
}

func (p *PPU) vBlank() {
	vBlankStatus := (p.ports.status & setStatusVBlank) > 0
	if p.ports.control.nmiEnabled && vBlankStatus {
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
		return spriteFetchingStateFetchBitsPlanes
	}
	return spriteFetchingStateIdle
}
