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
	d := up
	var cnt int
	for t := 0; t < 10_000; t++ {
		k := lib.PackInts(r, c)
		if _, ok := inf[k]; ok {
			d = d.right()
			delete(inf, k)
		} else {
			d = d.left()
			inf[k] = '#'
			cnt++
		}
		r, c = r+d.dr(), c+d.dc()
	}
	fmt.Println(cnt)

	// Part 2:
	inf = makeMap(init)
	r, c = len(init)/2, len(init[0])/2
	d = up
	var cnt2 int
	for t := 0; t < 10_000_000; t++ {
		k := lib.PackInts(r, c)
		switch inf[k] {
		case '#': // infected
			d = d.right()
			inf[k] = 'F'
		case 'W': // weakened
			inf[k] = '#'
			cnt2++
		case 'F': // flagged
			d = d.reverse()
			delete(inf, k)
		default: // clean
			d = d.left()
			inf[k] = 'W'
		}

		r, c = r+d.dr(), c+d.dc()
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

func (d dir) reverse() dir {
	return dir((int(d) + 2) % 4)
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
