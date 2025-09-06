package cpu

import (
	"fmt"
)

func Bcc(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BCC...")

	if !cpu.GetStatusFlag(StatusFlagCarry) {
		cpu.Pc = fetchedValue
	}
}

func Bcs(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BCS...")

	if cpu.GetStatusFlag(StatusFlagCarry) {
		cpu.Pc = fetchedValue
	}
}

func Beq(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BEQ...")

	if cpu.GetStatusFlag(StatusFlagZero) {
		cpu.Pc = fetchedValue
	}
}

func Bmi(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BMI...")

	if cpu.GetStatusFlag(StatusFlagNegative) {
		cpu.Pc = fetchedValue
	}
}

func Bne(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BNE...")

	if !cpu.GetStatusFlag(StatusFlagZero) {
		cpu.Pc = fetchedValue
	}
}

func Bpl(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BPL...")

	if !cpu.GetStatusFlag(StatusFlagNegative) {
		cpu.Pc = fetchedValue
	}
}

func Bvc(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BVC...")

	if !cpu.GetStatusFlag(StatusFlagOverflow) {
		cpu.Pc = fetchedValue
	}
}

func Bvs(cpu *CPU, fetchedValue uint16) {
	fmt.Println("Executing instruction BVS...")

	if cpu.GetStatusFlag(StatusFlagOverflow) {
		cpu.Pc = fetchedValue
	}
}
