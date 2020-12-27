package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var floor, basement int
	for i, ch := range lib.Line('(', ')') {
		switch ch {
		case '(':
			floor++
		case ')':
			floor--
		}
		if floor < 0 && basement == 0 {
			basement = i + 1
		}
	}
	fmt.Println(floor)
	fmt.Println(basement)
}
