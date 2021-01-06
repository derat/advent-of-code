package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	digits := lib.ExtractDigits(lib.Input("2017/1"))
	cnt := len(digits)
	var sum, sum2 int
	for i, n := range digits {
		sum += lib.If(n == digits[(i+1)%cnt], n, 0)
		sum2 += lib.If(n == digits[(i+cnt/2)%cnt], n, 0)
	}
	fmt.Println(sum)
	fmt.Println(sum2)
}
