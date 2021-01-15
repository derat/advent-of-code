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

	for gen := 0; gen < 20; gen++ {
		newPlants := make(map[int]struct{})
		seen := make(map[int]struct{})
		for p := range plants {
			// Iterate over this plant and the two pots on each side of it.
			for c := p - 2; c <= p+2; c++ {
				if _, ok := seen[c]; ok {
					continue // pot has already been handled
				}
				seen[c] = struct{}{}

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

		plants = newPlants
	}

	fmt.Println(lib.Sum(lib.MapIntKeys(plants)...))
}
