package cpu

import (
	"fmt"
)

func And(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction AND...")

	cpu.A &= uint8(fetchedValue)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) == 1)
}

func Bit(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BIT...")

	value := cpu.A & uint8(fetchedValue)

	cpu.SetStatusFlag(StatusFlagNegative, (fetchedValue>>7) == 1)
	cpu.SetStatusFlag(StatusFlagOverflow, ((fetchedValue>>6)&0b01) == 1)
	cpu.SetStatusFlag(StatusFlagZero, value == 0)
}

func Eor(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction EOR...")

	cpu.A ^= uint8(fetchedValue)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
}

func Ora(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction ORA...")

	cpu.A |= uint8(fetchedValue)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
}
