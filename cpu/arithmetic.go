package cpu

func Adc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction ADC...")

	var carryBit uint16
	if cpu.GetStatusFlag(StatusFlagCarry) {
		carryBit = 1
	}

	sum16 := uint16(cpu.A) + fetchedValue + carryBit
	sum := uint8(sum16)

	a := (cpu.A >> 7) == 1
	s := (sum >> 7) == 1
	m := (uint8(fetchedValue) >> 7) == 1

	cpu.SetStatusFlag(StatusFlagNegative, (sum>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, sum == 0)
	cpu.SetStatusFlag(StatusFlagCarry, sum16 > 0xFF)
	cpu.SetStatusFlag(StatusFlagOverflow, (!a && !m && s) || (a && m && !s))
	cpu.A = sum
}

func Cmp(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction CMP...")

	fv := uint8(fetchedValue)
	diff := cpu.A - fv

	cpu.SetStatusFlag(StatusFlagZero, diff == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (diff>>7) == 1)
	cpu.SetStatusFlag(StatusFlagCarry, fv <= cpu.A)
}

func Cpx(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction CPX...")

	fv := uint8(fetchedValue)
	diff := cpu.X - fv

	cpu.SetStatusFlag(StatusFlagZero, diff == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (diff>>7) == 1)
	cpu.SetStatusFlag(StatusFlagCarry, fv <= cpu.X)
}

func Cpy(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction CPY...")

	fv := uint8(fetchedValue)
	diff := cpu.Y - fv

	cpu.SetStatusFlag(StatusFlagZero, diff == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (diff>>7) == 1)
	cpu.SetStatusFlag(StatusFlagCarry, fv <= cpu.Y)
}

func Sbc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction SBC...")

	var carryBit uint16
	if cpu.GetStatusFlag(StatusFlagCarry) {
		carryBit = 1
	}

	value := fetchedValue ^ 0xFF
	sum16 := uint16(cpu.A) + value + carryBit
	sum := uint8(sum16)

	a := (cpu.A >> 7) == 1
	s := (sum >> 7) == 1
	m := (uint8(value) >> 7) == 1

	cpu.SetStatusFlag(StatusFlagNegative, (sum>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, sum == 0)
	cpu.SetStatusFlag(StatusFlagCarry, sum16 > 0xFF)
	cpu.SetStatusFlag(StatusFlagOverflow, (!a && !m && s) || (a && m && !s))
	cpu.A = sum
}
