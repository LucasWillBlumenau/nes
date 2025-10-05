package cpu

func Bcc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BCC...")

	if !cpu.GetStatusFlag(StatusFlagCarry) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Bcs(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BCS...")

	if cpu.GetStatusFlag(StatusFlagCarry) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Beq(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BEQ...")

	if cpu.GetStatusFlag(StatusFlagZero) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Bmi(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BMI...")

	if cpu.GetStatusFlag(StatusFlagNegative) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Bne(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BNE...")

	if !cpu.GetStatusFlag(StatusFlagZero) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Bpl(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BPL...")

	if !cpu.GetStatusFlag(StatusFlagNegative) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Bvc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BVC...")

	if !cpu.GetStatusFlag(StatusFlagOverflow) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}

func Bvs(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BVS...")

	if cpu.GetStatusFlag(StatusFlagOverflow) {
		if (cpu.Pc & 0xFF00) != (fetchedValue & 0xFF00) {
			cpu.ElapseCycle()
		}
		cpu.Pc = fetchedValue
		cpu.ElapseCycle()
	}
}
