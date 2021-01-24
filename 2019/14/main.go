package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	recipes := make(map[string]recipe)
	for _, ln := range lib.InputLines("2019/14") {
		r := recipe{ins: make(map[string]int)}
		var from, to, name string
		lib.Extract(ln, `^(.+) => (.+)$`, &from, &to)
		for _, s := range strings.Split(from, ", ") {
			var q int
			var n string
			lib.Extract(s, `^(\d+) ([A-Z]+)$`, &q, &n)
			r.ins[n] = q
		}
		lib.Extract(to, `^(\d+) ([A-Z]+)$`, &r.q, &name)
		recipes[name] = r
	}

	// Part 1: Print amount of ORE needed to build 1 FUEL.
	fmt.Println(build("FUEL", 1, recipes, make(map[string]int)))
}

type recipe struct {
	q   int            // quantity of output
	ins map[string]int // inputs (name to quantity)
}

func build(n string, q int, recipes map[string]recipe, inv map[string]int) int {
	if n == "ORE" {
		return q
	}

	if have := inv[n]; have >= q { // already have what we need
		inv[n] -= q
		return 0
	} else if have > 0 { // have some, but not enough
		q -= have
		delete(inv, n)
	}

	r, ok := recipes[n]
	lib.Assertf(ok, "No recipe for %q", n)

	nr := (q + r.q - 1) / r.q // number of reactions needed to produce q

	var ore int
	for in, iq := range r.ins {
		ore += build(in, iq*nr, recipes, inv)
	}

	inv[n] += nr*r.q - q // save whatever's left over

	return ore
}
