package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum int
	for _, m := range lib.InputInts("2019/1") {
		sum += m/3 - 2
	}
	fmt.Println(sum)
}
