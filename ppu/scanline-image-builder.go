package ppu

import (
	"image/color"
)

type state uint8

const (
	stateFetchNametableFirstCycle state = iota
	stateFetchNametableSecondCycle
	stateFetchAttrtableFirstCycle
	stateFetchAttrtableSecondCycle
	stateFetchPatternTableLowPlaneFirstCycle
	stateFetchPatternTableLowPlaneSecondCycle
	stateFetchPatternTableHighPlaneFirstCycle
	stateFetchPatternTableHighPlaneSecondCycle
	stateFetchesDone
)

func (s state) Next() state {
	if s == stateFetchesDone {
		return stateFetchNametableFirstCycle
	}
	return state(s + 1)
}

func (s state) String() string {
	switch s {
	case stateFetchNametableFirstCycle:
		return "stateFetchNametableFirstCycle"
	case stateFetchNametableSecondCycle:
		return "stateFetchNametableSecondCycle"
	case stateFetchAttrtableFirstCycle:
		return "stateFetchAttrtableFirstCycle"
	case stateFetchAttrtableSecondCycle:
		return "stateFetchAttrtableSecondCycle"
	case stateFetchPatternTableLowPlaneFirstCycle:
		return "stateFetchPatternTableLowPlaneFirstCycle"
	case stateFetchPatternTableLowPlaneSecondCycle:
		return "stateFetchPatternTableLowPlaneSecondCycle"
	case stateFetchPatternTableHighPlaneFirstCycle:
		return "stateFetchPatternTableHighPlaneFirstCycle"
	case stateFetchPatternTableHighPlaneSecondCycle:
		return "stateFetchPatternTableHighPlaneSecondCycle"
	case stateFetchesDone:
		return "stateFetchesDone"
	}
	return ""
}

const tileSize uint16 = 16
const highBitPlaneOffset uint16 = 8

var pixelMap = [4]color.RGBA{
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // black
	{R: 0x55, G: 0x55, B: 0xFF, A: 0xFF}, // blue
	{R: 0xFF, G: 0x55, B: 0x55, A: 0xFF}, // red
	{R: 0xFF, G: 0xFF, B: 0xAA, A: 0xFF}, // yellow/cream
}

type scanlineImageBuilder struct {
	bus *bus

	nametable uint16
	state     state
	fineX     uint8
	fineY     uint8

	coarseX  uint16
	coarseY  uint16
	scanline uint16

	// Buffered data
	currentAddr  uint16
	tileIndex    uint8
	lowBitPlane  uint8
	highBitPlane uint8
	paletteId    uint8

	patternTable uint16
}

func newScanlineImageBuilder(bus *bus) *scanlineImageBuilder {
	return &scanlineImageBuilder{bus: bus}
}

func (sb *scanlineImageBuilder) SetNewFrameState(coarseY uint8, fineY uint8, patternTable uint16) {

	sb.scanline = 0
	sb.coarseY = uint16(coarseY)
	sb.coarseX = 0
	sb.patternTable = patternTable
}

func (sb *scanlineImageBuilder) SetNewScanlineState(scanline uint16, fineX uint8) {
	sb.fineX = fineX
	sb.scanline = scanline
	sb.state = stateFetchNametableFirstCycle
	sb.coarseX = 0
	sb.fineY++
	if sb.fineY == 8 {
		sb.fineY = 0
		sb.coarseY++
	}
}

func (sb *scanlineImageBuilder) RunStep(dot uint16) []color.RGBA {
	if sb.scanline >= 240 || dot == 0 || dot > 256 {
		return nil
	}

	switch sb.state {
	case stateFetchNametableFirstCycle:
		sb.currentAddr = (0b1000|sb.nametable)<<10 | (sb.coarseY << 5) | sb.coarseX
	case stateFetchNametableSecondCycle:
		sb.tileIndex = sb.bus.read(sb.currentAddr)
	case stateFetchAttrtableFirstCycle:
		attrTableX := (dot - 1) / 32
		attrTableY := sb.scanline / 32
		sb.currentAddr = (0b1000|sb.nametable)<<10 + 0x3C0 + attrTableX + attrTableY*8

	case stateFetchAttrtableSecondCycle:
		attrTableByte := sb.bus.read(sb.currentAddr)
		tileX := (dot - 1) / 8
		tileY := sb.scanline / 8
		attrTableByteX := tileX % 2
		attrTableByteY := tileY % 2
		sb.paletteId = (attrTableByte >> uint8(attrTableByteY<<1|attrTableByteX)) & 0b11
	case stateFetchPatternTableLowPlaneFirstCycle:
		sb.currentAddr = (uint16(sb.tileIndex)*tileSize + uint16(sb.fineY)) | sb.patternTable
	case stateFetchPatternTableLowPlaneSecondCycle:
		sb.lowBitPlane = sb.bus.read(sb.currentAddr)
	case stateFetchPatternTableHighPlaneFirstCycle:
		sb.currentAddr = (uint16(sb.tileIndex)*tileSize + uint16(sb.fineY) + highBitPlaneOffset) | sb.patternTable
	case stateFetchPatternTableHighPlaneSecondCycle:
		sb.highBitPlane = sb.bus.read(sb.currentAddr)
	}

	sb.state = sb.state.Next()
	if sb.state != stateFetchesDone {
		return nil
	}

	var colors []color.RGBA
	if dot == 239 && sb.fineX > 0 {
		colors = make([]color.RGBA, 8+sb.fineX)
		addr := (0b1000|sb.nametable)<<10 | (sb.coarseY << 5) | sb.coarseX
		tileIndex := sb.bus.read(addr)
		lowBitPlane := sb.bus.read(uint16(tileIndex)*tileSize + uint16(sb.fineY))
		highBitPlane := sb.bus.read(uint16(tileIndex)*tileSize + uint16(sb.fineY) + highBitPlaneOffset)

		for i := range sb.fineX {
			shift := 7 - i
			hi := (highBitPlane >> shift) & 0b01
			lo := (lowBitPlane >> shift) & 0b01
			colors[8+i] = pixelMap[(hi<<1)|lo]
		}
	} else {
		colors = make([]color.RGBA, 8)
	}

	for i := range colors {
		shift := 7 - i
		hi := (sb.highBitPlane >> shift) & 0b01
		lo := (sb.lowBitPlane >> shift) & 0b01
		colorIdx := (hi << 1) | lo
		addr := 0x3F00 + uint16(sb.paletteId*4+colorIdx)

		colors[i] = nesPalette[sb.bus.read(addr)]
	}

	sb.state = sb.state.Next()
	sb.coarseX++
	return colors
}
