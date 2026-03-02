package cartridge

type ines2 struct {
	mirroring    MirroringType
	selectedBank int
	rom          *cartridgeRom
	headers      *cartridgeHeaders
}

func newINES2(rom *cartridgeRom, headers *cartridgeHeaders) mapper {
	return &ines2{
		mirroring:    headers.Mirroring,
		selectedBank: 0,
		rom:          rom,
		headers:      headers,
	}
}

func (m *ines2) Mirroring() MirroringType {
	return m.mirroring
}

func (m *ines2) ReadPrg(addr16 uint16) uint8 {
	addr := int(addr16) - 0x8000
	if addr < 0x4000 {
		addr += 16 * 1024 * m.selectedBank
		return m.rom.Program[addr]
	}
	addr += 16*1024*(m.headers.ProgramBanksQuantity-1) - 0x4000
	return m.rom.Program[addr]
}

func (m *ines2) WritePrg(addr uint16, data uint8) {
	if addr < 0x8000 {
		return
	}
	m.selectedBank = int(data & 0b1111)
}

func (m *ines2) ReadChr(addr uint16) uint8 {
	return m.rom.Character.Read(addr)
}

func (m *ines2) WriteChr(addr uint16, data uint8) {
	m.rom.Character.Write(addr, data)
}
