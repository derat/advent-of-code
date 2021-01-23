package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var input []moon
	for _, ln := range lib.InputLines("2019/12") {
		var m moon
		lib.Extract(ln, `^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`, &m.pos[0], &m.pos[1], &m.pos[2])
		input = append(input, m)
	}

	// Part 1: Print total energy in the system after 1000 steps.
	moons := append([]moon(nil), input...)
	for step := 0; step < 1000; step++ {
		update(moons)
	}
	var total int
	for _, m := range moons {
		pot := lib.Abs(m.pos[0]) + lib.Abs(m.pos[1]) + lib.Abs(m.pos[2])
		kin := lib.Abs(m.vel[0]) + lib.Abs(m.vel[1]) + lib.Abs(m.vel[2])
		total += pot * kin
	}
	fmt.Println(total)

	// Part 2: Print number of steps to reach a state that exactly matches
	// an earlier state. I just looked at X first and noticed that the first
	// earlier state that it loops around to is the initial state, so I made
	// the assumption that I could do the same for the other components as
	// well. This yielded the correct answer, but I'm not sure it's guaranteed.
	moons = append([]moon(nil), input...)
	var orig [3]string // serialization of pos and vel for all moons in each dimension
	for i := range orig {
		orig[i] = state(moons, i)
	}
	var loop [3]int // step at which orig was seen again for each dimension

	for step := 0; true; step++ {
		for i := range loop {
			if loop[i] == 0 && state(moons, i) == orig[i] {
				loop[i] = step
			}
		}
		if loop[0] > 0 && loop[1] > 0 && loop[2] > 0 {
			break
		}
		update(moons)
	}
	fmt.Println(lib.LCM(loop[:]...))
}

func update(moons []moon) {
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

func state(moons []moon, dim int) string {
	var vals []string
	for i := range moons {
		vals = append(vals, strconv.Itoa(moons[i].pos[dim]))
		vals = append(vals, strconv.Itoa(moons[i].vel[dim]))
	}
	return strings.Join(vals, "|")
}

type moon struct {
	pos [3]int
	vel [3]int
}
