package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sev int
	for _, ln := range lib.InputLines("2017/13") {
		var dep, rng int
		lib.Extract(ln, `^(\d+): (\d+)$`, &dep, &rng)
		n := rng*2 - 2                             // number of states
		st := dep % n                              // state when reached
		pos := lib.If(st < rng, st, st-(st%rng+1)) // pos when reached
		if pos == 0 {
			sev += dep * rng
		}
	}
	fmt.Println(sev)
}
