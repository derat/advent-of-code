package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const dim = 1000
	on := make([][]bool, dim)
	for r := range on {
		on[r] = make([]bool, dim)
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
				case "turn off":
					on[r][c] = false
				case "toggle":
					on[r][c] = !on[r][c]
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
}
