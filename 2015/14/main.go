package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

type reindeer struct {
	vel, fly, rest int
}

func main() {
	players := make(map[string]reindeer)
	for _, ln := range lib.InputLines("2015/14") {
		var name string
		var rd reindeer
		lib.Extract(ln, `^(\w+) can fly (\d+) km/s for (\d+) seconds?, but then must rest for (\d+) seconds?\.$`,
			&name, &rd.vel, &rd.fly, &rd.rest)
		players[name] = rd
	}

	const elapsed = 2503 // given in description
	dists := make(map[string]int, len(players))
	for name, rd := range players {
		cycle := rd.fly + rd.rest
		full := elapsed / cycle
		flyTime := full*rd.fly + lib.Min(rd.fly, elapsed%cycle)
		dists[name] = flyTime * rd.vel
	}
	fmt.Println(lib.Max(lib.MapIntVals(dists)...))
}
