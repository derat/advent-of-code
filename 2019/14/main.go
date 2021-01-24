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

	// Simplify recipes as much as possible: If a recipe calls for a quantity
	// of something that can be exactly produced, replace the ingredient.
	for {
		var nsimp int
		for n, r := range recipes {
			for _, r0 := range recipes {
				for in, iq := range r0.ins {
					if in == n && iq%r.q == 0 {
						for in0, iq0 := range r.ins {
							r0.ins[in0] += iq * iq0 / r.q
						}
						delete(r0.ins, n)
						nsimp++
					}
				}
			}
		}
		if nsimp == 0 {
			break
		}
	}

	// Delete recipes that are no longer needed. This drops my input from 54 to 48.
	for n := range recipes {
		if n == "FUEL" {
			continue
		}
		var used bool
		for _, r0 := range recipes {
			if _, ok := r0.ins[n]; ok {
				used = true
				break
			}
		}
		if !used {
			delete(recipes, n)
		}
	}

	// Part 1: Print amount of ORE needed to build 1 FUEL.
	fmt.Println(build("FUEL", 1, recipes, make(map[string]int)))

	// Part 2: Print amount of FUEL that can be produced with 1 trillion ORE. I spent way too long
	// trying to find a loop, in the form of the amount of fuel that has been produced once we end
	// up with zero leftovers. That worked for all of the examples but (presumably intentionally)
	// doesn't seem to terminate in a reasonable amonut of time with my input.

	// This approach of producing a fractional cost per unit of fuel and then just dividing and
	// truncating oddly produces the correct answer for my input. I would expect there to be issues
	// as you approach the maximum amount of ore.
	const maxOre = 1_000_000_000_000
	fmt.Println(int(maxOre / cost("FUEL", 1, recipes)))

	// After looking on Reddit, I saw people discussing this approach, which I should've figured
	// out: just binary search over different quanties of fuel until you find the max that consumes
	// a trillion ore or fewer. I think that I didn't even consider this since previous problems
	// usually had an iterative process where it'd be impossible to e.g. find the midpoint. Here,
	// finding the cost of a billion fuel is as easy as finding the cost of one.

	//n := sort.Search(maxOre, func(n int) bool {
	//	return build("FUEL", maxOre-n, recipes, make(map[string]int)) <= maxOre
	//})
	//fmt.Println(maxOre - n)
}

type recipe struct {
	q   int            // quantity of output
	ins map[string]int // inputs (name to quantity)
}

func build(n string, q int, recipes map[string]recipe, inv map[string]int) int {
	r, ok := recipes[n]
	lib.Assertf(ok, "No recipe for %q", n)
	nr := (q + r.q - 1) / r.q // reactions needed to produce q

	var ore int
	for in, iq := range r.ins {
		need := iq * nr
		if in == "ORE" { // building from raw ore
			ore += need
		} else if have := inv[in]; have >= need { // already have what we need
			if have == need {
				delete(inv, in)
			} else {
				inv[in] -= need
			}
		} else { // don't have enough; take what we can and build the rest
			need -= have
			delete(inv, in)
			ore += build(in, need, recipes, inv)
		}
	}
	if left := nr*r.q - q; left > 0 {
		inv[n] += left // save whatever's left over
	}
	return ore
}

func cost(n string, q float64, recipes map[string]recipe) float64 {
	r := recipes[n]
	fr := q / float64(r.q)
	var ore float64
	for in, iq := range r.ins {
		if in == "ORE" {
			ore += float64(iq) * fr
		} else {
			ore += cost(in, float64(iq)*fr, recipes)
		}
	}
	return ore
}
