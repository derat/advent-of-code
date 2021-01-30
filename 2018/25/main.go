package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var points [][]int
	for _, ln := range lib.InputLines("2018/25") {
		points = append(points, lib.ExtractInts(ln))
	}

	edges := make(map[int]map[int]struct{}, len(points)) // keys are indexes into points
	for i := range points {
		edges[i] = make(map[int]struct{})
	}
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			pi, pj := points[i], points[j]
			var dist int
			for k := range pi {
				dist += lib.Abs(pi[k] - pj[k])
			}
			if dist <= 3 {
				edges[i][j] = struct{}{}
				edges[j][i] = struct{}{}
			}
		}
	}

	rem := make(map[int]struct{})
	for i := range points {
		rem[i] = struct{}{}
	}

	var comps int
	for len(rem) > 0 {
		start := lib.MapSomeKey(rem).(int)
		inc, _ := lib.BFS([]interface{}{start}, func(si interface{}, m map[interface{}]struct{}) {
			for dst := range edges[si.(int)] {
				m[dst] = struct{}{}
			}
		}, nil)

		delete(rem, start)
		for i := range inc {
			delete(rem, i.(int))
		}

		comps++
	}
	fmt.Println(comps)
}
