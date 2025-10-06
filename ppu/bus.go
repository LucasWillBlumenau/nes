package ppu

import (
	"image/color"

	"github.com/LucasWillBlumenau/nes/cartridge"
)

var horizontalMirrors = []uint16{
	0: 0,
	1: 0,
	2: 0x400,
	3: 0x400,
}

var verticalMirrors = []uint16{
	0: 0,
	1: 0x400,
	2: 0,
	3: 0x400,
}

type PPUBus struct {
	cart              *cartridge.Cartridge
	ram               []uint8
	backgroundPalette [16]uint8
	foregroundPalette [16]uint8
}

func NewPPUBus(cart *cartridge.Cartridge) *PPUBus {
	ram := make([]uint8, 2*1024)
	return &PPUBus{
		cart: cart,
		ram:  ram,
	}
}

func (b *PPUBus) write(addr uint16, value uint8) {
	writeAddr := b.getAddress(addr)
	if writeAddr != nil {
		*writeAddr = value
	}

}

func (b *PPUBus) read(addr uint16) uint8 {
	readAddr := b.getAddress(addr)
	if readAddr != nil {
		return *readAddr
	}
	return 0
}

func (b *PPUBus) getAddress(addr uint16) *uint8 {
	if addr >= 0x4000 {
		addr &= 0x3FFF
	}

	isChrRomAddr := addr < 0x2000
	if isChrRomAddr {
		if int(addr) < len(b.cart.CharacterRom) {
			return &b.cart.CharacterRom[addr]
		}
		return nil
	}

	isNameTableAddress := addr < 0x03F00
	if isNameTableAddress {
		nameTableIndex := addr >> 10 & 0b11
		addr &= 0x3FF
		mirrors := horizontalMirrors
		if b.cart.UseVerticalMirroring {
			mirrors = verticalMirrors
		}
		return &b.ram[addr+mirrors[nameTableIndex]]
	}

	addr &= 0x1F
	isBackgroundColor := addr < 0x10
	if isBackgroundColor {
		return &b.backgroundPalette[addr]
	} else {
		return &b.foregroundPalette[addr-0x10]
	}
}

func (b *PPUBus) GetBackgroundColor(paletteIndex uint8, colorIndex uint8) color.RGBA {
	if colorIndex == 0 {
		return b.GetBackdropColor()
	}
	color := b.backgroundPalette[paletteIndex*4+colorIndex]
	return nesPalette[color]
}

func (b *PPUBus) GetSpriteColor(paletteIndex uint8, colorIndex uint8) color.RGBA {
	if colorIndex == 0 {
		return b.GetBackdropColor()
	}
	color := b.foregroundPalette[paletteIndex*4+colorIndex]
	return nesPalette[color]
}

func (b *PPUBus) GetBackdropColor() color.RGBA {
	color := b.backgroundPalette[0]
	return nesPalette[color]
}
