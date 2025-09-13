package ppu

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
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
	framesCount    uint
	tiles          []Tile
}

func NewPPU(chrRom []uint8) *PPU {
	vram := make([]uint8, 2*kb)
	tiles := generateTilesFromChrRom(chrRom)
	return &PPU{oddFrame: false, vram: vram, writeLatch: false, bus: newBus(chrRom), tiles: tiles}
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

	if p.framesCount == 5 {
		saveImage(p.tiles, p.bus.rom[256:])
	}

	p.framesCount++
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

func saveImage(tiles []Tile, indexes []uint8) {
	var pixelMap = [4]color.RGBA{
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // black
		{R: 0x55, G: 0x55, B: 0xFF, A: 0xFF}, // blue
		{R: 0xFF, G: 0x55, B: 0x55, A: 0xFF}, // red
		{R: 0xFF, G: 0xFF, B: 0xAA, A: 0xFF}, // yellow/cream
	}
	image := image.NewRGBA(image.Rect(0, 0, 256, 128))

	for i, tileIdx := range indexes {
		tile := tiles[tileIdx]
		tileX := (i % 32) * 8
		tileY := (i / 32) * 8

		for y := range 8 {
			for x := range 8 {
				index := tile[y][x]
				image.Set(tileX+x, tileY+y, pixelMap[index])
			}
		}
	}

	out, _ := os.Create("output.png")
	defer out.Close()

	png.Encode(out, image)

	fmt.Printf("Found %d tiles\n", len(tiles))

	os.Exit(0)
}
