package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	given := lib.InputInts("2015/20")[0]

	// Part 1: Elf n delivers 10n to each house
	// Naive solution; takes a few seconds to run.
	// Probably there's a more efficient way to find factors.
	// Maybe binary-search to find the house?
	for house := 1; true; house++ {
		var cnt int
		sqrt := int(math.Floor(math.Sqrt(float64(house))))
		for elf := 1; elf <= sqrt; elf++ {
			if house%elf == 0 {
				cnt += 10 * (elf + house/elf)
			}
		}
		if cnt >= given {
			fmt.Println(house)
			break
		}
	}

	// Part 2: Elves stop delivering presents after 50 houses and deliver 11n to each house
	// Simulate the elves! Not sure how to compute a lower bound on houses...
	presents := make([]int32, given+1) // number of presents indexed by house
Loop:
	for elf := 1; elf <= given; elf++ {
		visits := 0
		for house := elf; house <= given && visits < 50; house, visits = house+elf, visits+1 {
			presents[house] += int32(11 * elf)
			// We only know that we're completely done delivering to house n after elf n.
			if int(presents[house]) >= given && elf == house {
				fmt.Println(house)
				break Loop
			}
		}
	}

	/*
		// Naive solution that just checks each house as in part 1.
		// This takes a few seconds (compared to 0.5s for simulating elves.
		for house := 1; true; house++ {
			var cnt int
			sqrt := int(math.Floor(math.Sqrt(float64(house))))

			for elf := 1; elf <= sqrt; elf++ {
				if house%elf == 0 {
					if house/elf <= 50 {
						cnt += 11 * elf
					}
					if other := house / elf; house/other <= 50 {
						cnt += 11 * other
					}
				}
			}

			if cnt >= given {
				fmt.Println(house)
				break
			}
		}
	*/
}
