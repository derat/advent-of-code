package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var moons []moon
	for _, ln := range lib.InputLines("2019/12") {
		var m moon
		lib.Extract(ln, `^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`, &m.pos[0], &m.pos[1], &m.pos[2])
		moons = append(moons, m)
	}

	// Part 1: Print total energy in the system after 1000 steps.
	for step := 0; step < 1000; step++ {
		// Update velocity by considering each pair of moons.
		for i := range moons {
			mi := &moons[i]
			for j := i + 1; j < len(moons); j++ {
				mj := &moons[j]
				for k := range mi.vel {
					if mi.pos[k] < mj.pos[k] {
						mi.vel[k] += 1
						mj.vel[k] -= 1
					} else if mi.pos[k] > mj.pos[k] {
						mi.vel[k] -= 1
						mj.vel[k] += 1
					}
				}
			}
		}
		// Update positions.
		for i := range moons {
			mi := &moons[i]
			for j := range mi.vel {
				mi.pos[j] += mi.vel[j]
			}
		}
	}
	var total int
	for _, m := range moons {
		pot := lib.Abs(m.pos[0]) + lib.Abs(m.pos[1]) + lib.Abs(m.pos[2])
		kin := lib.Abs(m.vel[0]) + lib.Abs(m.vel[1]) + lib.Abs(m.vel[2])
		total += pot * kin
	}
	fmt.Println(total)
}

type moon struct {
	pos [3]int
	vel [3]int
}
