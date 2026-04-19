package cartridge

type ines4 struct {
	mirroring    MirroringType
	rom          *cartridgeRom
	ram          [1024 * 2]uint8
	selectedBank uint
	headers      *cartridgeHeaders
}

func newINES4(rom *cartridgeRom, headers *cartridgeHeaders) mapper {
	return &ines4{
		mirroring:    headers.Mirroring,
		selectedBank: 0,
		rom:          rom,
		ram:          [1024 * 2]uint8{},
		headers:      headers,
	}
}

func (m *ines4) Mirroring() MirroringType {
	return m.mirroring
}

func (m *ines4) ReadPrg(addr16 uint16) uint8 {
	if addr16 < 0x6000 {
		return 0
	}
	if addr16 < 0x8000 {
		return m.ram[addr16-0x6000]
	}

	addr := int(addr16) - 0x8000
	return m.rom.Program[addr]
}

func (m *ines4) WritePrg(addr uint16, data uint8) {
	if addr < 0x6000 {
		return
	}
	if addr < 0x8000 {
		m.ram[addr-0x6000] = data
	} else {
		addr -= 0x8000
		m.selectedBank = uint(data & m.rom.Program[addr] & 0b11)
	}
}

func (m *ines4) ReadChr(addr uint16) uint8 {
	const bankSize uint = 8 * 1024
	return m.rom.Character.Read(uint(addr) + m.selectedBank*bankSize)
}

func (m *ines4) WriteChr(addr uint16, data uint8) {
	m.rom.Character.Write(uint(addr), data)
}
