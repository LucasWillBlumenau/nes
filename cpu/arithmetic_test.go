package cpu_test

import (
	"testing"

	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/stretchr/testify/require"
)

func TestAdcInstruction(t *testing.T) {

	tests := []struct {
		name         string
		cpu          cpu.CPU
		value        uint16
		wantValue    uint8
		wantCarry    bool
		wantNegative bool
		wantZero     bool
		wantOverflow bool
	}{
		{
			name:         "Test sum",
			cpu:          cpu.CPU{A: 10},
			value:        100,
			wantValue:    110,
			wantCarry:    false,
			wantNegative: false,
			wantZero:     false,
			wantOverflow: false,
		},
		{
			name:         "Test sum with negative value caused by overflow",
			cpu:          cpu.CPU{A: 127},
			value:        30,
			wantValue:    157, // -99
			wantCarry:    false,
			wantNegative: true,
			wantZero:     false,
			wantOverflow: true,
		},
		{
			name:         "Test sum with positive value cause by overflow",
			cpu:          cpu.CPU{A: 254}, // -2
			value:        129,             // -127
			wantValue:    127,             // 127
			wantCarry:    true,
			wantNegative: false,
			wantZero:     false,
			wantOverflow: true,
		},
		{
			name:         "Test sum with carry by overflow",
			cpu:          cpu.CPU{A: 200},
			value:        56,
			wantValue:    0,
			wantCarry:    true,
			wantNegative: false,
			wantZero:     true,
			wantOverflow: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Adc(&test.cpu, test.value)
			require.Equal(t, test.cpu.A, test.wantValue)
			require.Equal(t, test.wantCarry, test.cpu.GetStatusFlag(cpu.StatusFlagCarry))
			require.Equal(t, test.wantNegative, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
			require.Equal(t, test.wantZero, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantOverflow, test.cpu.GetStatusFlag(cpu.StatusFlagOverflow))
		})

	}

}

func TestSbcInstruction(t *testing.T) {

	tests := []struct {
		name         string
		cpu          cpu.CPU
		value        uint16
		carry        bool
		wantValue    uint8
		wantCarry    bool
		wantNegative bool
		wantZero     bool
		wantOverflow bool
	}{
		{
			name:         "Test difference from two positive numbers",
			cpu:          cpu.CPU{A: 10},
			value:        5,
			wantValue:    5,
			carry:        true,
			wantCarry:    true,
			wantNegative: false,
			wantZero:     false,
			wantOverflow: false,
		},
		{
			name:         "Test difference from a negative number",
			cpu:          cpu.CPU{A: 10},
			value:        255,
			wantValue:    11,
			carry:        true,
			wantCarry:    false,
			wantNegative: false,
			wantZero:     false,
			wantOverflow: false,
		},
		{
			name:         "Test difference with overflow",
			cpu:          cpu.CPU{A: 128},
			value:        1,
			wantValue:    127,
			carry:        true,
			wantCarry:    true,
			wantNegative: false,
			wantZero:     false,
			wantOverflow: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cpu.SetStatusFlag(cpu.StatusFlagCarry, test.carry)
			cpu.Sbc(&test.cpu, test.value)
			require.Equal(t, test.cpu.A, test.wantValue)
			require.Equal(t, test.wantCarry, test.cpu.GetStatusFlag(cpu.StatusFlagCarry))
			require.Equal(t, test.wantNegative, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
			require.Equal(t, test.wantZero, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantOverflow, test.cpu.GetStatusFlag(cpu.StatusFlagOverflow))
		})
	}

}

func TestCmpInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantZFlag bool
		wantCFlag bool
		wantNFlag bool
	}{
		{
			name:      "Test comparasion with equal number",
			cpu:       cpu.CPU{A: 68},
			value:     68,
			wantZFlag: true,
			wantNFlag: false,
			wantCFlag: true,
		},
		{
			name:      "Test comparasion with greater number",
			cpu:       cpu.CPU{A: 68},
			value:     70,
			wantZFlag: false,
			wantNFlag: true,
			wantCFlag: false,
		},
		{
			name:      "Test comparasion with smaller number",
			cpu:       cpu.CPU{A: 68},
			value:     67,
			wantZFlag: false,
			wantNFlag: false,
			wantCFlag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Cmp(&test.cpu, test.value)
			require.Equal(t, test.wantCFlag, test.cpu.GetStatusFlag(cpu.StatusFlagCarry))
			require.Equal(t, test.wantZFlag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNFlag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
		})
	}
}

func TestCpxInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantZFlag bool
		wantCFlag bool
		wantNFlag bool
	}{
		{
			name:      "Test comparasion with equal number",
			cpu:       cpu.CPU{X: 68},
			value:     68,
			wantZFlag: true,
			wantNFlag: false,
			wantCFlag: true,
		},
		{
			name:      "Test comparasion with greater number",
			cpu:       cpu.CPU{X: 68},
			value:     70,
			wantZFlag: false,
			wantNFlag: true,
			wantCFlag: false,
		},
		{
			name:      "Test comparasion with smaller number",
			cpu:       cpu.CPU{X: 68},
			value:     67,
			wantZFlag: false,
			wantNFlag: false,
			wantCFlag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Cpx(&test.cpu, test.value)
			require.Equal(t, test.wantCFlag, test.cpu.GetStatusFlag(cpu.StatusFlagCarry))
			require.Equal(t, test.wantZFlag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNFlag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
		})
	}
}

func TestCpyInstruction(t *testing.T) {
	tests := []struct {
		name      string
		cpu       cpu.CPU
		value     uint16
		wantZFlag bool
		wantCFlag bool
		wantNFlag bool
	}{
		{
			name:      "Test comparasion with equal number",
			cpu:       cpu.CPU{Y: 68},
			value:     68,
			wantZFlag: true,
			wantNFlag: false,
			wantCFlag: true,
		},
		{
			name:      "Test comparasion with greater number",
			cpu:       cpu.CPU{Y: 68},
			value:     70,
			wantZFlag: false,
			wantNFlag: true,
			wantCFlag: false,
		},
		{
			name:      "Test comparasion with smaller number",
			cpu:       cpu.CPU{Y: 68},
			value:     67,
			wantZFlag: false,
			wantNFlag: false,
			wantCFlag: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu.Cpy(&test.cpu, test.value)
			require.Equal(t, test.wantCFlag, test.cpu.GetStatusFlag(cpu.StatusFlagCarry))
			require.Equal(t, test.wantZFlag, test.cpu.GetStatusFlag(cpu.StatusFlagZero))
			require.Equal(t, test.wantNFlag, test.cpu.GetStatusFlag(cpu.StatusFlagNegative))
		})
	}
}
