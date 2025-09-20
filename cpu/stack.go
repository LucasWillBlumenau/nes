package cpu

func Pha(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction PHA...")
	cpu.Push(cpu.A)
}

func Php(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction PHP...")
	cpu.Push(cpu.P)
}

func Pla(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction PLA...")
	cpu.A = cpu.Pop()
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.A>>7) != 0)
}

func Plp(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction PLP...")
	cpu.P = cpu.Pop()
}
