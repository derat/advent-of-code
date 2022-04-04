package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2018/12")

	var init string
	lib.Extract(lines[0], `^initial state: ([#.]+)$`, &init)
	plants := make(map[int]struct{})
	for i, ch := range init {
		if ch == '#' {
			plants[i] = struct{}{}
		}
	}

	pats := make(map[byte]struct{}) // patterns that yield plants
	for _, ln := range lines[1:] {
		var from, to string
		lib.Extract(ln, `^([#.]{5}) => ([#.])$`, &from, &to)
		if to != "#" {
			continue
		}
		var pat byte
		for _, ch := range from {
			pat <<= 1
			if ch == '#' {
				pat |= 1
			}
		}
		pats[pat] = struct{}{}
	}

	// Map from state to gen and value of leftmost plant.
	str, left := state(plants)
	seen := map[string][2]int{str: [2]int{0, left}}

	for gen := 1; true; gen++ {
		plants = update(plants, pats)

		// Part 1: Print sum of pots with plants after 20 generations.
		if gen == 20 {
			fmt.Println(lib.Sum(lib.MapKeys(plants)...))
		}

		// Part 2: Print sum of pots with plants after 30 billion generations.
		str, left := state(plants)
		if info, ok := seen[str]; ok {
			// My input ends up in a 1-generation loop. If the loop included more generations,
			// we'd need to repeat it as many times as possible and then manually play out the
			// remaining generations.
			pgen, pleft := info[0], info[1]
			lib.AssertEq(1, gen-pgen)
			shift := int64(left - pleft)       // shift amount each generation
			rem := 50_000_000_000 - int64(gen) // remaining generations
			val := int64(lib.Sum(lib.MapKeys(plants)...))
			fmt.Println(val + int64(len(plants))*rem*shift)
			break
		}
		seen[str] = [2]int{gen, left}
	}
}

func update(plants map[int]struct{}, pats map[byte]struct{}) map[int]struct{} {
	newPlants := make(map[int]struct{})
	done := make(map[int]struct{})
	for p := range plants {
		// Iterate over this plant and the two pots on each side of it.
		for c := p - 2; c <= p+2; c++ {
			if _, ok := done[c]; ok {
				continue // pot has already been handled
			}
			done[c] = struct{}{}

			var pat byte
			for i := c - 2; i <= c+2; i++ {
				pat <<= 1
				if _, ok := plants[i]; ok {
					pat |= 1
				}
			}
			if _, ok := pats[pat]; ok {
				newPlants[c] = struct{}{}
			}
		}
	}

	return newPlants
}

func state(plants map[int]struct{}) (string, int) {
	keys := lib.MapKeys(plants)
	min, max := lib.Min(keys...), lib.Max(keys...)
	b := make([]byte, max-min+1)
	for i := range b {
		if _, ok := plants[min+i]; ok {
			b[i] = '#'
		} else {
			b[i] = '.'
		}
	}
	return string(b), min
}
