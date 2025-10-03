package cpu

import (
	"errors"
	"fmt"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/interrupt"
)

var ErrInvalidInstruction = errors.New("invalid instruction")

const nmiLowByteAddress = 0xFFFA
const nmiHighByteAddress = 0xFFFB
const resetLowByteAddress = 0xFFFC
const resetHighByteAddress = 0xFFFD
const irqLowByteAddress = 0xFFFE
const irqHighByteAddress = 0xFFFF

type StatusFlag uint8

const (
	StatusFlagCarry            = 0
	StatusFlagZero             = 1
	StatusFlagInterruptDisable = 2
	StatusFlagDecimal          = 3
	StatusFlagOverflow         = 6
	StatusFlagNegative         = 7
)

type CPU struct {
	A           uint8
	X           uint8
	Y           uint8
	P           uint8
	Sp          uint8
	Pc          uint16
	bus         *bus.Bus
	extraCycles uint16
	dmaOccuring bool
	dmaPage     uint16
	dmaFetches  uint16
}

func NewCPU(bus *bus.Bus) *CPU {
	return &CPU{
		A:   0,
		X:   0,
		Y:   0,
		P:   0b00100100,
		Sp:  0xFD,
		Pc:  0,
		bus: bus,
	}
}

func (c *CPU) SetRomEntrypoint() {
	lo := c.bus.Read(resetLowByteAddress)
	hi := c.bus.Read(resetHighByteAddress)
	romEntryPoint := uint16(hi)<<8 | uint16(lo)
	c.Pc = romEntryPoint
}

func (c *CPU) GetStatusFlag(flag StatusFlag) bool {
	return ((c.P >> flag) & 0b00000001) == 1
}

func (c *CPU) SetStatusFlag(flag StatusFlag, value bool) {
	if value {
		c.P |= 1 << flag
	} else {
		c.P &= 0b11111111 ^ (1 << flag)
	}
}

func (c *CPU) Push(value uint8) {
	stackAddress := uint16(c.Sp) | 0x0100
	c.BusWrite(stackAddress, value)
	c.Sp -= 1
}

func (c *CPU) Pop() uint8 {
	c.Sp++
	stackAddress := uint16(c.Sp) | 0x0100
	value := c.BusRead(stackAddress)
	return value
}

func (c *CPU) ElapseCycle() {
	c.extraCycles++
}

func (c *CPU) Run() (uint16, error) {
	c.extraCycles = 0
	if interrupt, ok := interrupt.InterruptSignal.Read(); ok {
		c.attendInterrupt(interrupt)
		return 0, nil
	}

	if c.dmaOccuring {
		addr := c.dmaPage | c.dmaFetches
		value := c.BusRead(addr)
		c.bus.OAMWrite(value)
		c.dmaFetches++
		c.dmaOccuring = c.dmaFetches < 256
		return 2, nil
	}

	return c.executeInstruction()
}

func (c *CPU) executeInstruction() (uint16, error) {
	opcode := c.BusRead(c.Pc)
	instruction := instructionMap[opcode]
	if instruction.Dispatch == nil {
		return 0, fmt.Errorf("%w: invalid opcode %02X", ErrInvalidInstruction, opcode)
	}
	c.Pc++
	fmt.Println(instruction.Stringfy(c))
	value := c.fetchNextValue(instruction.AddressingMode)

	instruction.Dispatch(c, value)
	return instruction.Cycles + c.extraCycles, nil
}

func (c *CPU) attendInterrupt(interruptValue interrupt.Interrupt) {
	hi := uint8(c.Pc >> 8)
	lo := uint8(c.Pc & 0x00FF)

	c.Push(hi)
	c.Push(lo)
	c.Push(c.P)
	c.SetStatusFlag(StatusFlagInterruptDisable, true)

	var intHandlerLo, intHandlerHi uint8
	switch interruptValue {
	case interrupt.NonMaskableInterrupt:
		intHandlerLo = c.BusRead(nmiLowByteAddress)
		intHandlerHi = c.BusRead(nmiHighByteAddress)
	case interrupt.Irq:
		intHandlerLo = c.BusRead(irqLowByteAddress)
		intHandlerHi = c.BusRead(irqHighByteAddress)
	}
	c.Pc = uint16(intHandlerHi)<<8 + uint16(intHandlerLo)
}

func (c *CPU) fetchNextValue(addressingMode AddressingMode) uint16 {
	switch addressingMode {
	case Implied, Accumulator:
		return 0
	case Immediate:
		return uint16(c.getImmediateValue())
	case XIndexedAbsoluteValue:
		addr := c.getAbsoluteAddress()
		addrPage := addr & 0xFF00
		addr = addr + uint16(c.X)
		if addrPage != (addr & 0xFF00) {
			c.extraCycles++
		}
		return uint16(c.BusRead(addr))
	case XIndexedAbsolute:
		addr := c.getAbsoluteAddress()
		addrPage := addr & 0xFF00
		addr = addr + uint16(c.X)
		if addrPage != (addr & 0xFF00) {
			c.extraCycles++
		}
		return addr
	case YIndexedAbsoluteValue:
		addr := c.getAbsoluteAddress()
		addrPage := addr & 0xFF00
		addr = addr + uint16(c.Y)
		if addrPage != (addr & 0xFF00) {
			c.extraCycles++
		}
		return uint16(c.BusRead(addr))
	case YIndexedAbsolute:
		addr := c.getAbsoluteAddress()
		addrPage := addr & 0xFF00
		addr = addr + uint16(c.Y)
		if addrPage != (addr & 0xFF00) {
			c.extraCycles++
		}
		return addr
	case AbsoluteIndirect:
		loAddr := c.getAbsoluteAddress()
		hiAddr := loAddr + 1
		if loAddr&0x00FF == 0xFF {
			hiAddr -= 0x0100
		}
		lo := c.BusRead(loAddr)
		hi := c.BusRead(hiAddr)
		return uint16(hi)<<8 + uint16(lo)
	case AbsoluteValue:
		addr := c.getAbsoluteAddress()
		return uint16(c.BusRead(addr))
	case Absolute:
		return c.getAbsoluteAddress()
	case ZeroPageValue:
		addr := uint16(c.getImmediateValue())
		return uint16(c.BusRead(addr))
	case ZeroPage:
		return uint16(c.getImmediateValue())
	case XIndexedZeroPageValue:
		addr := uint16(c.getImmediateValue() + c.X)
		return uint16(c.BusRead(addr))
	case XIndexedZeroPage:
		return uint16(c.getImmediateValue() + c.X)
	case YIndexedZeroPageValue:
		addr := uint16(c.getImmediateValue() + c.Y)
		return uint16(c.BusRead(addr))
	case YIndexedZeroPage:
		return uint16(c.getImmediateValue() + c.Y)
	case Relative:
		offset := int8(c.BusRead(c.Pc))
		c.Pc++
		nextAddr := uint16(int16(c.Pc) + int16(offset))
		if (nextAddr & 0xFF00) != (c.Pc & 0xFF00) {
			c.extraCycles++
		}
		return nextAddr
	case XIndexedZeroPageIndirectValue:
		indirectAddr := c.getImmediateValue() + c.X
		lo := c.BusRead(uint16(indirectAddr))
		hi := c.BusRead(uint16(indirectAddr + 1))
		addr := uint16(hi)<<8 | uint16(lo)
		return uint16(c.BusRead(addr))
	case XIndexedZeroPageIndirect:
		indirectAddr := c.getImmediateValue() + c.X
		lo := c.BusRead(uint16(indirectAddr))
		hi := c.BusRead(uint16(indirectAddr + 1))
		return uint16(hi)<<8 | uint16(lo)
	case ZeroPageIndirectYIndexedValue:
		value := c.getImmediateValue()
		lo := c.BusRead(uint16(value))
		hi := c.BusRead(uint16(value + 1))
		baseAddr := (uint16(hi) << 8) + uint16(lo)
		baseAddrPage := baseAddr & 0xFF00
		addr := (uint16(hi) << 8) + uint16(lo) + uint16(c.Y)
		if (addr & 0xFF00) != baseAddrPage {
			c.extraCycles++
		}
		return uint16(c.BusRead(addr))
	case ZeroPageIndirectYIndexed:
		value := c.getImmediateValue()
		lo := c.BusRead(uint16(value))
		hi := c.BusRead(uint16(value + 1))
		baseAddr := (uint16(hi) << 8) + uint16(lo)
		baseAddrPage := baseAddr & 0xFF00
		addr := (uint16(hi) << 8) + uint16(lo) + uint16(c.Y)
		if (addr & 0xFF00) != baseAddrPage {
			c.extraCycles++
		}
		return addr
	}
	panic("should never get here")
}

func (c *CPU) getImmediateValue() uint8 {
	value := c.BusRead(c.Pc)
	c.Pc++
	return value
}

func (c *CPU) getAbsoluteAddress() uint16 {
	lo := c.BusRead(c.Pc)
	c.Pc++
	hi := c.BusRead(c.Pc)
	c.Pc++
	return uint16(hi)<<8 + uint16(lo)
}

func (c *CPU) BusRead(addr uint16) uint8 {
	return c.bus.Read(addr)
}

func (c *CPU) BusWrite(addr uint16, value uint8) {
	dmaRequested := c.bus.Write(addr, value)
	if dmaRequested {
		c.dmaOccuring = true
		c.dmaFetches = 0
		c.dmaPage = uint16(value) << 8
	}
}
