package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.Input("2022/6")
	for _, size := range []int{4, 14} {
		counts := make(map[byte]int)
		start := -1
		for i := range input {
			counts[input[i]]++
			if i >= size {
				if old := input[i-size]; counts[old] > 1 {
					counts[old]--
				} else {
					delete(counts, old)
				}
			}
			if len(counts) == size {
				start = i
				break
			}
		}
		lib.Assert(start >= 0)
		fmt.Println(start + 1)
	}
}
