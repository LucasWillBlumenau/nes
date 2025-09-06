package cpu

import (
	"fmt"
)

func Lda(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction LDA...")

	cpu.A = uint8(fetchedValue)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) == 1)
}

func Ldx(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction LDX...")

	cpu.X = uint8(fetchedValue)
	cpu.SetStatusFlag(StatusFlagZero, cpu.X == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.X>>7) == 1)
}

func Ldy(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction LDY...")

	cpu.Y = uint8(fetchedValue)
	cpu.SetStatusFlag(StatusFlagZero, cpu.Y == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.Y>>7) == 1)
}

// address
func Sta(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction STA...")

	cpu.Bus.Write(fetchedValue, cpu.A)
}

// address
func Stx(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction STX...")

	cpu.Bus.Write(fetchedValue, cpu.X)
}

// address
func Sty(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction STY...")

	cpu.Bus.Write(fetchedValue, cpu.Y)
}
