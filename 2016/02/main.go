package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Part 1:
	//  1 2 3
	//  4 5 6
	//  7 8 9
	r, c := 1, 1
	clamp := func(a int) int { return lib.Clamp(a, 0, 2) }

	// Part 2:
	//      1
	//    2 3 4
	//  5 6 7 8 9
	//    A B C
	//      D
	r2, c2 := 2, 0

	// Use the other coordinate's offset from the center (2) to determine this one's range.
	clamp2 := func(a, b int) int { return lib.Clamp(a, lib.Abs(2-b), 4-lib.Abs(2-b)) }

	buttons2 := [5]string{
		"  1  ",
		" 234 ",
		"56789",
		" ABC ",
		"  D  ",
	}

	var code, code2 string
	for _, ln := range lib.InputLines("2016/2") {
		for _, ch := range ln {
			switch ch {
			case 'U':
				r = clamp(r - 1)
				r2 = clamp2(r2-1, c2)
			case 'D':
				r = clamp(r + 1)
				r2 = clamp2(r2+1, c2)
			case 'L':
				c = clamp(c - 1)
				c2 = clamp2(c2-1, r2)
			case 'R':
				c = clamp(c + 1)
				c2 = clamp2(c2+1, r2)
			default:
				panic(fmt.Sprintf("Bad char %q", ch))
			}
		}

		code += fmt.Sprint(r*3 + c + 1)

		// Row contribution: [1, 3, 5, 3, 1]
		// Column contribution: c2 - lib.Abs(2-r2) + 1
		// I wrote this awful expression at first, but there must be a cleaner way:
		//  idx := lib.If(r2 > 0, 1, 0) + lib.If(r2 > 1, 3, 0) + lib.If(r2 > 2, 5, 0) +
		//  	lib.If(r2 > 3, 3, 0) + c2 - lib.Abs(2-r2) + 1
		//  code2 += fmt.Sprintf("%X", idx)
		code2 += buttons2[r2][c2 : c2+1]
	}
	fmt.Println(code)
	fmt.Println(code2)
}
