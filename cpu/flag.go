package cpu

func Clc(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction CLC...")
	cpu.SetStatusFlag(StatusFlagCarry, false)
}

func Cld(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction CLD...")
	cpu.SetStatusFlag(StatusFlagDecimal, false)
}

func Cli(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction CLI...")
	cpu.SetStatusFlag(StatusFlagInterruptDisable, false)
}

func Clv(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction CLV...")
	cpu.SetStatusFlag(StatusFlagOverflow, false)
}

func Sec(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction SEC...")
	cpu.SetStatusFlag(StatusFlagCarry, true)
}

func Sed(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction SED...")
	cpu.SetStatusFlag(StatusFlagDecimal, true)
}

func Sei(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction SEI...")
	cpu.SetStatusFlag(StatusFlagInterruptDisable, true)
}
