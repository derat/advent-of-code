package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	scanners := make(map[int]int)
	for _, ln := range lib.InputLines("2017/13") {
		var dep, rng int
		lib.Extract(ln, `^(\d+): (\d+)$`, &dep, &rng)
		scanners[dep] = rng
	}

	// Part 1: Print severity after leaving at time 0.
	var sev int
	for dep, rng := range scanners {
		n := rng*2 - 2 // number of states
		st := dep % n  // state at arrival time
		pos := lib.If(st < rng, st, st-(st%rng+1))
		if pos == 0 {
			sev += dep * rng
		}
	}
	fmt.Println(sev)

	// Part 2: Minimum delay without getting caught (not severity 0!).
	// This seems similar to part 2 of 2020/13, except there the goal
	// is to find a departure time where all scanners (buses) are at 0,
	// whereas here we want to find a time where none of them are at 0.
	// While I still don't understand the Chinese remainder theorom,
	// I'm also not sure if it helps here due to this difference. So,
	// I guess I'll just iterate through ascending delays.
	nstates := make(map[int]int, len(scanners))
	for dep, rng := range scanners {
		nstates[dep] = rng*2 - 2
	}
Loop:
	for start := 0; true; start++ {
		for dep, n := range nstates {
			if (start+dep)%n == 0 {
				continue Loop
			}
		}
		fmt.Println(start)
		break
	}
}
