package cpu_test

import (
	"testing"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/stretchr/testify/require"
)

func TestAslInstruction(t *testing.T) {
	tests := []struct {
		name       string
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantNFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with most significant bit set",
			inputValue: 0b10000000,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantNFlag:  false,
			wantZFlag:  true,
		},
		{
			name:       "test negative result",
			inputValue: 0b01100000,
			wantValue:  0b11000000,
			wantCFlag:  false,
			wantNFlag:  true,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		bus := bus.NewBus(nil, nil, nil, nil)
		c := cpu.NewCPU(bus)
		addr := uint16(0x1000)

		c.BusWrite(addr, test.inputValue)
		cpu.Asl(c, addr)
		require.Equal(t, test.wantValue, c.BusRead(addr))
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, test.wantNFlag, c.GetStatusFlag(cpu.StatusFlagNegative))
	}

}

func TestAslAccumulatorInstruction(t *testing.T) {
	tests := []struct {
		name       string
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantNFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with most significant bit set",
			inputValue: 0b10000000,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantNFlag:  false,
			wantZFlag:  true,
		},
		{
			name:       "test negative result",
			inputValue: 0b01100000,
			wantValue:  0b11000000,
			wantCFlag:  false,
			wantNFlag:  true,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		c := &cpu.CPU{A: test.inputValue}

		cpu.AslAccumulator(c, 0)
		require.Equal(t, test.wantValue, c.A)
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, test.wantNFlag, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}

func TestLsrInstruction(t *testing.T) {
	tests := []struct {
		name       string
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with carry and zero value",
			inputValue: 0b00000001,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantZFlag:  true,
		},
		{
			name:       "test with carry",
			inputValue: 0b00000011,
			wantValue:  0b00000001,
			wantCFlag:  true,
			wantZFlag:  false,
		},
		{
			name:       "test with carry without carry",
			inputValue: 0b11000000,
			wantValue:  0b01100000,
			wantCFlag:  false,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		bus := bus.NewBus(nil, nil, nil, nil)
		c := cpu.NewCPU(bus)
		addr := uint16(0x1000)

		c.BusWrite(addr, test.inputValue)
		cpu.Lsr(c, addr)
		require.Equal(t, test.wantValue, c.BusRead(addr))
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.False(t, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}

func TestLsrAccumulatorInstruction(t *testing.T) {
	tests := []struct {
		name       string
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with carry and zero value",
			inputValue: 0b00000001,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantZFlag:  true,
		},
		{
			name:       "test with carry",
			inputValue: 0b00000011,
			wantValue:  0b00000001,
			wantCFlag:  true,
			wantZFlag:  false,
		},
		{
			name:       "test with carry without carry",
			inputValue: 0b11000000,
			wantValue:  0b01100000,
			wantCFlag:  false,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		c := &cpu.CPU{A: test.inputValue, P: 0b10000000}

		cpu.LsrAccumulator(c, 0)
		require.Equal(t, test.wantValue, c.A)
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, false, c.GetStatusFlag(cpu.StatusFlagNegative))
		require.False(t, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}

func TestRolInstruction(t *testing.T) {
	tests := []struct {
		name       string
		carryFlag  bool
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantNFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with most significant bit set",
			carryFlag:  false,
			inputValue: 0b10000000,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantNFlag:  false,
			wantZFlag:  true,
		},
		{
			name:       "test negative result",
			carryFlag:  false,
			inputValue: 0b01100000,
			wantValue:  0b11000000,
			wantCFlag:  false,
			wantNFlag:  true,
			wantZFlag:  false,
		},
		{
			name:       "test rotation with carry set",
			carryFlag:  true,
			inputValue: 0b01100000,
			wantValue:  0b11000001,
			wantCFlag:  false,
			wantNFlag:  true,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		bus := bus.NewBus(nil, nil, nil, nil)
		c := cpu.NewCPU(bus)
		addr := uint16(0x1000)

		c.BusWrite(addr, test.inputValue)
		c.SetStatusFlag(cpu.StatusFlagCarry, test.carryFlag)
		cpu.Rol(c, addr)
		require.Equal(t, test.wantValue, c.BusRead(addr))
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, test.wantNFlag, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}

func TestRolAccumulatorInstruction(t *testing.T) {
	tests := []struct {
		name       string
		carryFlag  bool
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantNFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with most significant bit set",
			carryFlag:  false,
			inputValue: 0b10000000,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantNFlag:  false,
			wantZFlag:  true,
		},
		{
			name:       "test negative result",
			carryFlag:  false,
			inputValue: 0b01100000,
			wantValue:  0b11000000,
			wantCFlag:  false,
			wantNFlag:  true,
			wantZFlag:  false,
		},
		{
			name:       "test rotation with carry set",
			carryFlag:  true,
			inputValue: 0b01100000,
			wantValue:  0b11000001,
			wantCFlag:  false,
			wantNFlag:  true,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		c := &cpu.CPU{A: test.inputValue}
		c.SetStatusFlag(cpu.StatusFlagCarry, test.carryFlag)

		cpu.RolAccumulator(c, 0)
		require.Equal(t, test.wantValue, c.A)
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, test.wantNFlag, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}

func TestRorInstruction(t *testing.T) {
	tests := []struct {
		name       string
		carryFlag  bool
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantNFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with least significant bit set",
			carryFlag:  false,
			inputValue: 0b00000001,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantNFlag:  false,
			wantZFlag:  true,
		},
		{
			name:       "test rotation with carry set",
			carryFlag:  true,
			inputValue: 0b11000001,
			wantValue:  0b11100000,
			wantCFlag:  true,
			wantNFlag:  true,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		bus := bus.NewBus(nil, nil, nil, nil)
		c := cpu.NewCPU(bus)
		addr := uint16(0x1000)

		c.BusWrite(addr, test.inputValue)
		c.SetStatusFlag(cpu.StatusFlagCarry, test.carryFlag)
		cpu.Ror(c, addr)
		require.Equal(t, test.wantValue, c.BusRead(addr))
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, test.wantNFlag, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}

func TestRorAccumulatorInstruction(t *testing.T) {
	tests := []struct {
		name       string
		carryFlag  bool
		inputValue uint8
		wantValue  uint8
		wantCFlag  bool
		wantNFlag  bool
		wantZFlag  bool
	}{
		{
			name:       "test with least significant bit set",
			carryFlag:  false,
			inputValue: 0b00000001,
			wantValue:  0b00000000,
			wantCFlag:  true,
			wantNFlag:  false,
			wantZFlag:  true,
		},
		{
			name:       "test rotation with carry set",
			carryFlag:  true,
			inputValue: 0b11000001,
			wantValue:  0b11100000,
			wantCFlag:  true,
			wantNFlag:  true,
			wantZFlag:  false,
		},
	}

	for _, test := range tests {
		c := &cpu.CPU{A: test.inputValue}
		c.SetStatusFlag(cpu.StatusFlagCarry, test.carryFlag)

		cpu.RorAccumulator(c, 0)
		require.Equal(t, test.wantValue, c.A)
		require.Equal(t, test.wantCFlag, c.GetStatusFlag(cpu.StatusFlagCarry))
		require.Equal(t, test.wantZFlag, c.GetStatusFlag(cpu.StatusFlagZero))
		require.Equal(t, test.wantNFlag, c.GetStatusFlag(cpu.StatusFlagNegative))
	}
}
