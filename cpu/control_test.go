package cpu_test

import (
	"testing"

	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/stretchr/testify/require"
)

func TestJmpInstruction(t *testing.T) {
	c := cpu.CPU{Pc: 0x8000}
	cpu.Jmp(&c, 0x9000)
	require.Equal(t, uint16(0x9000), c.Pc)
}
