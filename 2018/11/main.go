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

	max := math.MinInt32
	xmax, ymax := -1, -1
	for r := 2; r < len(sums); r++ {
		for c := 2; c < len(sums[r]); c++ {
			sum := sums[r][c]
			if r > 2 {
				sum -= sums[r-3][c]
			}
			if c > 2 {
				sum -= sums[r][c-3]
			}
			if r > 2 && c > 2 {
				sum += sums[r-3][c-3]
			}
			if sum > max {
				max = sum
				xmax, ymax = c-1, r-1
			}
		}
	}
	fmt.Printf("%d,%d\n", xmax, ymax)
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
