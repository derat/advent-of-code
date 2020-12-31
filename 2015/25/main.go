package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2015/25")
	lib.AssertEq(len(input), 2)
	r, c := input[0], input[1]

	const (
		start = 20151125
		mult  = 252533
		mod   = 33554393
	)

	pos := r * c
	for i := 1; i <= r-2; i++ {
		pos += i
	}
	for i := 1; i <= c-1; i++ {
		pos += i
	}

	n := int64(start)
	for i := 0; i < pos-1; i++ {
		n = (n * mult) % mod
	}
	fmt.Println(n)
}
