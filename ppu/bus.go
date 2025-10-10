package ppu

import (
	"image/color"

	"github.com/LucasWillBlumenau/nes/cartridge"
)

type memoryDevice uint8

const (
	memoryDeviceRom memoryDevice = iota
	memoryDeviceNametable
	memoryDevicePalette
)

var horizontalMirroringOffset = []uint16{
	0: 0,
	1: 0,
	2: 0x400,
	3: 0x400,
}

var verticalMirroringOffset = []uint16{
	0: 0,
	1: 0x400,
	2: 0,
	3: 0x400,
}

type PPUBus struct {
	cart              *cartridge.Cartridge
	ram               []uint8
	backgroundPalette [4][4]uint8
	foregroundPalette [4][4]uint8
}

func NewPPUBus(cart *cartridge.Cartridge) *PPUBus {
	ram := make([]uint8, 2*1024)
	return &PPUBus{
		cart: cart,
		ram:  ram,
	}
}

func (b *PPUBus) write(addr uint16, value uint8) {
	writeAddr, device := b.getAddress(addr)
	if writeAddr != nil {
		if device == memoryDevicePalette {
			value &= 0b00111111
		}
		*writeAddr = value
	}

}

func (b *PPUBus) read(addr uint16) uint8 {
	readAddr, _ := b.getAddress(addr)
	if readAddr != nil {
		return *readAddr
	}
	return 0
}

func (b *PPUBus) getAddress(addr uint16) (*uint8, memoryDevice) {
	if addr >= 0x4000 {
		addr &= 0x3FFF
	}

	isChrRomAddr := addr < 0x2000
	if isChrRomAddr {
		if int(addr) < len(b.cart.CharacterRom) {
			return &b.cart.CharacterRom[addr], memoryDeviceRom
		}
		return nil, memoryDeviceRom
	}

	isNameTableAddress := addr < 0x03F00
	if isNameTableAddress {
		nameTableIndex := addr >> 10 & 0b11
		addr &= 0x3FF
		mirrors := horizontalMirroringOffset
		if b.cart.UseVerticalMirroring {
			mirrors = verticalMirroringOffset
		}
		return &b.ram[addr+mirrors[nameTableIndex]], memoryDeviceNametable
	}

	addr &= 0x1F
	palette := (addr >> 2) & 0b11
	color := addr & 0b11
	isBackgroundColor := addr < 0x10 || color == 0
	if isBackgroundColor {
		return &b.backgroundPalette[palette][color], memoryDevicePalette
	}
	return &b.foregroundPalette[palette][color], memoryDevicePalette
}

func (b *PPUBus) GetBackgroundColor(paletteIndex uint8, colorIndex uint8) color.RGBA {
	if colorIndex == 0 {
		return b.GetBackdropColor()
	}
	color := b.backgroundPalette[paletteIndex][colorIndex]
	return nesPalette[color]
}

func (b *PPUBus) GetSpriteColor(paletteIndex uint8, colorIndex uint8) color.RGBA {
	if colorIndex == 0 {
		return b.GetBackdropColor()
	}
	color := b.foregroundPalette[paletteIndex][colorIndex]
	return nesPalette[color]
}

func (b *PPUBus) GetBackdropColor() color.RGBA {
	color := b.backgroundPalette[0][0]
	return nesPalette[color]
}
