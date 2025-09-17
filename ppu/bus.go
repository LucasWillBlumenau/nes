package ppu

var nameTableMirrors = []uint16{
	0: 0,
	1: 0,
	2: 1,
	3: 1,
}

const nameTableOffset uint16 = 0x2000

type bus struct {
	patternTables     []uint8
	ram               []uint8
	backgroundPalette []uint8
	foregroundPalette []uint8
}

func newBus(rom []uint8) *bus {
	ram := make([]uint8, 2*kb)
	backgroundPalete := make([]uint8, 16)
	foregroundPalete := make([]uint8, 16)

	return &bus{
		patternTables:     rom,
		ram:               ram,
		backgroundPalette: backgroundPalete,
		foregroundPalette: foregroundPalete,
	}
}

func (b *bus) write(addr uint16, value uint8) {
	writeAddr := b.getAddress(addr)
	if writeAddr != nil {
		*writeAddr = value
	}

}

func (b *bus) read(addr uint16) uint8 {
	readAddr := b.getAddress(addr)
	if readAddr != nil {
		return *readAddr
	}
	return 0
}

func (b *bus) getAddress(tempAddr uint16) *uint8 {
	addr := tempAddr
	if addr >= 0x4000 {
		addr &= 0x3FFF
	}

	isChrRomAddr := addr < 0x2000
	if isChrRomAddr {
		return &b.patternTables[addr]
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
