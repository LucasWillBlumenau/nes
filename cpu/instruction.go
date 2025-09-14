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
		return "Accumulator"
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

type StatusFlag uint8

const (
	StatusFlagCarry            = 0
	StatusFlagZero             = 1
	StatusFlagInterruptDisable = 2
	StatusFlagDecimal          = 3
	StatusFlagOverflow         = 6
	StatusFlagNegative         = 7
)

type Instruction struct {
	AddressingMode AddressingMode
	Dispatch       func(*CPU, uint16)
	Cycles         uint8
}

var instructionMap = [256]*Instruction{
	0xA9: {Dispatch: Lda, AddressingMode: Immediate, Cycles: 2},
	0xAD: {Dispatch: Lda, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBD: {Dispatch: Lda, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xB9: {Dispatch: Lda, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xA5: {Dispatch: Lda, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB5: {Dispatch: Lda, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0xA1: {Dispatch: Lda, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xB1: {Dispatch: Lda, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// LDX
	0xA2: {Dispatch: Ldx, AddressingMode: Immediate, Cycles: 2},
	0xAE: {Dispatch: Ldx, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBE: {Dispatch: Ldx, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xA6: {Dispatch: Ldx, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB6: {Dispatch: Ldx, AddressingMode: YIndexedZeroPageValue, Cycles: 4},
	// LDY
	0xA0: {Dispatch: Ldy, AddressingMode: Immediate, Cycles: 2},
	0xAC: {Dispatch: Ldy, AddressingMode: AbsoluteValue, Cycles: 4},
	0xBC: {Dispatch: Ldy, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xA4: {Dispatch: Ldy, AddressingMode: ZeroPageValue, Cycles: 3},
	0xB4: {Dispatch: Ldy, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	// STA
	0x8D: {Dispatch: Sta, AddressingMode: Absolute, Cycles: 4},
	0x9D: {Dispatch: Sta, AddressingMode: XIndexedAbsolute, Cycles: 5},
	0x99: {Dispatch: Sta, AddressingMode: YIndexedAbsolute, Cycles: 5},
	0x85: {Dispatch: Sta, AddressingMode: ZeroPage, Cycles: 3},
	0x95: {Dispatch: Sta, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x81: {Dispatch: Sta, AddressingMode: XIndexedZeroPageIndirect, Cycles: 6},
	0x91: {Dispatch: Sta, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 6},
	// STX
	0x8E: {Dispatch: Stx, AddressingMode: Absolute, Cycles: 4},
	0x86: {Dispatch: Stx, AddressingMode: ZeroPage, Cycles: 3},
	0x96: {Dispatch: Stx, AddressingMode: YIndexedZeroPage, Cycles: 4},
	// STY
	0x8C: {Dispatch: Sty, AddressingMode: Absolute, Cycles: 4},
	0x84: {Dispatch: Sty, AddressingMode: ZeroPage, Cycles: 3},
	0x94: {Dispatch: Sty, AddressingMode: XIndexedZeroPage, Cycles: 4},
	// TAX
	0xAA: {Dispatch: Tax, AddressingMode: Implied, Cycles: 2},
	// TAY
	0xA8: {Dispatch: Tay, AddressingMode: Implied, Cycles: 2},
	// TSX
	0xBA: {Dispatch: Tsx, AddressingMode: Implied, Cycles: 2},
	// TXA
	0x8A: {Dispatch: Txa, AddressingMode: Implied, Cycles: 2},
	// TXS
	0x9A: {Dispatch: Txs, AddressingMode: Implied, Cycles: 2},
	// TYA
	0x98: {Dispatch: Tya, AddressingMode: Implied, Cycles: 2},
	// PHA
	0x48: {Dispatch: Pha, AddressingMode: Implied, Cycles: 3},
	// PHP - 🐘
	0x08: {Dispatch: Php, AddressingMode: Implied, Cycles: 3},
	// PLA
	0x68: {Dispatch: Pla, AddressingMode: Implied, Cycles: 4},
	// PLP
	0x28: {Dispatch: Plp, AddressingMode: Implied, Cycles: 4},
	// ASL
	0x0A: {Dispatch: AslAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x0E: {Dispatch: Asl, AddressingMode: Absolute, Cycles: 6},
	0x1E: {Dispatch: Asl, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x06: {Dispatch: Asl, AddressingMode: ZeroPage, Cycles: 5},
	0x16: {Dispatch: Asl, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// LSR
	0x4A: {Dispatch: LsrAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x4E: {Dispatch: Lsr, AddressingMode: Absolute, Cycles: 6},
	0x5E: {Dispatch: Lsr, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x46: {Dispatch: Lsr, AddressingMode: ZeroPage, Cycles: 5},
	0x56: {Dispatch: Lsr, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// ROL
	0x2A: {Dispatch: RolAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x2E: {Dispatch: Rol, AddressingMode: Absolute, Cycles: 6},
	0x3E: {Dispatch: Rol, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x26: {Dispatch: Rol, AddressingMode: ZeroPage, Cycles: 5},
	0x36: {Dispatch: Rol, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// ROR
	0x6A: {Dispatch: RorAccumulator, AddressingMode: Accumulator, Cycles: 2},
	0x6E: {Dispatch: Ror, AddressingMode: Absolute, Cycles: 6},
	0x7E: {Dispatch: Ror, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0x66: {Dispatch: Ror, AddressingMode: ZeroPage, Cycles: 5},
	0x76: {Dispatch: Ror, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// AND
	0x29: {Dispatch: And, AddressingMode: Immediate, Cycles: 2},
	0x2D: {Dispatch: And, AddressingMode: AbsoluteValue, Cycles: 4},
	0x3D: {Dispatch: And, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x39: {Dispatch: And, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0x25: {Dispatch: And, AddressingMode: ZeroPageValue, Cycles: 3},
	0x35: {Dispatch: And, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0x21: {Dispatch: And, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0x31: {Dispatch: And, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// BIT
	0x2C: {Dispatch: Bit, AddressingMode: AbsoluteValue, Cycles: 4},
	0x24: {Dispatch: Bit, AddressingMode: ZeroPageValue, Cycles: 3},
	// EOR
	0x49: {Dispatch: Eor, AddressingMode: Immediate, Cycles: 2},
	0x4D: {Dispatch: Eor, AddressingMode: Absolute, Cycles: 4},
	0x5D: {Dispatch: Eor, AddressingMode: XIndexedAbsolute, Cycles: 4},
	0x59: {Dispatch: Eor, AddressingMode: YIndexedAbsolute, Cycles: 4},
	0x45: {Dispatch: Eor, AddressingMode: ZeroPage, Cycles: 3},
	0x55: {Dispatch: Eor, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x41: {Dispatch: Eor, AddressingMode: XIndexedZeroPageIndirect, Cycles: 6},
	0x51: {Dispatch: Eor, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 5},
	// ORA
	0x09: {Dispatch: Ora, AddressingMode: Immediate, Cycles: 2},
	0x0D: {Dispatch: Ora, AddressingMode: Absolute, Cycles: 4},
	0x1D: {Dispatch: Ora, AddressingMode: XIndexedAbsolute, Cycles: 4},
	0x19: {Dispatch: Ora, AddressingMode: YIndexedAbsolute, Cycles: 4},
	0x05: {Dispatch: Ora, AddressingMode: ZeroPage, Cycles: 3},
	0x15: {Dispatch: Ora, AddressingMode: XIndexedZeroPage, Cycles: 4},
	0x01: {Dispatch: Ora, AddressingMode: XIndexedZeroPageIndirect, Cycles: 6},
	0x11: {Dispatch: Ora, AddressingMode: ZeroPageIndirectYIndexed, Cycles: 5},
	// ADC
	0x69: {Dispatch: Adc, AddressingMode: Immediate, Cycles: 2},
	0x6D: {Dispatch: Adc, AddressingMode: AbsoluteValue, Cycles: 4},
	0x7D: {Dispatch: Adc, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0x79: {Dispatch: Adc, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0x65: {Dispatch: Adc, AddressingMode: ZeroPageValue, Cycles: 3},
	0x75: {Dispatch: Adc, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0x61: {Dispatch: Adc, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0x71: {Dispatch: Adc, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// CMP
	0xC9: {Dispatch: Cmp, AddressingMode: Immediate, Cycles: 2},
	0xCD: {Dispatch: Cmp, AddressingMode: AbsoluteValue, Cycles: 4},
	0xDD: {Dispatch: Cmp, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xD9: {Dispatch: Cmp, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xC5: {Dispatch: Cmp, AddressingMode: ZeroPageValue, Cycles: 3},
	0xD5: {Dispatch: Cmp, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0xC1: {Dispatch: Cmp, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xD1: {Dispatch: Cmp, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// CPX
	0xE0: {Dispatch: Cpx, AddressingMode: Immediate, Cycles: 2},
	0xEC: {Dispatch: Cpx, AddressingMode: AbsoluteValue, Cycles: 4},
	0xE4: {Dispatch: Cpx, AddressingMode: ZeroPageValue, Cycles: 3},
	// CPY
	0xC0: {Dispatch: Cpy, AddressingMode: Immediate, Cycles: 2},
	0xCC: {Dispatch: Cpy, AddressingMode: AbsoluteValue, Cycles: 4},
	0xC4: {Dispatch: Cpy, AddressingMode: ZeroPageValue, Cycles: 3},
	// SBC
	0xE9: {Dispatch: Sbc, AddressingMode: Immediate, Cycles: 2},
	0xED: {Dispatch: Sbc, AddressingMode: AbsoluteValue, Cycles: 4},
	0xFD: {Dispatch: Sbc, AddressingMode: XIndexedAbsoluteValue, Cycles: 4},
	0xF9: {Dispatch: Sbc, AddressingMode: YIndexedAbsoluteValue, Cycles: 4},
	0xE5: {Dispatch: Sbc, AddressingMode: ZeroPageValue, Cycles: 3},
	0xF5: {Dispatch: Sbc, AddressingMode: XIndexedZeroPageValue, Cycles: 4},
	0xE1: {Dispatch: Sbc, AddressingMode: XIndexedZeroPageIndirectValue, Cycles: 6},
	0xF1: {Dispatch: Sbc, AddressingMode: ZeroPageIndirectYIndexedValue, Cycles: 5},
	// DEC
	0xCE: {Dispatch: Dec, AddressingMode: Absolute, Cycles: 6},
	0xDE: {Dispatch: Dec, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0xC6: {Dispatch: Dec, AddressingMode: ZeroPage, Cycles: 5},
	0xD6: {Dispatch: Dec, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// DEX
	0xCA: {Dispatch: Dex, AddressingMode: Implied, Cycles: 2},
	// DEY
	0x88: {Dispatch: Dey, AddressingMode: Implied, Cycles: 2},
	// INC
	0xEE: {Dispatch: Inc, AddressingMode: Absolute, Cycles: 6},
	0xFE: {Dispatch: Inc, AddressingMode: XIndexedAbsolute, Cycles: 7},
	0xE6: {Dispatch: Inc, AddressingMode: ZeroPage, Cycles: 5},
	0xF6: {Dispatch: Inc, AddressingMode: XIndexedZeroPage, Cycles: 6},
	// INX
	0xE8: {Dispatch: Inx, AddressingMode: Implied, Cycles: 2},
	// INY
	0xC8: {Dispatch: Iny, AddressingMode: Implied, Cycles: 2},
	// BRK
	0x00: {Dispatch: Brk, AddressingMode: Implied, Cycles: 7}, // TODO: check clock logic implementation for this instruction later
	// JMP
	0x4C: {Dispatch: Jmp, AddressingMode: Absolute, Cycles: 3},
	0x6C: {Dispatch: Jmp, AddressingMode: AbsoluteIndirect, Cycles: 5},
	// JSR
	0x20: {Dispatch: Jsr, AddressingMode: Absolute, Cycles: 6},
	// RTI
	0x40: {Dispatch: Rti, AddressingMode: Implied, Cycles: 6},
	// RTS
	0x60: {Dispatch: Rts, AddressingMode: Implied, Cycles: 6},
	// BCC
	0x90: {Dispatch: Bcc, AddressingMode: Relative, Cycles: 2},
	// BCS
	0xB0: {Dispatch: Bcs, AddressingMode: Relative, Cycles: 2},
	// BEQ
	0xF0: {Dispatch: Beq, AddressingMode: Relative, Cycles: 2},
	// BMI
	0x30: {Dispatch: Bmi, AddressingMode: Relative, Cycles: 2},
	// BNE
	0xD0: {Dispatch: Bne, AddressingMode: Relative, Cycles: 2},
	// BPL
	0x10: {Dispatch: Bpl, AddressingMode: Relative, Cycles: 2},
	// BVC
	0x50: {Dispatch: Bvc, AddressingMode: Relative, Cycles: 2},
	// BVS
	0x70: {Dispatch: Bvs, AddressingMode: Relative, Cycles: 2},
	// CLC
	0x18: {Dispatch: Clc, AddressingMode: Implied, Cycles: 2},
	// CLD
	0xD8: {Dispatch: Cld, AddressingMode: Implied, Cycles: 2},
	// CLI
	0x58: {Dispatch: Cli, AddressingMode: Implied, Cycles: 2},
	// CLV
	0xB8: {Dispatch: Clv, AddressingMode: Implied, Cycles: 2},
	// SEC
	0x38: {Dispatch: Sec, AddressingMode: Implied, Cycles: 2},
	// SED
	0xF8: {Dispatch: Sed, AddressingMode: Implied, Cycles: 2},
	// SEI
	0x78: {Dispatch: Sei, AddressingMode: Implied, Cycles: 2},
	// NOP 🛑
	0xEA: {Dispatch: Nop, AddressingMode: Implied, Cycles: 2},
}
