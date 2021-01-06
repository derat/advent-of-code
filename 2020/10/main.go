package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	outs := lib.InputInts("2020/10")
	sort.Ints(outs)
	outs = append(outs, outs[len(outs)-1]+3)

	diffs := make(map[int]int)
	last := 0
	for _, v := range outs {
		d := v - last
		lib.Assertf(d >= 1 && d <= 3, "Bad diff %v between %v and %v", d, last, v)
		diffs[d]++
		last = v
	}
	fmt.Println(diffs[1] * diffs[3])

	paths := map[int]int64{0: 1}
	for _, v := range outs {
		paths[v] = paths[v-3] + paths[v-2] + paths[v-1]
	}
	fmt.Println(paths[outs[len(outs)-1]])
}
