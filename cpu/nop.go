package cpu

import (
	"fmt"
)

func Nop(_ *CPU, _ uint16) {
	for i := 1; i < 100; i++ {
		fmt.Print("")
	}
	fmt.Println("...")
}
