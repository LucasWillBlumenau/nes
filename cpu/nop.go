package cpu

import (
	"fmt"
)

func Nop(_ *CPU, _ uint16) {
	fmt.Println("executing instruction NOP...")
}
