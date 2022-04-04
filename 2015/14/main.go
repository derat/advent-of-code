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

	// Part 1
	const total = 2503 // given in description
	dists := make(map[string]int, len(players))
	for name, rd := range players {
		cycle := rd.fly + rd.rest
		full := total / cycle
		flyTime := full*rd.fly + lib.Min(rd.fly, total%cycle)
		dists[name] = flyTime * rd.vel
	}
	fmt.Println(lib.Max(lib.MapVals(dists)...))

	// Part 2
	// Sigh, I misread this multiple times. First I thought that the extra
	// points should be added to the dists and used for determining the leader
	// at the end of each second. Then I thought that the dists and points
	// should be summed at the end instead of only looking at the points.
	dists = make(map[string]int, len(players))
	points := make(map[string]int, len(players))
	for t := 0; t < total; t++ {
		// Thoughts on possible optimizations instead of looping through each second:
		// What we really want to do (conceptually) is graph each reindeer's total
		// distance over time and find the times where there's an intersection in the
		// top line. At that point, the intersecting reindeer (two or more) each
		// receive a point, and going forward, the new leader gets 1 point/sec.
		// We could also use a heap to track the current leader.
		for name, rd := range players {
			if flying := (t % (rd.fly + rd.rest)) < rd.fly; flying {
				dists[name] += rd.vel
			}
		}
		// Award one point to current leader(s).
		max := lib.Max(lib.MapVals(dists)...)
		for _, name := range lib.MapKeysWithVal(dists, max) {
			points[name]++
		}
	}
	fmt.Println(lib.Max(lib.MapVals(points)...))
}
