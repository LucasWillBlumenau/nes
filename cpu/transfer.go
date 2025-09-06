package cpu

import (
	"fmt"
)

func Tax(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction TAX...")

	cpu.X = cpu.A
	cpu.SetStatusFlag(StatusFlagZero, cpu.X == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.X>>7) != 0)
}

func Tay(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction TAY...")

	cpu.Y = cpu.A
	cpu.SetStatusFlag(StatusFlagZero, cpu.Y == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.Y>>7) != 0)
}

func Tsx(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction TSX...")

	cpu.X = cpu.Sp
	cpu.SetStatusFlag(StatusFlagZero, cpu.X == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.X>>7) != 0)
}

func Txa(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction TXA...")

	cpu.A = cpu.X
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) != 0)
}

func Txs(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction TXS...")

	cpu.Sp = cpu.X
}

func Tya(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction TYA...")

	cpu.A = cpu.Y
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) != 0)
}
