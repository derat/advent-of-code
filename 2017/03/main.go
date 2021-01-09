package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	n := lib.InputInts("2017/3")[0]
	x, y := coords(n)
	fmt.Println(lib.Abs(x) + lib.Abs(y))

	// Part 2: Print first summed value larger than input.
	for i := 1; true; i++ {
		x, y := coords(i)
		if sv := sumValue(x, y); sv > n {
			fmt.Println(sv)
			break
		}
	}
}

func coords(n int) (x, y int) {
	// Numbers are placed in concentric hollow squares:
	//
	//  17 16 15 14 13
	//  18  5  4  3 12
	//  19  6  1  2 11
	//  20  7  8  9 10
	//  21 22 23 24 25
	//
	// So many fiddly off-by-one errors. :-(

	for i := 0; true; i++ { // iterate over distances from axes
		dim := i*2 + 1        // dimensions of square
		sq := lib.Pow(dim, 2) // total numbers contained within square
		if sq < n {
			continue
		}

		tr, tl, bl, br := corners(dim)

		// The round starts just above the bottom-right corner
		// and then wraps around counter-clockwise.
		//  ST:  i, -i+1
		//  TR:  i,  i
		//  TL: -i,  i
		//  BL: -i, -i
		//  BR:  i, -i
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
	panic("unreached")
}

// value is the inverse of coord.
func value(x, y int) int {
	i := lib.Max(lib.Abs(x), lib.Abs(y))
	dim := i*2 + 1
	tr, tl, bl, br := corners(dim)
	switch {
	case y == -i: // bottom
		return br - lib.Abs(i-x)
	case x == -i: // left
		return bl - lib.Abs(-i-y)
	case y == i: // top
		return tl - lib.Abs(-i-x)
	case x == i: // right
		return tr - lib.Abs(i-y)
	default:
		lib.Panicf("(%d, %d) not on border of square %d", x, y, i)
		return 0
	}
}

// Cached values computed by sumValue. Base case is 1 for the center.
var sumValues = map[uint64]int{lib.PackInts(0, 0): 1}

// sumValue returns the sum of the given location's already-filled
// neighbors (including diagonals) for part 2.
func sumValue(x, y int) int {
	k := lib.PackInts(x, y)
	if sv, ok := sumValues[k]; ok {
		return sv
	}

	v := value(x, y)
	var sum int
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue // skip center
			}
			if nv := value(x+dx, y+dy); nv < v {
				sum += sumValue(x+dx, y+dy)
			}
		}
	}
	sumValues[k] = sum
	return sum
}

// corners returns the values at the corners of the square
// of the supplied dimension.
func corners(dim int) (tr, tl, bl, br int) {
	tr = lib.Pow(dim-2, 2) + dim - 1
	tl = tr + dim - 1
	bl = tl + dim - 1
	br = bl + dim - 1
	return
}
