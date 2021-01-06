package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum int
	digits := lib.ExtractDigits(lib.Input("2017/1"))
	for i, n := range digits {
		sum += lib.If(n == digits[(i+1)%len(digits)], n, 0)
	}
	fmt.Println(sum)
}
