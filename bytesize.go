package main

import (
	"fmt"
)

func main() {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	fmt.Printf("KB=%.fBytes\n", KB)
	fmt.Printf("MB=%.fBytes\n", MB)
	fmt.Printf("GB=%.fBytes\n", GB)
	fmt.Printf("TB=%.fBytes\n", TB)
	fmt.Printf("PB=%.fBytes\n", PB)
}
