package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Part 1: Sum of signal strengths (cycle*reg) during various cycles.
	saveCycles := lib.AddSet(nil, 20, 60, 100, 140, 180, 220)
	var strengths []int

	// Part 2: Draw to 40x6 grid.
	const (
		rows = 6
		cols = 40
	)
	crt := lib.NewByteGrid(rows, cols, '.')

	reg := 1
	var cycle int
	incCycle := func() {
		cycle++
		if lib.MapHasKey(saveCycles, cycle) {
			strengths = append(strengths, cycle*reg)
		}
		if pixel := cycle - 1; pixel%cols >= reg-1 && pixel%cols <= reg+1 {
			r := pixel / cols
			c := pixel % cols
			crt[r][c] = '#'
		}
	}

	for _, ln := range lib.InputLines("2022/10") {
		var n int
		switch {
		case lib.TryExtract(ln, `^addx (-?\d+)$`, &n):
			incCycle()
			incCycle()
			reg += n
		case ln == "noop":
			incCycle()
		default:
			lib.Panicf("Bad instruction in %q", ln)
		}
		if cycle == rows*cols {
			break
		}
	}

	fmt.Println(lib.Sum(strengths...))
	fmt.Println(lib.OCR(crt, '.'))
	//fmt.Println(crt.Dump())
}
