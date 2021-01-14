package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	deps := make(map[string]map[string]struct{})
	rdeps := make(map[string][]string)
	for _, ln := range lib.InputLines("2018/7") {
		var pre, post string
		lib.Extract(ln, `^Step (\w+) must be finished before step (\w+) can begin\.$`, &pre, &post)

		m := deps[post]
		if m == nil {
			m = make(map[string]struct{})
		}
		m[pre] = struct{}{}
		deps[post] = m

		rdeps[pre] = append(rdeps[pre], post)
	}

	// Find initially-ready steps.
	ready := lib.NewHeap(func(a, b interface{}) bool { return a.(string) < b.(string) })
	for k := range rdeps {
		if _, ok := deps[k]; !ok {
			ready.Insert(k)
		}
	}

	var order string
	for ready.Len() > 0 {
		step := ready.Pop().(string)
		order += step
		for _, s := range rdeps[step] {
			delete(deps[s], step)
			if len(deps[s]) == 0 {
				ready.Insert(s)
			}
		}
	}
	fmt.Println(order)
}
