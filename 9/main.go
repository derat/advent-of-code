package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	const psize = 25 // preamble size
	ring := make([]int, 0, psize)
	idx := 0 // latest element in ring
	lookup := make(map[int]int8, psize)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			panic(err)
		}

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
				println(n)
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
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
}
