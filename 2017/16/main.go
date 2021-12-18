package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var moves []move
	for _, s := range strings.Split(lib.InputLines("2017/16")[0], ",") {
		var m move
		switch {
		case lib.TryExtract(s, `^(s)(\d+)$`, &m.op, &m.v1):
		case lib.TryExtract(s, `^(x)(\d+)/(\d+)$`, &m.op, &m.v1, &m.v2):
		case lib.TryExtract(s, `^(p)([a-p])/([a-p]+)$`, &m.op, &m.p1, &m.p2):
		default:
			lib.Panicf("Invalid move %q", s)
		}
		moves = append(moves, m)
	}

	const nprogs = 16
	progs := make([]byte, nprogs)
	lib.SetAscUint(progs, 'a')
	run(progs, moves)
	fmt.Println(string(progs))

	// Part 2: Perform the dance a billion times.
	// I initially figured out how many cycles it takes for each program to end up
	// back in its initial position, but that doesn't work due to the 'p' instruction
	// (which swaps based on program name rather than index).

	// First, find out how many cycles it takes to loop back to the initial state.
	orig := make([]byte, nprogs)
	lib.SetAscUint(orig, 'a')
	lib.SetAscUint(progs, 'a')
	var loop int
	for loop = 1; true; loop++ {
		run(progs, moves)
		if reflect.DeepEqual(progs, orig) {
			break
		}
	}

	// Now perform the remaining cycles.
	for rem := 1_000_000_000 % loop; rem > 0; rem-- {
		run(progs, moves)
	}
	fmt.Println(string(progs))
}

type move struct {
	op     string
	v1, v2 int
	p1, p2 string
}

func run(progs []byte, moves []move) {
	for _, m := range moves {
		switch m.op {
		case "s":
			lib.RotateSlice(progs, m.v1)
		case "x":
			progs[m.v1], progs[m.v2] = progs[m.v2], progs[m.v1]
		case "p":
			for i, ch := range progs {
				if ch == m.p1[0] {
					progs[i] = m.p2[0]
				} else if ch == m.p2[0] {
					progs[i] = m.p1[0]
				}
			}
		default:
			lib.Panicf("Invalid move %q", m.op)
		}
	}
}
