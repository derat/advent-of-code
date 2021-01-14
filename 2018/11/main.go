package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	serial := lib.InputInts("2018/11")[0]

	const dim = 300
	sums := make([][]int, dim) // sum of this cell and all above and to left
	for r := range sums {
		sums[r] = make([]int, dim)
		for c := range sums[r] {
			sum := power(c+1, r+1, serial) // x and y are 1-indexed

			if r > 0 { // add cells above
				sum += sums[r-1][c]
			}
			if c > 0 { // add cells to left
				sum += sums[r][c-1]
			}
			if r > 0 && c > 0 { // remove double-counted cells
				sum -= sums[r-1][c-1]
			}
			sums[r][c] = sum
		}
	}

	// Returns total power of square of given size and top-left 1-indexed coordinates.
	square := func(x, y, size int) int {
		r, c := y-1, x-1
		sum := sums[r+size-1][c+size-1]
		if r > 0 { // subtract cells above
			sum -= sums[r-1][c+size-1]
		}
		if c > 0 { // subtract cells to left
			sum -= sums[r+size-1][c-1]
		}
		if r > 0 && c > 0 { // add double-removed cells
			sum += sums[r-1][c-1]
		}
		return sum
	}

	// Part 1: Print top-left X,Y of 3x3 grid with largest total power.
	max := math.MinInt32
	var xmax, ymax int // 0-indexed
	for x := 1; x <= dim-2; x++ {
		for y := 1; y <= dim-2; y++ {
			if sum := square(x, y, 3); sum > max {
				max, xmax, ymax = sum, x, y
			}
		}
	}
	fmt.Printf("%d,%d\n", xmax, ymax)

	// Part 2: Print top-left X,Y,S of grid of any size with largest total power.
	max, xmax, ymax, smax := math.MinInt32, 0, 0, 0
	for size := 1; size <= dim; size++ {
		for x := 1; x <= dim-size+1; x++ {
			for y := 1; y <= dim-size+1; y++ {
				if sum := square(x, y, size); sum > max {
					max, xmax, ymax, smax = sum, x, y, size
				}
			}
		}
	}
	fmt.Printf("%d,%d,%d\n", xmax, ymax, smax)
}

func power(x, y, serial int) int {
	rack := x + 10
	pow := rack * y
	pow += serial
	pow *= rack
	pow = (pow / 100) % 10
	pow -= 5
	return pow
}
