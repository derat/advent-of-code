package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2022/20")
	fmt.Println(computeSum(mix(input, 1)))

	input2 := make([]int, len(input))
	for i, v := range input {
		input2[i] = v * 811589153
	}
	fmt.Println(computeSum(mix(input2, 10)))
}

// mix mixes orig count times.
func mix(orig []int, count int) []int {
	mixed := make([]int, len(orig)) // orig indexes in current mixed order
	for i := range mixed {
		mixed[i] = i
	}

	for c := 0; c < count; c++ {
		for i, n := range orig {
			var src int
			for ; src < len(mixed); src++ {
				if mixed[src] == i {
					break
				}
			}
			lib.AssertLess(src, len(mixed))

			// This is mod len-1 to handle wrapping: when index 0 moves left,
			// it needs to end up to the left of the len-1 number.
			dst := (src + n) % (len(mixed) - 1)
			if dst < 0 {
				dst += (len(mixed) - 1)
			}

			if dst < src {
				copy(mixed[dst+1:], mixed[dst:src])
				mixed[dst] = i
			} else if dst > src {
				copy(mixed[src:], mixed[src+1:dst+1])
				mixed[dst] = i
			}
		}
	}

	out := make([]int, len(mixed))
	for i, oi := range mixed {
		out[i] = orig[oi]
	}
	return out
}

// computeSum returns the sum of the 1000th, 2000th, and 3000th numbers
// appearing after 0 in vals.
func computeSum(vals []int) int {
	for i, v := range vals {
		if v == 0 {
			var sum int
			for j := 0; j <= 3000; j++ {
				if j == 1000 || j == 2000 || j == 3000 {
					sum += vals[(i+j)%len(vals)]
				}
			}
			return sum
		}
	}
	panic("0 not found")
}
