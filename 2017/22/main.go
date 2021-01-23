package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	init := lib.InputLinesBytes("2017/22", '#', '.')

	// Part 1:
	inf := makeMap(init)
	r, c := len(init)/2, len(init[0])/2
	d := lib.Up
	var cnt int
	for t := 0; t < 10_000; t++ {
		k := lib.PackInts(r, c)
		if _, ok := inf[k]; ok {
			d = d.Right()
			delete(inf, k)
		} else {
			d = d.Left()
			inf[k] = '#'
			cnt++
		}
		r, c = r+d.DR(), c+d.DC()
	}
	fmt.Println(cnt)

	// Part 2:
	inf = makeMap(init)
	r, c = len(init)/2, len(init[0])/2
	d = lib.Up
	var cnt2 int
	for t := 0; t < 10_000_000; t++ {
		k := lib.PackInts(r, c)
		switch inf[k] {
		case '#': // infected
			d = d.Right()
			inf[k] = 'F'
		case 'W': // weakened
			inf[k] = '#'
			cnt2++
		case 'F': // flagged
			d = d.Reverse()
			delete(inf, k)
		default: // clean
			d = d.Left()
			inf[k] = 'W'
		}

		r, c = r+d.DR(), c+d.DC()
	}
	fmt.Println(cnt2)
}

func makeMap(init [][]byte) map[uint64]byte {
	inf := make(map[uint64]byte)
	for r, row := range init {
		for c, ch := range row {
			if ch == '#' {
				inf[lib.PackInts(r, c)] = '#'
			}
		}
	}
	return inf
}
