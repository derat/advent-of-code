package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const dim = 1000

	// Part 1
	on := make([][]bool, dim)
	for r := range on {
		on[r] = make([]bool, dim)
	}

	// Part 2
	brt := make([][]uint16, dim)
	for r := range brt {
		brt[r] = make([]uint16, dim)
	}

	for _, ln := range lib.InputLines("2015/6") {
		var op string
		var r1, c1, r2, c2 int
		lib.Extract(ln, `^(turn on|turn off|toggle) (\d+),(\d+) through (\d+),(\d+)$`, &op, &r1, &c1, &r2, &c2)
		r1, r2 = lib.Min(r1, r2), lib.Max(r1, r2)
		c1, c2 = lib.Min(c1, c2), lib.Max(c1, c2)
		for r := r1; r <= r2; r++ {
			for c := c1; c <= c2; c++ {
				switch op {
				case "turn on":
					on[r][c] = true
					brt[r][c]++
				case "turn off":
					on[r][c] = false
					if brt[r][c] > 0 {
						brt[r][c]--
					}
				case "toggle":
					on[r][c] = !on[r][c]
					brt[r][c] += 2
				}
			}
		}
	}

	var cnt int
	for _, r := range on {
		for _, v := range r {
			if v {
				cnt++
			}
		}
	}
	fmt.Println(cnt)

	var cnt2 int64
	for _, r := range brt {
		for _, v := range r {
			cnt2 += int64(v)
		}
	}
	fmt.Println(cnt2)
}
