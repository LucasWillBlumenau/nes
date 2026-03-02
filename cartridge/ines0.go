package cartridge

type nrom struct {
	banksQuantity int
	rom           *cartridgeRom
	mirroring     MirroringType
}

func newINES0(rom *cartridgeRom, headers *cartridgeHeaders) mapper {
	return &nrom{
		banksQuantity: headers.ProgramBanksQuantity,
		rom:           rom,
		mirroring:     headers.Mirroring,
	}
}

func (m *nrom) Mirroring() MirroringType {
	return m.mirroring
}

func (m *nrom) ReadPrg(addr uint16) uint8 {
	if addr < 0x8000 {
		return 0
	}
	addr -= 0x8000
	romSize := len(m.rom.Program)
	if addr >= uint16(romSize) && m.banksQuantity == 1 {
		addr -= prgBankSize
	}
	return m.rom.Program[addr]
}

func (m *nrom) WritePrg(_ uint16, _ uint8) {
}

func (m *nrom) ReadChr(addr uint16) uint8 {
	return m.rom.Character.Read(addr)
}

func (m *nrom) WriteChr(addr uint16, data uint8) {
	m.rom.Character.Write(addr, data)
}
