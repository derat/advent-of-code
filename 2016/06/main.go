package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var counts [][256]int
	for _, ln := range lib.InputLines("2016/6") {
		if counts == nil {
			counts = make([][256]int, len(ln))
		}
		for i, c := range ln {
			counts[i][byte(c)]++
		}
	}
	var code, code2 string
	for i := range counts {
		max, min := 0, math.MaxInt32
		var best, worst byte
		for ch, cnt := range counts[i] {
			if cnt > max {
				max = cnt
				best = byte(ch)
			}
			if cnt > 0 && cnt < min {
				min = cnt
				worst = byte(ch)
			}
		}
		code += fmt.Sprintf("%c", best)
		code2 += fmt.Sprintf("%c", worst)
	}
	fmt.Println(code)
	fmt.Println(code2)
}
