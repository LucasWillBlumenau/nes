package bus

import (
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/ppu"
)

const ppuControlPortAddr = 0x2000
const ppuMaskPortAddr = 0x2001
const ppuStatusPortAddr = 0x2002
const ppuOAMAddressPortAddr = 0x2003
const ppuOAMDataPortAddr = 0x2004
const ppuScrollPortAddr = 0x2005
const ppuVRamAddressPortAddr = 0x2006
const ppuVRamDataPortAddr = 0x2007
const ppuOAMDMAPortAddr = 0x4014
const ppuJoypadOnePortAddr = 0x4016
const ppuJoypadTwoPortAddr = 0x4017

type Bus struct {
	ram       []uint8
	cartridge *cartridge.Cartridge
	ppu       *ppu.PPU
	joypadOne *joypad.Joypad
	joypadTwo *joypad.Joypad
}

func NewBus(
	ppu *ppu.PPU,
	cartridge *cartridge.Cartridge,
	joypadOne *joypad.Joypad,
	joypadTwo *joypad.Joypad,
) *Bus {
	ram := make([]byte, 2*1024)
	return &Bus{
		cartridge: cartridge,
		ram:       ram, ppu: ppu,
		joypadOne: joypadOne,
		joypadTwo: joypadTwo,
	}
}

func (b *Bus) Write(addr uint16, value uint8) bool {
	valueAddress := b.getRamAddress(addr)
	if valueAddress != nil {
		*valueAddress = value
	}

	if addr < 0x4000 {
		addr &= 0x2007
		switch addr {
		case ppuControlPortAddr:
			b.ppu.WritePPUControlPort(value)
		case ppuMaskPortAddr:
			b.ppu.WritePPUMaskPort(value)
		case ppuOAMAddressPortAddr:
			b.ppu.WriteOAMAddrPort(value)
		case ppuOAMDataPortAddr:
			b.OAMWrite(value)
		case ppuScrollPortAddr:
			b.ppu.WritePPUScrollPort(value)
		case ppuVRamAddressPortAddr:
			b.ppu.WritePPUAddrPort(value)
		case ppuVRamDataPortAddr:
			b.ppu.WritePPUDataPort(value)
		}
	} else if addr < 0x8000 {
		switch addr {
		case ppuOAMDMAPortAddr:
			return true
		case ppuJoypadOnePortAddr,
			ppuJoypadTwoPortAddr:
			b.joypadOne.Write(value)
			b.joypadTwo.Write(value)
		}
	}
	addr -= 0x8000
	b.cartridge.WritePrgRom(addr, value)
	return false
}

func (b *Bus) OAMWrite(value uint8) {
	b.ppu.WriteOAMDataPort(value)
}

func (b *Bus) Read(addr uint16) uint8 {
	valueAddress := b.getRamAddress(addr)
	if valueAddress != nil {
		return *valueAddress
	}
	if addr < 0x4000 {
		addr &= 0x2007
		switch addr {
		case ppuStatusPortAddr:
			return b.ppu.ReadStatusPort()
		case ppuOAMDataPortAddr:
			return b.ppu.ReadOAMDataPort()
		case ppuVRamDataPortAddr:
			return b.ppu.ReadVRamDataPort()
		}
	} else if addr < 0x8000 {
		switch addr {
		case ppuJoypadOnePortAddr:
			value := b.joypadOne.Read()
			return value
		case ppuJoypadTwoPortAddr:
			value := b.joypadTwo.Read()
			return value
		}
	}
	addr -= 0x8000
	return b.cartridge.ReadPrgRom(addr)
}

func (b *Bus) getRamAddress(addr uint16) *uint8 {
	isRam := addr < 0x2000
	if isRam {
		addr &= 0x07FF
		return &b.ram[addr]
	}
	return nil
}
