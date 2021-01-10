package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var parts []particle
	for _, ln := range lib.InputLines("2017/20") {
		vals := lib.ExtractInts(ln)
		lib.AssertEq(len(vals), 9)
		var p particle
		copy(p.pos[:], vals[0:3])
		copy(p.vel[:], vals[3:6])
		copy(p.acc[:], vals[6:9])
		parts = append(parts, p)
	}

	minID := -1
	var minPart particle
	for id, p := range parts {
		if minID == -1 || p.less(minPart) {
			minID = id
			minPart = p
		}
	}
	fmt.Println(minID)
}

type particle struct {
	pos, vel, acc [3]int
}

func (p particle) less(o particle) bool {
	total := func(v [3]int) int { return lib.Sum(lib.Abs(v[0]), lib.Abs(v[1]), lib.Abs(v[2])) }
	if pa, oa := total(p.acc), total(o.acc); pa < oa {
		return true
	} else if pa > oa {
		return false
	}
	if pv, ov := total(p.vel), total(o.vel); pv < ov {
		return true
	} else if pv > ov {
		return false
	}
	return total(p.pos) < total(o.pos)
}
