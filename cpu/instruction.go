package cpu

type AddressingMode uint8

const (
	Implied = iota
	Accumulator
	Immediate
	XIndexedAbsolute
	YIndexedAbsolute
	AbsoluteIndirect
	Absolute
	ZeroPage
	XIndexedZeroPage
	YIndexedZeroPage
	Relative
	XIndexedZeroPageIndirect
	ZeroPageIndirectYIndexed
	XIndexedAbsoluteValue
	YIndexedAbsoluteValue
	AbsoluteValue
	ZeroPageValue
	XIndexedZeroPageValue
	YIndexedZeroPageValue
	XIndexedZeroPageIndirectValue
	ZeroPageIndirectYIndexedValue
)

func (a AddressingMode) String() string {
	switch a {
	case Implied:
		return "Implied"
	case Accumulator:
		return ""
	case Immediate:
		return "Immediate"
	case XIndexedAbsolute:
		return "XIndexedAbsolute"
	case YIndexedAbsolute:
		return "YIndexedAbsolute"
	case AbsoluteIndirect:
		return "AbsoluteIndirect"
	case Absolute:
		return "Absolute"
	case ZeroPage:
		return "ZeroPage"
	case XIndexedZeroPage:
		return "XIndexedZeroPage"
	case YIndexedZeroPage:
		return "YIndexedZeroPage"
	case Relative:
		return "Relative"
	case XIndexedZeroPageIndirect:
		return "XIndexedZeroPageIndirect"
	case ZeroPageIndirectYIndexed:
		return "ZeroPageIndirectYIndexed"
	case XIndexedAbsoluteValue:
		return "XIndexedAbsoluteValue"
	case YIndexedAbsoluteValue:
		return "YIndexedAbsoluteValue"
	case AbsoluteValue:
		return "AbsoluteValue"
	case ZeroPageValue:
		return "ZeroPageValue"
	case XIndexedZeroPageValue:
		return "XIndexedZeroPageValue"
	case YIndexedZeroPageValue:
		return "YIndexedZeroPageValue"
	case XIndexedZeroPageIndirectValue:
		return "XIndexedZeroPageIndirectValue"
	case ZeroPageIndirectYIndexedValue:
		return "ZeroPageIndirectYIndexedValue"
	}
	return ""
}

type Instruction struct {
	Name           string
	AddressingMode AddressingMode
	Dispatch       func(*CPU, uint16)
	Cycles         uint16
}

var instructionMap = [256]Instruction{
	0xA9: {Name: "LDA", Dispatch: Lda, AddressingMode: Immediate, Cycles: 2},
	0xAD: {Name: "LDA", Dispatch: Lda, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBD: {Name: "LDA", Dispatch: Lda, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xB9: {Name: "LDA", Dispatch: Lda, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xA5: {Name: "LDA", Dispatch: Lda, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB5: {Name: "LDA", Dispatch: Lda, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0xA1: {Name: "LDA", Dispatch: Lda, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xB1: {Name: "LDA", Dispatch: Lda, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// LDX
	0xA2: {Name: "LDX", Dispatch: Ldx, AddressingMode: Immediate, Cycles: 2},
	0xAE: {Name: "LDX", Dispatch: Ldx, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBE: {Name: "LDX", Dispatch: Ldx, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xA6: {Name: "LDX", Dispatch: Ldx, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB6: {Name: "LDX", Dispatch: Ldx, AddressingMode: YIndexedZeroPageValue, Cycles: 4},
	// LDY
	0xA0: {Name: "LDY", Dispatch: Ldy, AddressingMode: Immediate, Cycles: 2},
	0xAC: {Name: "LDY", Dispatch: Ldy, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBC: {Name: "LDY", Dispatch: Ldy, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xA4: {Name: "LDY", Dispatch: Ldy, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB4: {Name: "LDY", Dispatch: Ldy, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	// LAX
	0xAB: {Name: "*LAX", Dispatch: Lax, AddressingMode: Immediate, Cycles: 2},
	0xAF: {Name: "*LAX", Dispatch: Lax, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBF: {Name: "*LAX", Dispatch: Lax, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xA7: {Name: "*LAX", Dispatch: Lax, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB7: {Name: "*LAX", Dispatch: Lax, AddressingMode: YIndexedZeroPageValue, Cycles: 4},
	0xA3: {Name: "*LAX", Dispatch: Lax, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xB3: {Name: "*LAX", Dispatch: Lax, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},

	// STA
	0x8D: {Name: "STA", Dispatch: Sta, AddressingMode: Absolute, Cycles: 4},
	0x9D: {Name: "STA", Dispatch: Sta, AddressingMode: XIndexedAbsolute, Cycles: 5},
	0x99: {Name: "STA", Dispatch: Sta, AddressingMode: YIndexedAbsolute, Cycles: 5},
	0x85: {Name: "STA", Dispatch: Sta, AddressingMode: ZeroPage, Cycles: 3},
	0x95: {Name: "STA", Dispatch: Sta, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x81: {Name: "STA", Dispatch: Sta, AddressingMode: XIndexedZeroPageIndirect, Cycles: 6},
	0x91: {Name: "STA", Dispatch: Sta, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 6},
	// STX
	0x8E: {Name: "STX", Dispatch: Stx, AddressingMode: Absolute, Cycles: 4},
	0x86: {Name: "STX", Dispatch: Stx, AddressingMode: ZeroPage, Cycles: 3},
	0x96: {Name: "STX", Dispatch: Stx, AddressingMode: YIndexedZeroPage, Cycles: 4},
	// STY
	0x8C: {Name: "STY", Dispatch: Sty, AddressingMode: Absolute, Cycles: 4},
	0x84: {Name: "STY", Dispatch: Sty, AddressingMode: ZeroPage, Cycles: 3},
	0x94: {Name: "STY", Dispatch: Sty, AddressingMode: XIndexedZeroPage, Cycles: 4},
	// SAX
	0x8F: {Name: "*SAX", Dispatch: Sax, AddressingMode: Absolute, Cycles: 4},
	0x87: {Name: "*SAX", Dispatch: Sax, AddressingMode: ZeroPage, Cycles: 3},
	0x97: {Name: "*SAX", Dispatch: Sax, AddressingMode: YIndexedZeroPage, Cycles: 4},
	0x83: {Name: "*SAX", Dispatch: Sax, AddressingMode: XIndexedZeroPageIndirect, Cycles: 6},
	// TAX
	0xAA: {Name: "TAX", Dispatch: Tax, AddressingMode: Implied, Cycles: 2},
	// TAY
	0xA8: {Name: "TAY", Dispatch: Tay, AddressingMode: Implied, Cycles: 2},
	// TSX
	0xBA: {Name: "TSX", Dispatch: Tsx, AddressingMode: Implied, Cycles: 2},
	// TXA
	0x8A: {Name: "TXA", Dispatch: Txa, AddressingMode: Implied, Cycles: 2},
	// TXS
	0x9A: {Name: "TXS", Dispatch: Txs, AddressingMode: Implied, Cycles: 2},
	// TYA
	0x98: {Name: "TYA", Dispatch: Tya, AddressingMode: Implied, Cycles: 2},
	// PHA
	0x48: {Name: "PHA", Dispatch: Pha, AddressingMode: Implied, Cycles: 3},
	// PHP - 🐘
	0x08: {Name: "PHP ", Dispatch: Php, AddressingMode: Implied, Cycles: 3},
	// PLA
	0x68: {Name: "PLA", Dispatch: Pla, AddressingMode: Implied, Cycles: 4},
	// PLP
	0x28: {Name: "PLP", Dispatch: Plp, AddressingMode: Implied, Cycles: 4},
	// ASL
	0x0A: {Name: "ASL", Dispatch: AslAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x0E: {Name: "ASL", Dispatch: Asl, AddressingMode: Absolute, Cycles: 6},
	0x1E: {Name: "ASL", Dispatch: Asl, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x06: {Name: "ASL", Dispatch: Asl, AddressingMode: ZeroPage, Cycles: 5},
	0x16: {Name: "ASL", Dispatch: Asl, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// LSR
	0x4A: {Name: "LSR", Dispatch: LsrAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x4E: {Name: "LSR", Dispatch: Lsr, AddressingMode: Absolute, Cycles: 6},
	0x5E: {Name: "LSR", Dispatch: Lsr, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x46: {Name: "LSR", Dispatch: Lsr, AddressingMode: ZeroPage, Cycles: 5},
	0x56: {Name: "LSR", Dispatch: Lsr, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// ROL
	0x2A: {Name: "ROL", Dispatch: RolAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x2E: {Name: "ROL", Dispatch: Rol, AddressingMode: Absolute, Cycles: 6},
	0x3E: {Name: "ROL", Dispatch: Rol, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x26: {Name: "ROL", Dispatch: Rol, AddressingMode: ZeroPage, Cycles: 5},
	0x36: {Name: "ROL", Dispatch: Rol, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// ROR
	0x6A: {Name: "ROR", Dispatch: RorAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x6E: {Name: "ROR", Dispatch: Ror, AddressingMode: Absolute, Cycles: 6},
	0x7E: {Name: "ROR", Dispatch: Ror, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x66: {Name: "ROR", Dispatch: Ror, AddressingMode: ZeroPage, Cycles: 5},
	0x76: {Name: "ROR", Dispatch: Ror, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// SLO
	0x0F: {Name: "*SLO", Dispatch: Slo, AddressingMode: Absolute, Cycles: 6},
	0x1F: {Name: "*SLO", Dispatch: Slo, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x1B: {Name: "*SLO", Dispatch: Slo, AddressingMode: YIndexedAbsolute, Cycles: 7},
	0x07: {Name: "*SLO", Dispatch: Slo, AddressingMode: ZeroPage, Cycles: 5},
	0x17: {Name: "*SLO", Dispatch: Slo, AddressingMode: XIndexedZeroPage, Cycles: 6},
	0x03: {Name: "*SLO", Dispatch: Slo, AddressingMode: XIndexedZeroPageIndirect, Cycles: 8},
	0x13: {Name: "*SLO", Dispatch: Slo, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 8},
	// RLA
	0x2F: {Name: "*RLA", Dispatch: Rla, AddressingMode: Absolute, Cycles: 6},
	0x3F: {Name: "*RLA", Dispatch: Rla, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x3B: {Name: "*RLA", Dispatch: Rla, AddressingMode: YIndexedAbsolute, Cycles: 7},
	0x27: {Name: "*RLA", Dispatch: Rla, AddressingMode: ZeroPage, Cycles: 5},
	0x37: {Name: "*RLA", Dispatch: Rla, AddressingMode: XIndexedZeroPage, Cycles: 6},
	0x23: {Name: "*RLA", Dispatch: Rla, AddressingMode: XIndexedZeroPageIndirect, Cycles: 8},
	0x33: {Name: "*RLA", Dispatch: Rla, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 8},
	// SRE
	0x4F: {Name: "*SRE", Dispatch: Sre, AddressingMode: Absolute, Cycles: 6},
	0x5F: {Name: "*SRE", Dispatch: Sre, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x5B: {Name: "*SRE", Dispatch: Sre, AddressingMode: YIndexedAbsolute, Cycles: 7},
	0x47: {Name: "*SRE", Dispatch: Sre, AddressingMode: ZeroPage, Cycles: 5},
	0x57: {Name: "*SRE", Dispatch: Sre, AddressingMode: XIndexedZeroPage, Cycles: 6},
	0x43: {Name: "*SRE", Dispatch: Sre, AddressingMode: XIndexedZeroPageIndirect, Cycles: 8},
	0x53: {Name: "*SRE", Dispatch: Sre, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 8},
	// RRA
	0x6F: {Name: "*RRA", Dispatch: Rra, AddressingMode: Absolute, Cycles: 6},
	0x7F: {Name: "*RRA", Dispatch: Rra, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x7B: {Name: "*RRA", Dispatch: Rra, AddressingMode: YIndexedAbsolute, Cycles: 7},
	0x67: {Name: "*RRA", Dispatch: Rra, AddressingMode: ZeroPage, Cycles: 5},
	0x77: {Name: "*RRA", Dispatch: Rra, AddressingMode: XIndexedZeroPage, Cycles: 6},
	0x63: {Name: "*RRA", Dispatch: Rra, AddressingMode: XIndexedZeroPageIndirect, Cycles: 8},
	0x73: {Name: "*RRA", Dispatch: Rra, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 8},
	// AND
	0x29: {Name: "AND", Dispatch: And, AddressingMode: Immediate, Cycles: 2},
	0x2D: {Name: "AND", Dispatch: And, AddressingMode: AbsoluteValue, Cycles: 4},
	0x3D: {Name: "AND", Dispatch: And, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x39: {Name: "AND", Dispatch: And, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0x25: {Name: "AND", Dispatch: And, AddressingMode: ZeroPageValue, Cycles: 3},
	0x35: {Name: "AND", Dispatch: And, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0x21: {Name: "AND", Dispatch: And, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0x31: {Name: "AND", Dispatch: And, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// BIT
	0x2C: {Name: "BIT", Dispatch: Bit, AddressingMode: AbsoluteValue, Cycles: 4},
	0x24: {Name: "BIT", Dispatch: Bit, AddressingMode: ZeroPageValue, Cycles: 3},
	// EOR
	0x49: {Name: "EOR", Dispatch: Eor, AddressingMode: Immediate, Cycles: 2},
	0x4D: {Name: "EOR", Dispatch: Eor, AddressingMode: AbsoluteValue, Cycles: 4},
	0x5D: {Name: "EOR", Dispatch: Eor, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x59: {Name: "EOR", Dispatch: Eor, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0x45: {Name: "EOR", Dispatch: Eor, AddressingMode: ZeroPageValue, Cycles: 3},
	0x55: {Name: "EOR", Dispatch: Eor, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0x41: {Name: "EOR", Dispatch: Eor, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0x51: {Name: "EOR", Dispatch: Eor, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// ORA
	0x09: {Name: "ORA", Dispatch: Ora, AddressingMode: Immediate, Cycles: 2},
	0x0D: {Name: "ORA", Dispatch: Ora, AddressingMode: AbsoluteValue, Cycles: 4},
	0x1D: {Name: "ORA", Dispatch: Ora, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x19: {Name: "ORA", Dispatch: Ora, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0x05: {Name: "ORA", Dispatch: Ora, AddressingMode: ZeroPageValue, Cycles: 3},
	0x15: {Name: "ORA", Dispatch: Ora, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0x01: {Name: "ORA", Dispatch: Ora, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0x11: {Name: "ORA", Dispatch: Ora, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// ADC
	0x69: {Name: "ADC", Dispatch: Adc, AddressingMode: Immediate, Cycles: 2},
	0x6D: {Name: "ADC", Dispatch: Adc, AddressingMode: AbsoluteValue, Cycles: 4},
	0x7D: {Name: "ADC", Dispatch: Adc, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x79: {Name: "ADC", Dispatch: Adc, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0x65: {Name: "ADC", Dispatch: Adc, AddressingMode: ZeroPageValue, Cycles: 3},
	0x75: {Name: "ADC", Dispatch: Adc, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0x61: {Name: "ADC", Dispatch: Adc, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0x71: {Name: "ADC", Dispatch: Adc, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// CMP
	0xC9: {Name: "CMP", Dispatch: Cmp, AddressingMode: Immediate, Cycles: 2},
	0xCD: {Name: "CMP", Dispatch: Cmp, AddressingMode: AbsoluteValue, Cycles: 4},
	0xDD: {Name: "CMP", Dispatch: Cmp, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xD9: {Name: "CMP", Dispatch: Cmp, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xC5: {Name: "CMP", Dispatch: Cmp, AddressingMode: ZeroPageValue, Cycles: 3},
	0xD5: {Name: "CMP", Dispatch: Cmp, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0xC1: {Name: "CMP", Dispatch: Cmp, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xD1: {Name: "CMP", Dispatch: Cmp, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// CPX
	0xE0: {Name: "CPX", Dispatch: Cpx, AddressingMode: Immediate, Cycles: 2},
	0xEC: {Name: "CPX", Dispatch: Cpx, AddressingMode: AbsoluteValue, Cycles: 4},
	0xE4: {Name: "CPX", Dispatch: Cpx, AddressingMode: ZeroPageValue, Cycles: 3},
	// CPY
	0xC0: {Name: "CPY", Dispatch: Cpy, AddressingMode: Immediate, Cycles: 2},
	0xCC: {Name: "CPY", Dispatch: Cpy, AddressingMode: AbsoluteValue, Cycles: 4},
	0xC4: {Name: "CPY", Dispatch: Cpy, AddressingMode: ZeroPageValue, Cycles: 3},
	// SBC
	0xE9: {Name: "SBC", Dispatch: Sbc, AddressingMode: Immediate, Cycles: 2},
	0xEB: {Name: "*SBC", Dispatch: Sbc, AddressingMode: Immediate, Cycles: 2},
	0xED: {Name: "SBC", Dispatch: Sbc, AddressingMode: AbsoluteValue, Cycles: 4},
	0xFD: {Name: "SBC", Dispatch: Sbc, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xF9: {Name: "SBC", Dispatch: Sbc, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xE5: {Name: "SBC", Dispatch: Sbc, AddressingMode: ZeroPageValue, Cycles: 3},
	0xF5: {Name: "SBC", Dispatch: Sbc, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0xE1: {Name: "SBC", Dispatch: Sbc, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xF1: {Name: "SBC", Dispatch: Sbc, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// DCP
	0xCF: {Name: "*DCP", Dispatch: Dcp, AddressingMode: Absolute, Cycles: 6},
	0xDF: {Name: "*DCP", Dispatch: Dcp, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0xDB: {Name: "*DCP", Dispatch: Dcp, AddressingMode: YIndexedAbsolute, Cycles: 7},
	0xC7: {Name: "*DCP", Dispatch: Dcp, AddressingMode: ZeroPage, Cycles: 5},
	0xD7: {Name: "*DCP", Dispatch: Dcp, AddressingMode: XIndexedZeroPage, Cycles: 6},
	0xC3: {Name: "*DCP", Dispatch: Dcp, AddressingMode: XIndexedZeroPageIndirect, Cycles: 8},
	0xD3: {Name: "*DCP", Dispatch: Dcp, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 8},
	// ISC
	0xEF: {Name: "*ISB", Dispatch: Isc, AddressingMode: Absolute, Cycles: 6},
	0xFF: {Name: "*ISB", Dispatch: Isc, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0xFB: {Name: "*ISB", Dispatch: Isc, AddressingMode: YIndexedAbsolute, Cycles: 7},
	0xE7: {Name: "*ISB", Dispatch: Isc, AddressingMode: ZeroPage, Cycles: 5},
	0xF7: {Name: "*ISB", Dispatch: Isc, AddressingMode: XIndexedZeroPage, Cycles: 6},
	0xE3: {Name: "*ISB", Dispatch: Isc, AddressingMode: XIndexedZeroPageIndirect, Cycles: 8},
	0xF3: {Name: "*ISB", Dispatch: Isc, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 8},

	// DEC
	0xCE: {Name: "DEC", Dispatch: Dec, AddressingMode: Absolute, Cycles: 6},
	0xDE: {Name: "DEC", Dispatch: Dec, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0xC6: {Name: "DEC", Dispatch: Dec, AddressingMode: ZeroPage, Cycles: 5},
	0xD6: {Name: "DEC", Dispatch: Dec, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// DEX
	0xCA: {Name: "DEX", Dispatch: Dex, AddressingMode: Implied, Cycles: 2},
	// DEY
	0x88: {Name: "DEY", Dispatch: Dey, AddressingMode: Implied, Cycles: 2},
	// INC
	0xEE: {Name: "INC", Dispatch: Inc, AddressingMode: Absolute, Cycles: 6},
	0xFE: {Name: "INC", Dispatch: Inc, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0xE6: {Name: "INC", Dispatch: Inc, AddressingMode: ZeroPage, Cycles: 5},
	0xF6: {Name: "INC", Dispatch: Inc, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// INX
	0xE8: {Name: "INX", Dispatch: Inx, AddressingMode: Implied, Cycles: 2},
	// INY
	0xC8: {Name: "INY", Dispatch: Iny, AddressingMode: Implied, Cycles: 2},
	// BRK
	0x00: {Name: "BRK", Dispatch: Brk, AddressingMode: Implied, Cycles: 0},
	// JMP
	0x4C: {Name: "JMP", Dispatch: Jmp, AddressingMode: Absolute, Cycles: 3},
	0x6C: {Name: "JMP", Dispatch: Jmp, AddressingMode: AbsoluteIndirect, Cycles: 5},
	// JSR
	0x20: {Name: "JSR", Dispatch: Jsr, AddressingMode: Absolute, Cycles: 6},
	// RTI
	0x40: {Name: "RTI", Dispatch: Rti, AddressingMode: Implied, Cycles: 6},
	// RTS
	0x60: {Name: "RTS", Dispatch: Rts, AddressingMode: Implied, Cycles: 6},
	// BCC
	0x90: {Name: "BCC", Dispatch: Bcc, AddressingMode: Relative, Cycles: 2},
	// BCS
	0xB0: {Name: "BCS", Dispatch: Bcs, AddressingMode: Relative, Cycles: 2},
	// BEQ
	0xF0: {Name: "BEQ", Dispatch: Beq, AddressingMode: Relative, Cycles: 2},
	// BMI
	0x30: {Name: "BMI", Dispatch: Bmi, AddressingMode: Relative, Cycles: 2},
	// BNE
	0xD0: {Name: "BNE", Dispatch: Bne, AddressingMode: Relative, Cycles: 2},
	// BPL
	0x10: {Name: "BPL", Dispatch: Bpl, AddressingMode: Relative, Cycles: 2},
	// BVC
	0x50: {Name: "BVC", Dispatch: Bvc, AddressingMode: Relative, Cycles: 2},
	// BVS
	0x70: {Name: "BVS", Dispatch: Bvs, AddressingMode: Relative, Cycles: 2},
	// CLC
	0x18: {Name: "CLC", Dispatch: Clc, AddressingMode: Implied, Cycles: 2},
	// CLD
	0xD8: {Name: "CLD", Dispatch: Cld, AddressingMode: Implied, Cycles: 2},
	// CLI
	0x58: {Name: "CLI", Dispatch: Cli, AddressingMode: Implied, Cycles: 2},
	// CLV
	0xB8: {Name: "CLV", Dispatch: Clv, AddressingMode: Implied, Cycles: 2},
	// SEC
	0x38: {Name: "SEC", Dispatch: Sec, AddressingMode: Implied, Cycles: 2},
	// SED
	0xF8: {Name: "SED", Dispatch: Sed, AddressingMode: Implied, Cycles: 2},
	// SEI
	0x78: {Name: "SEI", Dispatch: Sei, AddressingMode: Implied, Cycles: 2},
	// NOP 🛑
	0xEA: {Name: "NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0x1A: {Name: "*NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0x3A: {Name: "*NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0x5A: {Name: "*NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0x7A: {Name: "*NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0xDA: {Name: "*NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0xFA: {Name: "*NOP", Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
	0x80: {Name: "*NOP", Dispatch: Nop, AddressingMode: Immediate, Cycles: 2},
	0x82: {Name: "*NOP", Dispatch: Nop, AddressingMode: Immediate, Cycles: 2},
	0x89: {Name: "*NOP", Dispatch: Nop, AddressingMode: Immediate, Cycles: 2},
	0xC2: {Name: "*NOP", Dispatch: Nop, AddressingMode: Immediate, Cycles: 2},
	0xE2: {Name: "*NOP", Dispatch: Nop, AddressingMode: Immediate, Cycles: 2},
	0x0C: {Name: "*NOP", Dispatch: Nop, AddressingMode: Absolute, Cycles: 4},
	0x1C: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x3C: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x5C: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x7C: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xDC: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xFC: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x04: {Name: "*NOP", Dispatch: Nop, AddressingMode: ZeroPage, Cycles: 3},
	0x44: {Name: "*NOP", Dispatch: Nop, AddressingMode: ZeroPage, Cycles: 3},
	0x64: {Name: "*NOP", Dispatch: Nop, AddressingMode: ZeroPage, Cycles: 3},
	0x14: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x34: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x54: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x74: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0xD4: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0xF4: {Name: "*NOP", Dispatch: Nop, AddressingMode: XIndexedZeroPage, Cycles: 4},
}
