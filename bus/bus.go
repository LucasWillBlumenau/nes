package bus

import (
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/ppu"
)

// memory unit constants
const kb = 1024

// ppu register constants
const ppuControlPortAddr uint16 = 0x2000
const ppuMaskPortAddr uint16 = 0x2001
const ppuStatusPortAddr uint16 = 0x2002
const ppuOAMAddressPortAddr = 0x2003
const ppuOAMDataPortAddr = 0x2004
const ppuScrollPortAddr = 0x2005
const ppuVRamAddressPortAddr = 0x2006
const ppuVRamDataPortAddr = 0x2007

// low byte, followed by high byte
var nmiAddressLocation = []uint16{0xFFFA, 0xFFFB}
var resetAddressLocation = []uint16{0xFFFC, 0xFFFD}
var irqAddressLocation = []uint16{0xFFFE, 0xFFFF}

type Bus struct {
	ram       []uint8
	cartridge *cartridge.Cartridge
	ppu       *ppu.PPU
}

func NewBus(ppu *ppu.PPU, cartridge *cartridge.Cartridge) *Bus {
	ram := make([]byte, 2*kb)
	return &Bus{cartridge: cartridge, ram: ram, ppu: ppu}
}

func (b *Bus) Write(addr uint16, value uint8) {
	valueAddress := b.getValueAddress(addr)
	if valueAddress != nil {
		*valueAddress = value
	}

	addr &= 0x2007
	switch addr {
	case ppuControlPortAddr:
		b.ppu.WritePPUControlPort(value)
	case ppuMaskPortAddr:
		b.ppu.WritePPUMaskPort(value)
	case ppuOAMAddressPortAddr:
		b.ppu.WriteOAMAddrPort(value)
	case ppuOAMDataPortAddr:
		b.ppu.WriteOAMDataPort(value)
	case ppuScrollPortAddr:
		b.ppu.WritePPUScrollPort(value)
	case ppuVRamAddressPortAddr:
		b.ppu.WritePPUAddrPort(value)
	case ppuVRamDataPortAddr:
		b.ppu.WritePPUDataPort(value)
	}

	// TODO: check out behavior when doing writes to addresses that should only be read
	// panic("invalid address found")
}

func (b *Bus) Read(addr uint16) uint8 {
	valueAddress := b.getValueAddress(addr)
	if valueAddress != nil {
		return *valueAddress
	}

	addr &= 0x2007
	switch addr {
	case ppuStatusPortAddr:
		return b.ppu.ReadStatusPort()
	case ppuOAMDataPortAddr:
		return b.ppu.ReadOAMDataPort()
	case ppuVRamDataPortAddr:
		return b.ppu.ReadVRamDataPort()
	}

	// TODO: check out behavior on read to write-only ports
	// panic("invalid address found")
	return 0
}

func (b *Bus) getValueAddress(addr uint16) *uint8 {
	isRam := addr < 0x2000
	if isRam {
		addr &= 0x07FF
		return &b.ram[addr]
	}

	isReadFromRom := addr >= 0x8000
	if isReadFromRom {
		if addr >= 0xC000 && b.cartridge.ProgramBanks == 1 {
			addr &= 0xBFFF
		}
		addr -= 0x8000
		return &b.cartridge.ProgramRom[addr]
	}
	return nil
}
