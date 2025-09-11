package cpu

import (
	"errors"
	"fmt"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/interruption"
)

var ErrInvalidInstruction = errors.New("invalid instruction")

const nmiLowByteAddress = 0xFFFA
const nmiHighByteAddress = 0xFFFB

type CPU struct {
	A   uint8
	X   uint8
	Y   uint8
	P   uint8
	Sp  uint8
	Pc  uint16
	bus *bus.Bus
}

func NewCPU(bus *bus.Bus) *CPU {
	return &CPU{
		A:   0,
		X:   0,
		Y:   0,
		P:   0,
		Sp:  0,
		Pc:  0x08000,
		bus: bus,
	}
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
	c.Sp -= 1
	stackAddress := uint16(c.Sp) | 0x0100
	c.BusWrite(stackAddress, value)
}

func (c *CPU) Pop() uint8 {
	stackAddress := uint16(c.Sp) | 0x0100
	value := c.BusRead(stackAddress)
	c.Sp++
	return value
}

func (c *CPU) Run() {
	for {
		select {
		case interrupt := <-interruption.InterruptionHandler:
			c.attendInterrupt(interrupt)
			continue
		default:
			c.executeInstruction()
		}
	}
}

func (c *CPU) executeInstruction() error {
	opcode := c.BusRead(c.Pc)

	instruction := instructionMap[opcode]
	if instruction == nil {
		return fmt.Errorf("%w: invalid opcode %X", ErrInvalidInstruction, opcode)
	}

	c.Pc++
	value := c.fetchNextValue(instruction.AddressingMode)
	instruction.Dispatch(c, value)
	return nil
}

func (c *CPU) attendInterrupt(interrupt interruption.Interruption) {
	fmt.Printf("Attenting interrupt %d\n", interrupt)

	switch interrupt {
	case interruption.NonMaskableInterrupt:
		lo := c.BusRead(nmiLowByteAddress)
		hi := c.BusRead(nmiHighByteAddress)
		c.Pc = uint16(hi)<<8 + uint16(lo)
	}

}

func (c *CPU) fetchNextValue(addressingMode AddressingMode) uint16 {
	switch addressingMode {
	case Implied, Accumulator:
		return 0
	case Immediate:
		return uint16(c.getImmediateValue())
	case XIndexedAbsoluteValue:
		addr := c.getAbsoluteAddress() + uint16(c.X)
		return uint16(c.BusRead(addr))
	case XIndexedAbsolute:
		return c.getAbsoluteAddress() + uint16(c.X)
	case YIndexedAbsoluteValue:
		addr := c.getAbsoluteAddress() + uint16(c.Y)
		return uint16(c.BusRead(addr))
	case YIndexedAbsolute:
		return c.getAbsoluteAddress() + uint16(c.Y)
	case AbsoluteIndirect:
		indirectAddr := c.getAbsoluteAddress()
		lo := c.BusRead(indirectAddr)
		hi := c.BusRead(indirectAddr + 1)
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
		return uint16(int16(c.Pc) + int16(offset))
	case XIndexedZeroPageIndirectValue:
		addr := uint16(c.getImmediateValue() + c.X)
		lo := c.BusRead(addr)
		hi := c.BusRead(addr + 1)
		addr = uint16(hi)<<8 + uint16(lo)
		return uint16(c.BusRead(addr))
	case XIndexedZeroPageIndirect:
		addr := uint16(c.getImmediateValue() + c.X)
		lo := c.BusRead(addr)
		hi := c.BusRead(addr + 1)
		return uint16(hi)<<8 + uint16(lo)
	case ZeroPageIndirectYIndexedValue:
		addr := uint16(c.getImmediateValue()) + uint16(c.Y)
		lo := c.BusRead(addr)
		hi := c.BusRead(addr + 1)
		addr = uint16(hi)<<8 + uint16(lo)
		return uint16(c.BusRead(addr))
	case ZeroPageIndirectYIndexed:
		addr := uint16(c.getImmediateValue()) + uint16(c.Y)
		lo := c.BusRead(addr)
		hi := c.BusRead(addr + 1)
		return uint16(hi)<<8 + uint16(lo)
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
	c.bus.Write(addr, value)
}
