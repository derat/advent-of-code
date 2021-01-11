package main

import (
	"fmt"
	"math/bits"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var comps []comp
	apins := make(map[int][]int)
	bpins := make(map[int][]int)
	for _, ln := range lib.InputLines("2017/24") {
		var c comp
		lib.Extract(ln, `^(\d+)/(\d+)$`, &c.a, &c.b)
		id := len(comps)
		comps = append(comps, c)
		apins[c.a] = append(apins[c.a], id)
		bpins[c.b] = append(bpins[c.b], id)
	}

	// Returns true if the first set of components is better than the second.
	type betterFunc func(uint64, uint64) bool

	var add func(int, uint64, betterFunc) uint64
	add = func(pins int, used uint64, better betterFunc) uint64 {
		best := used
		for _, id := range apins[pins] {
			if used&(1<<id) == 0 {
				if u := add(comps[id].b, used|(1<<id), better); better(u, best) {
					best = u
				}
			}
		}
		for _, id := range bpins[pins] {
			if used&(1<<id) == 0 {
				if u := add(comps[id].a, used|(1<<id), better); better(u, best) {
					best = u
				}
			}
		}
		return best
	}

	strength := func(used uint64) int {
		var sum int
		for i := 0; used != 0; i, used = i+1, used>>1 {
			if used&0x1 != 0 {
				sum += comps[i].a + comps[i].b
			}
		}
		return sum
	}

	// Part 1: Print strength of strongest possible bridge.
	fmt.Println(strength(add(0, 0, func(a, b uint64) bool {
		return strength(a) > strength(b)
	})))

	// Part 2: Print strength of longest possible bridge, using strength to break ties.
	fmt.Println(strength(add(0, 0, func(a, b uint64) bool {
		if la, lb := bits.OnesCount64(a), bits.OnesCount64(b); la > lb {
			return true
		} else if la < lb {
			return false
		}
		return strength(a) > strength(b)
	})))
}

type comp struct {
	a, b int
}
