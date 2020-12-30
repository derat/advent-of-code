package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	given := lib.InputInts("2015/20")[0]
	for house := 1; true; house++ {
		if cnt := presents(house); cnt >= given {
			fmt.Println(house)
			break
		}
	}
}

func presents(house int) int {
	var total int
	sqrt := int(math.Floor(math.Sqrt(float64(house))))
	for elf := 1; elf <= sqrt; elf++ {
		if house%elf == 0 {
			total += 10 * (elf + house/elf)
		}
	}
	return total
}
