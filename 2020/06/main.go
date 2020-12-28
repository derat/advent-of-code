package main

import (
	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum, sum2 int
	for _, pg := range lib.InputParagraphs("2020/6") {
		sum += count(pg, false)
		sum2 += count(pg, true)
	}
	println(sum)
	println(sum2)
}

func count(grp []string, all bool) int {
	seen := make(map[rune]int)
	for _, s := range grp {
		for _, ch := range s {
			seen[ch] += 1
		}
	}
	if all {
		sum := 0
		for _, n := range seen {
			if n == len(grp) {
				sum++
			}
		}
		return sum
	}
	return len(seen)
}
