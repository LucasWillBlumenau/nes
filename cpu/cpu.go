package cpu

import (
	"errors"
	"fmt"

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
	A               uint8
	X               uint8
	Y               uint8
	P               uint8
	Sp              uint8
	Pc              uint16
	elapsedCycles   int64
	bus             *Bus
	extraCycles     uint16
	dmaOccuring     bool
	dmaPage         uint16
	dmaFetches      uint16
	lastInstruction *Instruction
	firstOperand    uint8
	secondOperand   uint8
}

func NewCPU(bus *Bus) *CPU {
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

func (c *CPU) ElapsedCycles() int64 {
	return c.elapsedCycles
}

func (c *CPU) Reset() {
	interrupt.InterruptSignal.Send(interrupt.Reset)
}

func (c *CPU) ResetState() {
	lo := c.BusRead(resetLowByteAddress)
	hi := c.BusRead(resetHighByteAddress)
	c.A = 0
	c.X = 0
	c.Y = 0
	c.P = 0b00100100
	c.Sp = 0xFD
	c.Pc = uint16(hi)<<8 + uint16(lo)
}

func (c *CPU) Run() (uint16, error) {
	c.extraCycles = 0
	if interrupt, ok := interrupt.InterruptSignal.Read(); ok {
		c.attendInterrupt(interrupt)
		c.elapsedCycles += 7
		return 7, nil
	}

	if c.dmaOccuring {
		addr := c.dmaPage | c.dmaFetches
		value := c.BusRead(addr)
		c.bus.OAMWrite(value)
		c.dmaFetches++
		c.dmaOccuring = c.dmaFetches < 256
		c.elapsedCycles += 2
		return 2, nil
	}

	cyclesTaken, err := c.executeInstruction()
	if err != nil {
		return 0, err
	}
	c.elapsedCycles += int64(cyclesTaken)
	return cyclesTaken, err
}

func (c *CPU) executeInstruction() (uint16, error) {
	c.firstOperand = 0
	c.secondOperand = 0
	opcode := c.BusRead(c.Pc)
	instruction := instructionMap[opcode]
	if instruction.Dispatch == nil {
		return 0, fmt.Errorf("%w: invalid opcode %02X", ErrInvalidInstruction, opcode)
	}
	c.Pc++
	c.lastInstruction = &instruction

	value := c.fetchNextValue(instruction.AddressingMode)

	instruction.Dispatch(c, value)
	return instruction.Cycles + c.extraCycles, nil
}

func (c *CPU) attendInterrupt(interruptValue interrupt.Interrupt) {
	if interruptValue == interrupt.Reset {
		c.ResetState()
		return
	}

	hi := uint8(c.Pc >> 8)
	lo := uint8(c.Pc & 0x00FF)

	c.Push(hi)
	c.Push(lo)
	c.Push(c.P)
	c.SetStatusFlag(StatusFlagInterruptDisable, true)

	var interruptHandlerLowAddrPar, interruptHandlerHighAddrPart uint8
	switch interruptValue {
	case interrupt.NonMaskableInterrupt:
		interruptHandlerLowAddrPar = c.BusRead(nmiLowByteAddress)
		interruptHandlerHighAddrPart = c.BusRead(nmiHighByteAddress)
	case interrupt.Irq:
		interruptHandlerLowAddrPar = c.BusRead(irqLowByteAddress)
		interruptHandlerHighAddrPart = c.BusRead(irqHighByteAddress)
	}
	c.Pc = uint16(interruptHandlerHighAddrPart)<<8 + uint16(interruptHandlerLowAddrPar)
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
		addr = addr + uint16(c.X)
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
		addr = addr + uint16(c.Y)
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
		offset := c.BusRead(c.Pc)
		c.firstOperand = uint8(offset)
		c.Pc++
		positivePart := uint16(offset & 0b01111111)
		negativePart := uint16(offset & 0b10000000)
		nextAddr := c.Pc + positivePart - negativePart
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
		addr := (uint16(hi) << 8) + uint16(lo) + uint16(c.Y)
		return addr
	}
	panic("should never get here")
}

func (c *CPU) getImmediateValue() uint8 {
	value := c.BusRead(c.Pc)
	c.Pc++
	c.firstOperand = value
	return value
}

func (c *CPU) getAbsoluteAddress() uint16 {
	lo := c.BusRead(c.Pc)
	c.Pc++
	c.firstOperand = lo
	hi := c.BusRead(c.Pc)
	c.Pc++
	c.secondOperand = hi
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

func (c *CPU) State() string {
	return fmt.Sprintf(
		"A: %02x, X: %02x, Y: %02x, P: %02x, SP: %02x, PC: %04x",
		c.A,
		c.X,
		c.Y,
		c.P,
		c.Sp,
		c.Pc,
	)
}

func (c *CPU) GetLastInstruction() string {
	if c.lastInstruction == nil {
		return ""
	}
	return c.formatInstrucitionAsAsm(c.lastInstruction)
}

func (c *CPU) formatInstrucitionAsAsm(instruction *Instruction) string {
	var asm string
	switch instruction.AddressingMode {
	case ZeroPage, ZeroPageValue:
		asm = fmt.Sprintf("%s $%02X", instruction.Name, c.firstOperand)
	case XIndexedZeroPage, XIndexedZeroPageValue:
		asm = fmt.Sprintf("%s $%02X, X", instruction.Name, c.firstOperand)
	case YIndexedZeroPage, YIndexedZeroPageValue:
		asm = fmt.Sprintf("%s $%02X, Y", instruction.Name, c.firstOperand)
	case Absolute, AbsoluteValue:
		asm = fmt.Sprintf("%s $%02X%02X", instruction.Name, c.secondOperand, c.firstOperand)
	case Relative:
		addr := int(c.Pc-2) + int(c.firstOperand) + 2
		asm = fmt.Sprintf("%s $%04X", instruction.Name, addr)
	case XIndexedAbsolute, XIndexedAbsoluteValue:
		asm = fmt.Sprintf("%s $%02X%02X, X", instruction.Name, c.secondOperand, c.firstOperand)
	case YIndexedAbsolute, YIndexedAbsoluteValue:
		asm = fmt.Sprintf("%s $%02X%02X, Y", instruction.Name, c.secondOperand, c.firstOperand)
	case AbsoluteIndirect:
		asm = fmt.Sprintf("%s ($%02X%02X)", instruction.Name, c.secondOperand, c.firstOperand)
	case Implied:
		asm = instruction.Name
	case Accumulator:
		asm = fmt.Sprintf("%s A", instruction.Name)
	case Immediate:
		asm = fmt.Sprintf("%s #$%02X", instruction.Name, c.firstOperand)
	case XIndexedZeroPageIndirect, XIndexedZeroPageIndirectValue:
		asm = fmt.Sprintf("%s ($%02X%02X, X)", instruction.Name, c.secondOperand, c.firstOperand)
	case ZeroPageIndirectYIndexed, ZeroPageIndirectYIndexedValue:
		asm = fmt.Sprintf("%s ($%02X), Y", instruction.Name, c.firstOperand)
	default:
		panic("invalid addressing mode")
	}
	return asm
}
