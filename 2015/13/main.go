package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	people := make(map[string]struct{})
	rules := make(map[string]int)
	for _, ln := range lib.InputLines("2015/13") {
		var recip, verb, neighbor string
		var amount int
		lib.Extract(ln, `^(\w+) would (gain|lose) (\d+) happiness units by sitting next to (\w+)\.$`,
			&recip, &verb, &amount, &neighbor)
		if verb == "lose" {
			amount = -amount
		}
		people[recip] = struct{}{}
		people[neighbor] = struct{}{}
		rules[recip+"|"+neighbor] = amount
	}

	var max int
	var check func(seated []string, remain []string)
	check = func(seated []string, remain []string) {
		if len(remain) == 0 {
			var total int
			for i, a := range seated {
				b := seated[(i+1)%len(seated)]
				total += rules[a+"|"+b] + rules[b+"|"+a]
			}
			max = lib.Max(max, total)
			return
		}
		for i, p := range remain {
			others := append([]string{}, remain[:i]...)
			others = append(others, remain[i+1:]...)
			check(append(seated, p), others)
		}
	}
	check(nil, lib.MapStringKeys(people))
	fmt.Println(max)

	// Part 2: Add yourself.
	max = 0
	people["Me"] = struct{}{}
	check(nil, lib.MapStringKeys(people))
	fmt.Println(max)
}
