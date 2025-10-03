package cpu

func Asl(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction ASL...")

	value := cpu.BusRead(fetchedValue)
	result := value << 1

	cpu.BusWrite(fetchedValue, result)

	cpu.SetStatusFlag(StatusFlagCarry, (value&0x80) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (result&0x80) != 0)
}

func AslAccumulator(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction ASL...")
	value := cpu.A
	result := value << 1
	cpu.A = result
	cpu.SetStatusFlag(StatusFlagCarry, (value&0x80) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, (result&0x80) != 0)
}

func Lsr(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction LSR...")
	value := cpu.BusRead(fetchedValue)
	result := value >> 1
	cpu.BusWrite(fetchedValue, result)
	cpu.SetStatusFlag(StatusFlagCarry, (value&0x01) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, false)
}

func LsrAccumulator(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction LSR...")
	value := cpu.A
	result := value >> 1
	cpu.A = result
	cpu.SetStatusFlag(StatusFlagCarry, (value&0x01) != 0)
	cpu.SetStatusFlag(StatusFlagZero, result == 0)
	cpu.SetStatusFlag(StatusFlagNegative, false)
}

func Rol(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction ROL...")
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
	// fmt.Println("Executing instruction ROL...")
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
	// fmt.Println("Executing instruction ROR...")
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
	// fmt.Println("Executing instruction ROR...")
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

func Slo(cpu *CPU, fetchedValue uint16) {
	value := cpu.BusRead(fetchedValue)
	carry := value>>7 == 1
	value = value << 1

	cpu.SetStatusFlag(StatusFlagCarry, carry)
	cpu.BusWrite(fetchedValue, value)

	cpu.A |= value
	cpu.SetStatusFlag(StatusFlagNegative, cpu.A>>7 == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
}

func Rla(cpu *CPU, fetchedValue uint16) {
	value := cpu.BusRead(fetchedValue)
	var currentCarry uint8 = 0
	if cpu.GetStatusFlag(StatusFlagCarry) {
		currentCarry = 1
	}

	newCarry := value>>7 == 1
	value = value<<1 | currentCarry
	cpu.SetStatusFlag(StatusFlagCarry, newCarry)

	cpu.BusWrite(fetchedValue, value)

	cpu.A &= value
	cpu.SetStatusFlag(StatusFlagNegative, cpu.A>>7 == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
}

func Sre(cpu *CPU, fetchedValue uint16) {
	value := cpu.BusRead(fetchedValue)
	carry := value&1 == 1
	value = value >> 1

	cpu.SetStatusFlag(StatusFlagCarry, carry)
	cpu.BusWrite(fetchedValue, value)

	cpu.A ^= value
	cpu.SetStatusFlag(StatusFlagNegative, cpu.A>>7 == 1)
	cpu.SetStatusFlag(StatusFlagZero, cpu.A == 0)
}

func Rra(cpu *CPU, fetchedValue uint16) {
	value := cpu.BusRead(fetchedValue)
	var currentCarry uint8 = 0
	if cpu.GetStatusFlag(StatusFlagCarry) {
		currentCarry = 1
	}

	newCarry := value&0b01 == 1
	value = (value >> 1) | (currentCarry << 7)
	cpu.SetStatusFlag(StatusFlagCarry, newCarry)
	cpu.BusWrite(fetchedValue, value)

	var carryBit uint16
	if newCarry {
		carryBit = 1
	}

	sum16 := uint16(cpu.A) + uint16(value) + carryBit
	sum := uint8(sum16)

	a := (cpu.A >> 7) == 1
	s := (sum >> 7) == 1
	m := (value >> 7) == 1

	cpu.SetStatusFlag(StatusFlagNegative, (sum>>7) == 1)
	cpu.SetStatusFlag(StatusFlagZero, sum == 0)
	cpu.SetStatusFlag(StatusFlagCarry, sum16 > 0xFF)
	cpu.SetStatusFlag(StatusFlagOverflow, (!a && !m && s) || (a && m && !s))
	cpu.A = sum
}
