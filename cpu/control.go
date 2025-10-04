package cpu

import "github.com/LucasWillBlumenau/nes/interrupt"

func Brk(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction BRK...")
	cpu.Pc++
	interrupt.InterruptSignal.Send(interrupt.Irq)
}

func Jmp(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction JMP...")
	cpu.Pc = fetchedValue
}

func Jsr(cpu *CPU, fetchedValue uint16) {
	// fmt.Println("Executing instruction JSR...")

	programCounter := cpu.Pc - 1

	hi := uint8((programCounter & 0xFF00) >> 8)
	lo := uint8(programCounter & 0x00FF)

	cpu.Push(hi)
	cpu.Push(lo)

	cpu.Pc = fetchedValue
}

func Rti(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction RTI...")

	cpu.P = (cpu.Pop() & 0b11001111) | (cpu.P & 0b00110000)

	lo := cpu.Pop()
	hi := cpu.Pop()

	cpu.Pc = (uint16(hi) << 8) | uint16(lo)
}

func Rts(cpu *CPU, _ uint16) {
	// fmt.Println("Executing instruction RTS...")

	lo := cpu.Pop()
	hi := cpu.Pop()

	programCounter := (uint16(hi) << 8) | uint16(lo)
	cpu.Pc = programCounter + 1
}
