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

	var lowest, start, cnt int
	for _, bl := range blocks {
		if bl.min <= lowest && bl.max >= lowest {
			lowest = bl.max + 1
		}
		if bl.min > start {
			cnt += bl.min - start
		}
		start = lib.Max(bl.max+1, start)
	}
	cnt += lib.Max(4294967295-start, 0)
	fmt.Println(lowest)
	fmt.Println(cnt)
}
