package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	changes := lib.InputInts("2018/1")

	// Part 1: Sum of all changes.
	fmt.Println(lib.Sum(changes...))

	// Part 2: First value seen twice.
	var val int
	seen := map[int]struct{}{0: struct{}{}}
Loop:
	for {
		for _, ch := range changes {
			val += ch
			if _, ok := seen[val]; ok {
				fmt.Println(val)
				break Loop
			}
			seen[val] = struct{}{}
		}
	}
}
