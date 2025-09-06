package bus

import "github.com/LucasWillBlumenau/nes/cartridge"

const kb = 1024

// low byte, followed by high byte
var nmiAddressLocation = []uint16{0xFFFA, 0xFFFB}
var resetAddressLocation = []uint16{0xFFFC, 0xFFFD}
var irqAddressLocation = []uint16{0xFFFE, 0xFFFF}

type Bus struct {
	ram          []byte
	vram         []byte
	ppuRegisters ppuRegisters
	cartridge    *cartridge.Cartridge
}

type ppuRegisters struct {
	ControlRegister1      uint8
	ControlRegister2      uint8
	StatusRegister        uint8
	SPRRamAddressRegister uint8
	SPRRamIORegister      uint8
	VRamRegister1         uint8
	VRamRegister2         uint8
	VRamIORegister        uint8
}

func NewBus(cartridge *cartridge.Cartridge) *Bus {
	ram := make([]byte, 2*kb)
	vram := make([]byte, 16*kb)
	return &Bus{cartridge: cartridge, ram: ram, vram: vram}
}

func (b *Bus) Write(addr uint16, value uint8) {

}

func (b *Bus) Read(addr uint16) uint8 {
	isReadFromRam := addr < 0x2000
	if isReadFromRam {
		addr = addr & 0x07FF
		return b.ram[addr]
	}

	isReadFromRom := addr >= 0x8000
	if isReadFromRom {
		if addr >= 0xC000 && b.cartridge.ProgramBanks == 1 {
			addr &= 0xBFFF
		}
		addr -= 0x8000
		return b.cartridge.ProgramRom[addr]
	}

	return 0
}
