package cartridge

type nrom struct {
	BanksQuantity int
	PrgRom        []uint8
	ChrRom        []uint8
}

func (m *nrom) ReadPrg(addr uint16) uint8 {
	romSize := len(m.PrgRom)
	if addr >= uint16(romSize) && m.BanksQuantity == 1 {
		addr -= prgBankSize
	}
	return m.PrgRom[addr]
}

func (m *nrom) WritePrg(_ uint16, _ uint8) {
}

func (m *nrom) ReadChr(addr uint16) uint8 {
	return m.ChrRom[addr]
}

func (m *nrom) WriteChr(_ uint16, _ uint8) {
}
