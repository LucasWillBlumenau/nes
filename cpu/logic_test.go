package cpu_test

import (
	"testing"

	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/stretchr/testify/require"
)

func TestAndInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantValue uint8
		wantZflag bool
		wantNflag bool
	}{
		{
			name:      "test when all flags should be reset",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b00010110,
			wantValue: 0b00000010,
			wantZflag: false,
			wantNflag: false,
		},
		{
			name:      "test when z flag should be set",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b00010100,
			wantValue: 0b00000000,
			wantZflag: true,
			wantNflag: false,
		},
		{
			name:      "test when n flag should be set",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b10010110,
			wantValue: 0b10000010,
			wantZflag: false,
			wantNflag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.And(&test.cpu, test.value)
			require.Equal(t, test.wantValue, test.cpu.A)
			require.Equal(t, test.wantZflag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNflag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
		})
	}

}

func TestBitInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantValue uint8
		wantZflag bool
		wantNflag bool
		wantVflag bool
	}{
		{
			name:      "test when all flags should be reset",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b00010110,
			wantValue: 0b10101010,
			wantZflag: false,
			wantNflag: false,
			wantVflag: false,
		},
		{
			name:      "test when z flag should be set",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b00010100,
			wantValue: 0b10101010,
			wantZflag: true,
			wantNflag: false,
			wantVflag: false,
		},
		{
			name:      "test when n flag should be set",
			cpu:       cpu.CPU{A: 0b00101010},
			value:     0b10010110,
			wantValue: 0b00101010,
			wantZflag: false,
			wantNflag: true,
			wantVflag: false,
		},
		{
			name:      "test when v flag should be set",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b01010110,
			wantValue: 0b10101010,
			wantZflag: false,
			wantNflag: false,
			wantVflag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Bit(&test.cpu, test.value)
			require.Equal(t, test.wantValue, test.cpu.A)
			require.Equal(t, test.wantZflag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNflag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
			require.Equal(t, test.wantVflag, test.cpu.GetStatusFlag(cpu.StatusFlagOverflow))
		})
	}
}

func TestEorInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantValue uint8
		wantZflag bool
		wantNflag bool
	}{
		{
			name:      "test when n should be set",
			cpu:       cpu.CPU{A: 0b10010010},
			value:     0b01111111,
			wantValue: 0b11101101,
			wantZflag: false,
			wantNflag: true,
		},
		{
			name:      "test when z flag should be set",
			cpu:       cpu.CPU{A: 0b01010101},
			value:     0b01010101,
			wantValue: 0b00000000,
			wantZflag: true,
			wantNflag: false,
		},
		{
			name:      "test set all bits",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b01010101,
			wantValue: 0b11111111,
			wantZflag: false,
			wantNflag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Eor(&test.cpu, test.value)
			require.Equal(t, test.wantValue, test.cpu.A)
			require.Equal(t, test.wantZflag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNflag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
		})
	}
}

func TestOraInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantValue uint8
		wantZflag bool
		wantNflag bool
	}{
		{
			name:      "test when n should be set",
			cpu:       cpu.CPU{A: 0b10000000},
			value:     0b00000000,
			wantValue: 0b10000000,
			wantZflag: false,
			wantNflag: true,
		},
		{
			name:      "test when z flag should be set",
			cpu:       cpu.CPU{A: 0b00000000},
			value:     0b00000000,
			wantValue: 0b00000000,
			wantZflag: true,
			wantNflag: false,
		},
		{
			name:      "test set all bits",
			cpu:       cpu.CPU{A: 0b10101010},
			value:     0b01010101,
			wantValue: 0b11111111,
			wantZflag: false,
			wantNflag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Ora(&test.cpu, test.value)
			require.Equal(t, test.wantValue, test.cpu.A)
			require.Equal(t, test.wantZflag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNflag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
		})
	}
}
