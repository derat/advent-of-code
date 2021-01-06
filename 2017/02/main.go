package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum int
	for _, ln := range lib.InputLines("2017/2") {
		nums := lib.ExtractInts(ln)
		sum += lib.Max(nums...) - lib.Min(nums...)
	}
	fmt.Println(sum)
}
