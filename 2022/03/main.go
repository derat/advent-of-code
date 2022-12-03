package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum, sum2 int
	group := make(map[rune]int)
	for i, ln := range lib.InputLines("2022/3") {
		// Part 1: sum of priorities of items that appear in both compartments
		// of each elf's rucksack.
		lib.Assertf(len(ln)%2 == 0, "Line %q has %d items", ln, len(ln))
		first := make(map[rune]int)
		for _, item := range ln[:len(ln)/2] {
			first[item]++
		}
		second := make(map[rune]int)
		for _, item := range ln[len(ln)/2:] {
			second[item]++
		}
		both := lib.Intersect(first, second)
		lib.AssertEq(len(both), 1)
		sum += priority(lib.MapSomeKey(both))

		// Part 2: sum of priorities of single items that are common across each 3-elf group.
		for item := range lib.AddSet(lib.Set(first), lib.MapKeys(second)...) {
			group[item]++
		}
		if i%3 == 2 {
			for item, cnt := range group {
				if cnt == 3 {
					sum2 += priority(item)
					break
				}
			}
			group = make(map[rune]int)
		}
	}
	fmt.Println(sum)
	fmt.Println(sum2)
}

func priority(item rune) int {
	switch {
	case item >= 'a' && item <= 'z':
		return int(item-'a') + 1
	case item >= 'A' && item <= 'Z':
		return int(item-'A') + 27
	default:
		lib.Panicf("Invalid item %q", item)
		return 0
	}
}
