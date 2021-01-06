package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	n := lib.InputInts("2017/3")[0]
	x, y := coords(n)
	fmt.Println(lib.Abs(x) + lib.Abs(y))
}

func coords(n int) (x, y int) {
	// Numbers are placed in concentric hollow squares:
	//
	//  17  16  15  14  13
	//  18   5   4   3  12
	//  19   6   1   2  11
	//  20   7   8   9  10
	//  21  22  23---> ...
	//
	// So many fiddly off-by-one errors. :-(

	for i := 0; true; i++ { // iterate over distances from axes
		dim := i*2 + 1        // dimensions of square
		sq := lib.Pow(dim, 2) // total numbers contained within square
		if sq < n {
			continue
		}

		// The round starts just above the bottom-right corner
		// and then wraps around counter-clockwise.
		//  ST:  i, -i+1
		//  TR:  i,  i
		//  TL: -i,  i
		//  BL: -i, -i
		//  BR:  i, -i

		prev := lib.Pow(dim-2, 2) // numbers in previous square
		tr := prev + dim - 1
		tl := tr + dim - 1
		bl := tl + dim - 1
		br := sq

		switch {
		case n <= tr:
			return i, i - (tr - n)
		case n <= tl:
			return -i + (tl - n), i
		case n <= bl:
			return -i, -i + (bl - n)
		case n <= br:
			return i - (br - n), -i
		default:
			lib.Panicf("%d beyond bottom right %v", n, br)
		}
	}
}
