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

func AslAccumulator(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction ASL...")
	value := cpu.A
	result := value << 1
	cpu.A = result
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

func LsrAccumulator(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction LSR...")
	value := cpu.A
	result := value >> 1
	cpu.A = result
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

func RolAccumulator(cpu *CPU, _ uint16) {
	fmt.Println("Executing instruction ROL...")
	var carryBit uint8
	if cpu.GetStatusFlag(StatusFlagCarry) {
		carryBit = 1
	} else {
		carryBit = 0
	}

	value := cpu.A
	result := (value << 1) | carryBit

	cpu.A = result
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

func RorAccumulator(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction ROR...")
	var carryBit uint8
	if cpu.GetStatusFlag(StatusFlagCarry) {
		carryBit = 0x80
	} else {
		carryBit = 0
	}

	value := cpu.A
	result := (value >> 1) | carryBit
	cpu.A = result

	cpu.SetStatusFlag(StatusFlagCarry, (value&0x01) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (result&0x80) != 0)
}
