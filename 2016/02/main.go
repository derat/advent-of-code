package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	r, c := 1, 1
	for _, ln := range lib.InputLines("2016/2") {
		for _, ch := range ln {
			switch ch {
			case 'U':
				r = lib.Max(r-1, 0)
			case 'D':
				r = lib.Min(r+1, 2)
			case 'L':
				c = lib.Max(c-1, 0)
			case 'R':
				c = lib.Min(c+1, 2)
			default:
				panic(fmt.Sprintf("Bad char %q", ch))
			}
		}
		fmt.Print(r*3 + c + 1)
	}
	fmt.Println()
}
