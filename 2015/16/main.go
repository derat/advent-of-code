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

Loop:
	for _, ln := range lib.InputLines("2015/16") {
		var id int
		var rest string
		lib.Extract(ln, `^Sue (\d+): (.+)$`, &id, &rest)
		for _, s := range strings.Split(rest, ", ") {
			var item string
			var cnt int
			lib.Extract(s, `^(\w+): (\d+)$`, &item, &cnt)
			if cnt != clues[item] {
				continue Loop
			}
		}
		fmt.Println(id)
		break
	}
}
