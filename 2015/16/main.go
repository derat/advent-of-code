package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Attributes given in puzzle.
	clues := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	// Additional stipulations given in part 2.
	more := lib.AddSet(nil, "cats", "trees")
	less := lib.AddSet(nil, "pomeranians", "pomeranians")

	var match, match2 int
	for _, ln := range lib.InputLines("2015/16") {
		var id int
		var rest string
		lib.Extract(ln, `^Sue (\d+): (.+)$`, &id, &rest)
		var bad, bad2 bool
		for _, s := range strings.Split(rest, ", ") {
			var item string
			var cnt int
			lib.Extract(s, `^(\w+): (\d+)$`, &item, &cnt)

			exp := clues[item]

			// Part 1 needs an exact match.
			if cnt != exp {
				bad = true
			}

			// Part 2 matches different items differently.
			if _, ok := more[item]; ok {
				if cnt <= exp {
					bad2 = true
				}
			} else if _, ok := less[item]; ok {
				if cnt >= exp {
					bad2 = true
				}
			} else {
				if cnt != exp {
					bad2 = true
				}
			}
		}

		if !bad {
			match = id
		}
		if !bad2 {
			match2 = id
		}
	}

	fmt.Println(match)
	fmt.Println(match2)
}
