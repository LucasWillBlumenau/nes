package cpu

func Dec(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction DEC...")

	currentValue := cpu.BusRead(fetchedValue)
	result := currentValue - 1

	cpu.SetStatusFlag(StatusFlagNegative, (result>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)

	cpu.BusWrite(fetchedValue, result)
}

func Dex(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction DEX...")

	cpu.X--
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.X>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.X == 0)
}

func Dey(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction DEY...")

	cpu.Y--
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.Y>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.Y == 0)
}

func Inc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction INC...")

	currentValue := cpu.BusRead(fetchedValue)
	result := currentValue + 1

	cpu.SetStatusFlag(StatusFlagNegative, (result>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)

	cpu.BusWrite(fetchedValue, result)
}

func Inx(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction INX...")

	cpu.X++
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.X>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.X == 0)
}

func Iny(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction INY...")

	cpu.Y++
	cpu.SetStatusFlag(StatusFlagNegative, (cpu.Y>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.Y == 0)
}
