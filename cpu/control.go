package cpu

import (
	"fmt"
)

func Brk(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BRK...")

	programCounter := cpu.Pc + 1

	hi := uint8(programCounter >> 8)
	lo := uint8(programCounter & 0x00FF)

	cpu.Push(hi)
	cpu.Push(lo)
	cpu.Push(cpu.P)

	cpu.SetStatusFlag(StatusFlagInterruptDisable, true)

	lo = cpu.BusRead(0xFFFE)
	hi = cpu.BusRead(0xFFFF)

	cpu.Pc = (uint16(hi) << 8) | uint16(lo)
}

func Jmp(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction JMP...")
	cpu.Pc = fetchedValue
}

func Jsr(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction JSR...")

	programCounter := cpu.Pc - 1

	hi := uint8((programCounter & 0xFF00) >> 8)
	lo := uint8(programCounter & 0x00FF)

	cpu.Push(hi)
	cpu.Push(lo)

	cpu.Pc = fetchedValue
}

func Rti(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction RTI...")

	cpu.P = cpu.Pop()

	programCounterLo := cpu.Pop()
	programCounterHi := cpu.Pop()

	cpu.Pc = (uint16(programCounterHi) << 8) | uint16(programCounterLo)
}

func Rts(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction RTS...")

	lo := cpu.Pop()
	hi := cpu.Pop()

	programCounter := (uint16(hi) << 8) | uint16(lo)
	cpu.Pc = programCounter + 1
}
