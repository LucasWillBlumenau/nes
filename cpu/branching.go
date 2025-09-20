package cpu

func Bcc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BCC...")

	if !cpu.GetStatusFlag(StatusFlagCarry) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Bcs(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BCS...")

	if cpu.GetStatusFlag(StatusFlagCarry) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Beq(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BEQ...")

	if cpu.GetStatusFlag(StatusFlagZero) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Bmi(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BMI...")

	if cpu.GetStatusFlag(StatusFlagNegative) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Bne(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BNE...")

	if !cpu.GetStatusFlag(StatusFlagZero) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Bpl(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BPL...")

	if !cpu.GetStatusFlag(StatusFlagNegative) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Bvc(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BVC...")

	if !cpu.GetStatusFlag(StatusFlagOverflow) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}

func Bvs(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BVS...")

	if cpu.GetStatusFlag(StatusFlagOverflow) {
		cpu.Pc = fetchedValue
		cpu.SetBranchTaken()
	}
}
