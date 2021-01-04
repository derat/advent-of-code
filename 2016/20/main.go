package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	type span struct {
		min, max int
	}
	var blocks []span
	for _, ln := range lib.InputLines("2016/20") {
		vals := lib.ExtractInts(ln)
		lib.AssertEq(len(vals), 2)
		blocks = append(blocks, span{vals[0], vals[1]})
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].min < blocks[j].min
	})

	lowest := 0
	for _, bl := range blocks {
		if bl.min <= lowest && bl.max >= lowest {
			lowest = bl.max + 1
		}
	}
	fmt.Println(lowest)
}
