package cpu_test

import (
	"testing"

	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/stretchr/testify/require"
)

func TestBccInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when C flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x8000,
		},
		{
			name:        "branch when C flag is reset",
			cpu:         cpu.CPU{Pc: 0x8001},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x9000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagCarry, test.flag)
			cpu.Bcc(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBcsInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when C flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x9000,
		},
		{
			name:        "branch when C flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x8000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagCarry, test.flag)
			cpu.Bcs(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBeqInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when Z flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x9000,
		},
		{
			name:        "branch when Z flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x8000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagZero, test.flag)
			cpu.Beq(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBmiInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when N flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x9000,
		},
		{
			name:        "branch when N flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x8000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagNegative, test.flag)
			cpu.Bmi(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBneInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when Z flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x8000,
		},
		{
			name:        "branch when Z flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x9000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagZero, test.flag)
			cpu.Bne(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBplInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when N flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x8000,
		},
		{
			name:        "branch when N flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x9000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagNegative, test.flag)
			cpu.Bpl(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBvcInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when V flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x8000,
		},
		{
			name:        "branch when V flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x9000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagOverflow, test.flag)
			cpu.Bvc(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}

func TestBvsInstruction(t *testing.T) {
	tests := []struct {
		name        string
		cpu         cpu.CPU
		flag        bool
		value       uint16
		wantAddress uint16
	}{
		{
			name:        "branch when V flag is set",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        true,
			value:       0x9000,
			wantAddress: 0x9000,
		},
		{
			name:        "branch when V flag is reset",
			cpu:         cpu.CPU{Pc: 0x8000},
			flag:        false,
			value:       0x9000,
			wantAddress: 0x8000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagOverflow, test.flag)
			cpu.Bvs(&test.cpu, test.value)
			require.Equal(t, test.wantAddress, test.cpu.Pc)
		})
	}
}
