package cpu

import (
	"fmt"
)

func Asl(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction ASL...")

	value := cpu.BusRead(fetchedValue)
	result := value << 1

	cpu.BusWrite(fetchedValue, result)

	cpu.SetStatusFlag(StatusFlagCarry, (value&0x80) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (result&0x80) != 0)
}

func Lsr(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction LSR...")

	value := cpu.BusRead(fetchedValue)
	result := value >> 1

	cpu.BusWrite(fetchedValue, result)

	cpu.SetStatusFlag(StatusFlagCarry, (value&0x01) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
}

func Rol(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction ROL...")

	var carryBit uint8
	if cpu.GetStatusFlag(StatusFlagCarry) {
		carryBit = 1
	} else {
		carryBit = 0
	}

	value := cpu.BusRead(fetchedValue)
	result := (value << 1) | carryBit

	cpu.BusWrite(fetchedValue, result)

	cpu.SetStatusFlag(StatusFlagCarry, (value&0x80) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (result&0x80) != 0)
}

func Ror(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction ROR...")

	var carryBit uint8
	if cpu.GetStatusFlag(StatusFlagCarry) {
		carryBit = 0x80
	} else {
		carryBit = 0
	}

	value := cpu.BusRead(fetchedValue)
	result := (value >> 1) | carryBit

	cpu.BusWrite(fetchedValue, result)

	cpu.SetStatusFlag(StatusFlagCarry, (value&0x01) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (result&0x80) != 0)
}
