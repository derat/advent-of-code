package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	floor := 0
	for _, ch := range lib.Line('(', ')') {
		switch ch {
		case '(':
			floor++
		case ')':
			floor--
		}
	}
	fmt.Println(floor)
}
