package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	init := lib.InputLinesBytes("2017/22", '#', '.')
	inf := make(map[uint64]struct{})
	for r, row := range init {
		for c, ch := range row {
			if ch == '#' {
				inf[lib.PackInts(r, c)] = struct{}{}
			}
		}
	}

	r, c := len(init)/2, len(init[0])/2
	d := up

	var cnt int
	for t := 0; t < 10_000; t++ {
		k := lib.PackInts(r, c)
		if _, ok := inf[k]; ok {
			d = d.right()
			delete(inf, k)
		} else {
			d = d.left()
			inf[k] = struct{}{}
			cnt++
		}
		r, c = r+d.dr(), c+d.dc()
	}
	fmt.Println(cnt)
}

type dir int

const (
	up dir = iota
	left
	down
	right
)

func (d dir) left() dir {
	return dir((int(d) + 1) % 4)
}

func (d dir) right() dir {
	return dir((int(d) + 3) % 4)
}

func (d dir) dr() int {
	switch d {
	case up:
		return -1
	case down:
		return 1
	case left, right:
		return 0
	default:
		panic("Invalid dir")
	}
}

func (d dir) dc() int {
	switch d {
	case up, down:
		return 0
	case left:
		return -1
	case right:
		return 1
	default:
		panic("Invalid dir")
	}
}
