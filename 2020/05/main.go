package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	max := -1
	seen := make(map[int]struct{})
	for _, s := range lib.InputLinesBytes("2020/5", 'F', 'B', 'R', 'L') {
		if len(s) != 10 {
			panic(fmt.Sprintf("bad line %q", s))
		}
		row := find(s[:7])
		col := find(s[7:])
		id := row*8 + col
		if id > max {
			max = id
		}
		seen[id] = struct{}{}
	}
	fmt.Printf("max: %d\n", max)

	for i := 1; i < 8*128; i++ {
		if _, ok := seen[i]; ok {
			continue
		}
		if _, ok := seen[i-1]; !ok {
			continue
		}
		if _, ok := seen[i+1]; !ok {
			continue
		}
		fmt.Printf("yours: %d\n", i)
	}
}

func find(s []byte) int {
	min, max := 0, int(math.Pow(2, float64(len(s)))-1)
	for _, ch := range s {
		half := (max-min)/2 + 1
		switch ch {
		case 'F', 'L':
			max -= half
		case 'B', 'R':
			min += half
		default:
			panic(fmt.Sprintf("bad line %q", s))
		}
	}
	if min != max {
		panic(fmt.Sprintf("didn't find row for %q: [%d, %d]\n", s, min, max))
	}
	return min
}
