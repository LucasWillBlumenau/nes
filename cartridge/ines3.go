package cartridge

type ines3 struct {
	mirroring    MirroringType
	rom          *cartridgeRom
	ram          [1024 * 2]uint8
	selectedBank int
	headers      *cartridgeHeaders
}

func newINES3(rom *cartridgeRom, headers *cartridgeHeaders) mapper {
	return &ines3{
		mirroring:    headers.Mirroring,
		selectedBank: 0,
		rom:          rom,
		ram:          [1024 * 2]uint8{},
		headers:      headers,
	}
}

func (m *ines3) Mirroring() MirroringType {
	return m.mirroring
}

func (m *ines3) ReadPrg(addr16 uint16) uint8 {
	if addr16 < 0x6000 {
		return 0
	}
	if addr16 < 0x8000 {
		return m.ram[addr16&0x7FF]
	}

	addr := int(addr16) - 0x8000
	return m.rom.Program[addr]
}

func (m *ines3) WritePrg(addr uint16, data uint8) {
	if addr < 0x6000 {
		return
	}
	if addr < 0x8000 {
		m.ram[addr&0x7FF] = data
	} else {
		addr -= 0x8000
		m.selectedBank = int(data & m.rom.Program[addr] & 0b11)
	}
}

func (m *ines3) ReadChr(addr uint16) uint8 {
	return m.rom.Character.Read(int(addr) + m.selectedBank*8*1024)
}

func (m *ines3) WriteChr(addr uint16, data uint8) {
	m.rom.Character.Write(int(addr), data)
}
