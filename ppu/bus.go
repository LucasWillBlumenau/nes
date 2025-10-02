package ppu

import (
	"image/color"

	"github.com/LucasWillBlumenau/nes/cartridge"
)

var nameTableMirrors = []uint16{
	0: 0,
	1: 0,
	2: 1,
	3: 1,
}

const nameTableOffset uint16 = 0x2000

type PPUBus struct {
	cart              *cartridge.Cartridge
	ram               []uint8
	backgroundPalette []uint8
	foregroundPalette []uint8
}

func NewPPUBus(cart *cartridge.Cartridge) *PPUBus {
	ram := make([]uint8, 2*1024)
	backgroundPalete := make([]uint8, 16)
	foregroundPalete := make([]uint8, 16)

	return &PPUBus{
		cart:              cart,
		ram:               ram,
		backgroundPalette: backgroundPalete,
		foregroundPalette: foregroundPalete,
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

	// TODO: handle horizontal and vertical mirroring, for now it's assumed it is horizontal
	isNameTableAddress := addr < 0x03F00
	if isNameTableAddress {
		nameTableIndex := addr >> 10 & 0b11
		nmTbl := nameTableMirrors[nameTableIndex]
		addr = (nmTbl << 10) | (addr & 0b1111001111111111)
		return &b.ram[addr-nameTableOffset]
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
