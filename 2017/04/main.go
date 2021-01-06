package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nvalid, nvalid2 int
	for _, ln := range lib.InputLines("2017/4") {
		seen := make(map[string]struct{})
		seen2 := make(map[string]struct{}) // Part 2: sorted letters

		valid, valid2 := true, true
		for _, s := range strings.Fields(ln) {
			if _, ok := seen[s]; ok {
				valid = false
			} else {
				seen[s] = struct{}{}
			}

			b := []byte(s)
			sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
			ana := string(b)
			if _, ok := seen2[ana]; ok {
				valid2 = false
			} else {
				seen2[ana] = struct{}{}
			}
		}

		if valid {
			nvalid++
		}
		if valid2 {
			nvalid2++
		}
	}
	fmt.Println(nvalid)
	fmt.Println(nvalid2)
}
