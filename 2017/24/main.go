package main

import (
	"fmt"

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

	var add func(int, int, uint64) int
	add = func(need, strength int, used uint64) int {
		max := strength
		for _, id := range apins[need] {
			if used&(1<<id) == 0 {
				max = lib.Max(max, add(comps[id].b, strength+comps[id].strength(), used|(1<<id)))
			}
		}
		for _, id := range bpins[need] {
			if used&(1<<id) == 0 {
				max = lib.Max(max, add(comps[id].a, strength+comps[id].strength(), used|(1<<id)))
			}
		}
		return max
	}

	fmt.Println(add(0, 0, 0))
}

type comp struct {
	a, b int
}

func (c comp) strength() int {
	return c.a + c.b
}
