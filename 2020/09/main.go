package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const psize = 25 // preamble size
	ring := make([]int64, 0, psize)
	idx := 0 // latest element in ring
	lookup := make(map[int64]int8, psize)
	var all []int64 // so much for efficiency

	for _, n := range lib.ReadInt64s() {
		// Only look for the number once we're past the preamble.
		if len(ring) == psize {
			found := false
			for _, v := range ring {
				targ := n - v
				if _, ok := lookup[targ]; ok && targ != v {
					found = true
					break
				}
			}
			if !found {
				fmt.Println(n)

				i, j := 0, 1
				total := all[i] + all[j]
			Find:
				for j < len(all) {
					switch {
					case total == n:
						min, max := int64(math.MaxInt64), int64(math.MinInt64)
						for _, v := range all[i : j+1] {
							if v < min {
								min = v
							}
							if v > max {
								max = v
							}
						}
						fmt.Println(min+max, all[i:j+1])
						break Find
					case total > n && i < j-1:
						total -= all[i]
						i++
					default:
						j++
						total += all[j]
					}
				}
				break
			}
		}

		if len(ring) < psize {
			ring = append(ring, n)
			idx = len(ring) - 1
		} else {
			idx = (idx + 1) % psize
			old := ring[idx]
			if v := lookup[old]; v == 1 {
				delete(lookup, old)
			} else {
				lookup[old]--
			}
			ring[idx] = n
		}
		lookup[n]++
		all = append(all, n)
	}
}
