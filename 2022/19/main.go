package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var bps []blueprint
	for i, ln := range lib.InputLines("2022/19") {
		var id int
		var bp blueprint
		lib.Extract(ln, `Blueprint (\d+): `+
			`Each ore robot costs (\d+) ore. `+
			`Each clay robot costs (\d+) ore. `+
			`Each obsidian robot costs (\d+) ore and (\d+) clay. `+
			`Each geode robot costs (\d+) ore and (\d+) obsidian.`,
			&id, &bp.oreOre, &bp.clayOre, &bp.obsOre, &bp.obsClay, &bp.geoOre, &bp.geoObs)
		lib.AssertEq(id, i+1)
		bps = append(bps, bp)
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
	oreOre          int // ore-collecting
	clayOre         int // clay-collecting
	obsOre, obsClay int // obsidian-collecting
	geoOre, geoObs  int // geode-cracking
}

type state struct {
	ore, clay, obs                  int // resources
	robOre, robClay, robObs, robGeo int // robots
	min                             int // elapsed minutes
}

// maxGeodes returns the maximum number of geodes that can be produced by bp in min.
func maxGeodes(bp blueprint, min int) int {
	return -1 * lib.AStar(
		[]state{{robOre: 1, min: 0}},
		func(s state) bool { return s.min == min },
		func(s state, next map[state]int) {
			ns := s
			ns.ore += s.robOre
			ns.clay += s.robClay
			ns.obs += s.robObs
			ns.min++
			cost := -s.robGeo

			// Transition to states where different types of robots were built.
			if s.ore >= bp.oreOre {
				n := ns
				n.ore -= bp.oreOre
				n.robOre++
				next[n] = cost
			}
			if s.ore >= bp.clayOre {
				n := ns
				n.ore -= bp.clayOre
				n.robClay++
				next[n] = cost
			}
			if s.ore >= bp.obsOre && s.clay >= bp.obsClay {
				n := ns
				n.ore -= bp.obsOre
				n.clay -= bp.obsClay
				n.robObs++
				next[n] = cost
			}
			if s.ore >= bp.geoOre && s.obs >= bp.geoObs {
				n := ns
				n.ore -= bp.geoOre
				n.obs -= bp.geoObs
				n.robGeo++
				next[n] = cost
			}
			next[ns] = cost // do nothing
		},
		func(s state) int {
			// Calculate the number of geodes that would be produced if we built ore,
			// obsidian, and geode robots as soon as possible, ignoring that types require
			// overlapping resources and that we can only build one robot per turn.

			// Ore and obsidian devoted to building robots of different types.
			oreOre := s.ore
			obsOre := s.ore
			geoOre := s.ore
			geoObs := s.obs

			var geo int
			for i := s.min; i < min; i++ {
				// Spend resources to build first.
				var madeOre, madeObs, madeGeo int
				if oreOre >= bp.oreOre {
					oreOre -= bp.oreOre
					madeOre = 1
				}
				if obsOre >= bp.obsOre {
					obsOre -= bp.obsOre
					madeObs = 1
				}
				if geoOre >= bp.geoOre && geoObs >= bp.geoObs {
					geoOre -= bp.geoOre
					geoObs -= bp.geoObs
					madeGeo = 1
				}

				// Get new resources.
				oreOre += s.robOre
				obsOre += s.robOre
				geoOre += s.robOre
				geoObs += s.robObs
				geo += s.robGeo

				// Add the newly-built robots.
				s.robOre += madeOre
				s.robObs += madeObs
				s.robGeo += madeGeo
			}
			return -geo
		})
}
