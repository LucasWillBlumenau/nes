package cpu

var instructionMap = [256]*Instruction{
	0xA9: {Dispatch: Lda, AddressingMode: Immediate},
	0xAD: {Dispatch: Lda, AddressingMode: AbsoluteValue},
	0xBD: {Dispatch: Lda, AddressingMode: XIndexedAbsoluteValue},
	0xB9: {Dispatch: Lda, AddressingMode: YIndexedAbsoluteValue},
	0xA5: {Dispatch: Lda, AddressingMode: ZeroPageValue},
	0xB5: {Dispatch: Lda, AddressingMode: XIndexedZeroPageValue},
	0xA1: {Dispatch: Lda, AddressingMode: XIndexedZeroPageIndirectValue},
	0xB1: {Dispatch: Lda, AddressingMode: ZeroPageIndirectYIndexedValue},
	// LDX
	0xA2: {Dispatch: Ldx, AddressingMode: Immediate},
	0xAE: {Dispatch: Ldx, AddressingMode: AbsoluteValue},
	0xBE: {Dispatch: Ldx, AddressingMode: YIndexedAbsoluteValue},
	0xA6: {Dispatch: Ldx, AddressingMode: ZeroPageValue},
	0xB6: {Dispatch: Ldx, AddressingMode: YIndexedZeroPageValue},
	// LDY
	0xA0: {Dispatch: Ldy, AddressingMode: Immediate},
	0xAC: {Dispatch: Ldy, AddressingMode: AbsoluteValue},
	0xBC: {Dispatch: Ldy, AddressingMode: XIndexedAbsoluteValue},
	0xA4: {Dispatch: Ldy, AddressingMode: ZeroPageValue},
	0xB4: {Dispatch: Ldy, AddressingMode: XIndexedZeroPageValue},
	// STA
	0x8D: {Dispatch: Sta, AddressingMode: Absolute},
	0x9D: {Dispatch: Sta, AddressingMode: XIndexedAbsolute},
	0x99: {Dispatch: Sta, AddressingMode: YIndexedAbsolute},
	0x85: {Dispatch: Sta, AddressingMode: ZeroPage},
	0x95: {Dispatch: Sta, AddressingMode: XIndexedZeroPage},
	0x81: {Dispatch: Sta, AddressingMode: XIndexedZeroPageIndirect},
	0x91: {Dispatch: Sta, AddressingMode: ZeroPageIndirectYIndexed},
	// STX
	0x8E: {Dispatch: Stx, AddressingMode: Absolute},
	0x86: {Dispatch: Stx, AddressingMode: ZeroPage},
	0x96: {Dispatch: Stx, AddressingMode: YIndexedZeroPage},
	// STY
	0x8C: {Dispatch: Sty, AddressingMode: Absolute},
	0x84: {Dispatch: Sty, AddressingMode: ZeroPage},
	0x94: {Dispatch: Sty, AddressingMode: XIndexedZeroPage},
	// TAX
	0xAA: {Dispatch: Tax, AddressingMode: Implied},
	// TAY
	0xA8: {Dispatch: Tay, AddressingMode: Implied},
	// TSX
	0xBA: {Dispatch: Tsx, AddressingMode: Implied},
	// TXA
	0x8A: {Dispatch: Txa, AddressingMode: Implied},
	// TXS
	0x9A: {Dispatch: Txs, AddressingMode: Implied},
	// TYA
	0x98: {Dispatch: Tya, AddressingMode: Implied},
	// PHA
	0x48: {Dispatch: Pha, AddressingMode: Implied},
	// PHP - 🐘
	0x08: {Dispatch: Php, AddressingMode: Implied},
	// PLA
	0x68: {Dispatch: Pla, AddressingMode: Implied},
	// PLP
	0x28: {Dispatch: Plp, AddressingMode: Implied},
	// ASL
	0x0A: {Dispatch: Asl, AddressingMode: Accumulator},
	0x0E: {Dispatch: Asl, AddressingMode: Absolute},
	0x1E: {Dispatch: Asl, AddressingMode: XIndexedAbsolute},
	0x06: {Dispatch: Asl, AddressingMode: ZeroPage},
	0x16: {Dispatch: Asl, AddressingMode: XIndexedZeroPage},
	// LSR
	0x4A: {Dispatch: Lsr, AddressingMode: Accumulator},
	0x4E: {Dispatch: Lsr, AddressingMode: Absolute},
	0x5E: {Dispatch: Lsr, AddressingMode: XIndexedAbsolute},
	0x46: {Dispatch: Lsr, AddressingMode: ZeroPage},
	0x56: {Dispatch: Lsr, AddressingMode: XIndexedZeroPage},
	// ROL
	0x2A: {Dispatch: Rol, AddressingMode: Accumulator},
	0x2E: {Dispatch: Rol, AddressingMode: Absolute},
	0x3E: {Dispatch: Rol, AddressingMode: XIndexedAbsolute},
	0x26: {Dispatch: Rol, AddressingMode: ZeroPage},
	0x36: {Dispatch: Rol, AddressingMode: XIndexedZeroPage},
	// ROR
	0x6A: {Dispatch: Ror, AddressingMode: Accumulator},
	0x6E: {Dispatch: Ror, AddressingMode: Absolute},
	0x7E: {Dispatch: Ror, AddressingMode: XIndexedAbsolute},
	0x66: {Dispatch: Ror, AddressingMode: ZeroPage},
	0x76: {Dispatch: Ror, AddressingMode: XIndexedZeroPage},
	// AND
	0x29: {Dispatch: And, AddressingMode: Immediate},
	0x2D: {Dispatch: And, AddressingMode: AbsoluteValue},
	0x3D: {Dispatch: And, AddressingMode: XIndexedAbsoluteValue},
	0x39: {Dispatch: And, AddressingMode: YIndexedAbsoluteValue},
	0x25: {Dispatch: And, AddressingMode: ZeroPageValue},
	0x35: {Dispatch: And, AddressingMode: XIndexedZeroPageValue},
	0x21: {Dispatch: And, AddressingMode: XIndexedZeroPageIndirectValue},
	0x31: {Dispatch: And, AddressingMode: ZeroPageIndirectYIndexedValue},
	// BIT
	0x2C: {Dispatch: Bit, AddressingMode: AbsoluteValue},
	0x24: {Dispatch: Bit, AddressingMode: ZeroPageValue},
	// EOR
	0x49: {Dispatch: Eor, AddressingMode: Immediate},
	0x4D: {Dispatch: Eor, AddressingMode: Absolute},
	0x5D: {Dispatch: Eor, AddressingMode: XIndexedAbsolute},
	0x59: {Dispatch: Eor, AddressingMode: YIndexedAbsolute},
	0x45: {Dispatch: Eor, AddressingMode: ZeroPage},
	0x55: {Dispatch: Eor, AddressingMode: XIndexedZeroPage},
	0x41: {Dispatch: Eor, AddressingMode: XIndexedZeroPageIndirect},
	0x51: {Dispatch: Eor, AddressingMode: ZeroPageIndirectYIndexed},
	// ORA
	0x09: {Dispatch: Ora, AddressingMode: Immediate},
	0x0D: {Dispatch: Ora, AddressingMode: Absolute},
	0x1D: {Dispatch: Ora, AddressingMode: XIndexedAbsolute},
	0x19: {Dispatch: Ora, AddressingMode: YIndexedAbsolute},
	0x05: {Dispatch: Ora, AddressingMode: ZeroPage},
	0x15: {Dispatch: Ora, AddressingMode: XIndexedZeroPage},
	0x01: {Dispatch: Ora, AddressingMode: XIndexedZeroPageIndirect},
	0x11: {Dispatch: Ora, AddressingMode: ZeroPageIndirectYIndexed},
	// ADC
	0x69: {Dispatch: Adc, AddressingMode: Immediate},
	0x6D: {Dispatch: Adc, AddressingMode: AbsoluteValue},
	0x7D: {Dispatch: Adc, AddressingMode: XIndexedAbsoluteValue},
	0x79: {Dispatch: Adc, AddressingMode: YIndexedAbsoluteValue},
	0x65: {Dispatch: Adc, AddressingMode: ZeroPageValue},
	0x75: {Dispatch: Adc, AddressingMode: XIndexedZeroPageValue},
	0x61: {Dispatch: Adc, AddressingMode: XIndexedZeroPageIndirectValue},
	0x71: {Dispatch: Adc, AddressingMode: ZeroPageIndirectYIndexedValue},
	// CMP
	0xC9: {Dispatch: Cmp, AddressingMode: Immediate},
	0xCD: {Dispatch: Cmp, AddressingMode: AbsoluteValue},
	0xDD: {Dispatch: Cmp, AddressingMode: XIndexedAbsoluteValue},
	0xD9: {Dispatch: Cmp, AddressingMode: YIndexedAbsoluteValue},
	0xC5: {Dispatch: Cmp, AddressingMode: ZeroPageValue},
	0xD5: {Dispatch: Cmp, AddressingMode: XIndexedZeroPageValue},
	0xC1: {Dispatch: Cmp, AddressingMode: XIndexedZeroPageIndirectValue},
	0xD1: {Dispatch: Cmp, AddressingMode: ZeroPageIndirectYIndexedValue},
	// CPX
	0xE0: {Dispatch: Cpx, AddressingMode: Immediate},
	0xEC: {Dispatch: Cpx, AddressingMode: AbsoluteValue},
	0xE4: {Dispatch: Cpx, AddressingMode: ZeroPageValue},
	// CPY
	0xC0: {Dispatch: Cpy, AddressingMode: Immediate},
	0xCC: {Dispatch: Cpy, AddressingMode: AbsoluteValue},
	0xC4: {Dispatch: Cpy, AddressingMode: ZeroPageValue},
	// SBC
	0xE9: {Dispatch: Sbc, AddressingMode: Immediate},
	0xED: {Dispatch: Sbc, AddressingMode: AbsoluteValue},
	0xFD: {Dispatch: Sbc, AddressingMode: XIndexedAbsoluteValue},
	0xF9: {Dispatch: Sbc, AddressingMode: YIndexedAbsoluteValue},
	0xE5: {Dispatch: Sbc, AddressingMode: ZeroPageValue},
	0xF5: {Dispatch: Sbc, AddressingMode: XIndexedZeroPageValue},
	0xE1: {Dispatch: Sbc, AddressingMode: XIndexedZeroPageIndirectValue},
	0xF1: {Dispatch: Sbc, AddressingMode: ZeroPageIndirectYIndexedValue},
	// DEC
	0xCE: {Dispatch: Dec, AddressingMode: Absolute},
	0xDE: {Dispatch: Dec, AddressingMode: XIndexedAbsolute},
	0xC6: {Dispatch: Dec, AddressingMode: ZeroPage},
	0xD6: {Dispatch: Dec, AddressingMode: XIndexedZeroPage},
	// DEX
	0xCA: {Dispatch: Dex, AddressingMode: Implied},
	// DEY
	0x88: {Dispatch: Dey, AddressingMode: Implied},
	// INC
	0xEE: {Dispatch: Inc, AddressingMode: Absolute},
	0xFE: {Dispatch: Inc, AddressingMode: XIndexedAbsolute},
	0xE6: {Dispatch: Inc, AddressingMode: ZeroPage},
	0xF6: {Dispatch: Inc, AddressingMode: XIndexedZeroPage},
	// INX
	0xE8: {Dispatch: Inx, AddressingMode: Implied},
	// INY
	0xC8: {Dispatch: Iny, AddressingMode: Implied},
	// BRK
	0x00: {Dispatch: Brk, AddressingMode: Implied},
	// JMP
	0x4C: {Dispatch: Jmp, AddressingMode: Absolute},
	0x6C: {Dispatch: Jmp, AddressingMode: AbsoluteIndirect},
	// JSR
	0x20: {Dispatch: Jsr, AddressingMode: Absolute},
	// RTI
	0x40: {Dispatch: Rti, AddressingMode: Implied},
	// RTS
	0x60: {Dispatch: Rts, AddressingMode: Implied},
	// BCC
	0x90: {Dispatch: Bcc, AddressingMode: Relative},
	// BCS
	0xB0: {Dispatch: Bcs, AddressingMode: Relative},
	// BEQ
	0xF0: {Dispatch: Beq, AddressingMode: Relative},
	// BMI
	0x30: {Dispatch: Bmi, AddressingMode: Relative},
	// BNE
	0xD0: {Dispatch: Bne, AddressingMode: Relative},
	// BPL
	0x10: {Dispatch: Bpl, AddressingMode: Relative},
	// BVC
	0x50: {Dispatch: Bvc, AddressingMode: Relative},
	// BVS
	0x70: {Dispatch: Bvs, AddressingMode: Relative},
	// CLC
	0x18: {Dispatch: Clc, AddressingMode: Implied},
	// CLD
	0xD8: {Dispatch: Cld, AddressingMode: Implied},
	// CLI
	0x58: {Dispatch: Cli, AddressingMode: Implied},
	// CLV
	0xB8: {Dispatch: Clv, AddressingMode: Implied},
	// SEC
	0x38: {Dispatch: Sec, AddressingMode: Implied},
	// SED
	0xF8: {Dispatch: Sed, AddressingMode: Implied},
	// SEI
	0x78: {Dispatch: Sei, AddressingMode: Implied},
	// NOP 🛑
	0xEA: {Dispatch: Nop, AddressingMode: Implied},
}

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
}
