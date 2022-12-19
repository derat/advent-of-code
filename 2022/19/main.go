package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var bps []blueprint
	for i, ln := range lib.InputLines("2022/19") {
		var id int
		// This is dumb: I want to use uint8 in the struct so I don't need a bunch
		// of ugly casts later, but my lib.Extract function expects a single-char
		// string rather than an int when extracting to uint8 (since byte and uint8
		// are the same thing in Go).
		var oreOre, clayOre, obsOre, obsClay, geoOre, geoObs int
		lib.Extract(ln, `Blueprint (\d+): `+
			`Each ore robot costs (\d+) ore. `+
			`Each clay robot costs (\d+) ore. `+
			`Each obsidian robot costs (\d+) ore and (\d+) clay. `+
			`Each geode robot costs (\d+) ore and (\d+) obsidian.`,
			&id, &oreOre, &clayOre, &obsOre, &obsClay, &geoOre, &geoObs)
		lib.AssertEq(id, i+1)
		bps = append(bps, blueprint{
			oreOre:  uint8(oreOre),
			clayOre: uint8(clayOre),
			obsOre:  uint8(obsOre),
			obsClay: uint8(obsClay),
			geoOre:  uint8(geoOre),
			geoObs:  uint8(geoObs),
		})
	}

	// Part 1: Sum of quality levels (1-based index times largest number of geodes
	// that can be opened in 24 minutes) of all blueprints.
	var sum int
	for i, bp := range bps {
		geo := maxGeodes(bp, 24)
		sum += (i + 1) * geo
	}
	fmt.Println(sum)

	// Part 2: Product of max geodes than can be opened in 32 minutes from first
	// three blueprints.
	var prod int64 = 1
	for _, bp := range bps[:3] {
		geo := maxGeodes(bp, 32)
		prod *= int64(geo)
	}
	fmt.Println(prod)
}

// blueprint holds the costs of building different types of robots.
type blueprint struct {
	oreOre          uint8 // ore-collecting
	clayOre         uint8 // clay-collecting
	obsOre, obsClay uint8 // obsidian-collecting
	geoOre, geoObs  uint8 // geode-cracking
}

// maxGeodes returns the maximum number of geodes that can be produced by bp in min.
func maxGeodes(bp blueprint, min uint8) int {
	ests := make(map[[2]uint8]int) // keys are geode robots and starting minute
	return -1 * lib.AStar(
		[]state{{robOre: 1, min: 0}},
		func(s state) bool { return s.min == min },
		func(s state, next map[state]int) {
			ns := s
			inc(&ns.ore, s.robOre)
			inc(&ns.clay, s.robClay)
			inc(&ns.obs, s.robObs)
			ns.min++
			cost := -int(s.robGeo)

			if s.ore >= bp.oreOre {
				n := ns
				n.ore -= bp.oreOre
				inc(&n.robOre, 1)
				next[n] = cost
			}
			if s.ore >= bp.clayOre {
				n := ns
				n.ore -= bp.clayOre
				inc(&n.robClay, 1)
				next[n] = cost
			}
			if s.ore >= bp.obsOre && s.clay >= bp.obsClay {
				n := ns
				n.ore -= bp.obsOre
				n.clay -= bp.obsClay
				inc(&n.robObs, 1)
				next[n] = cost
			}
			if s.ore >= bp.geoOre && s.obs >= bp.geoObs {
				n := ns
				n.ore -= bp.geoOre
				n.obs -= bp.geoObs
				inc(&n.robGeo, 1)
				next[n] = cost
			}
			next[ns] = cost // do nothing
		},
		func(s state) int {
			key := [2]uint8{s.robGeo, s.min}
			if v, ok := ests[key]; ok {
				return v
			}
			// As a bad lower (i.e. upper) bound, use the geodes that will be produced
			// in the remaining time by the existing geode-cracking robots plus the
			// geodes that we'd get if we created a new geode-cracking robot in each
			// remaining minute.
			// TODO: Try to improve this. Maybe I could look at robot costs and existing
			// resources and use that to further limit the maximum number of geode-cracking
			// robots that can be created.
			rem := int(min - s.min)
			geo := int(s.robGeo)*rem + ((rem-1)*rem)/2
			cost := -geo
			ests[key] = cost
			return cost
		})
}

// inc increments the uint8 at v by n while checking for overflow.
func inc(v *uint8, n uint8) {
	old := *v
	*v += n
	lib.AssertGreaterEq(*v, old)
}

type state struct {
	ore, clay, obs                  uint8 // resources
	robOre, robClay, robObs, robGeo uint8 // robots
	min                             uint8 // elapsed minutes
}
