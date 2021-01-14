package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	deps := make(map[string][]string)
	rdeps := make(map[string][]string)
	for _, ln := range lib.InputLines("2018/7") {
		var pre, post string
		lib.Extract(ln, `^Step (\w+) must be finished before step (\w+) can begin\.$`, &pre, &post)
		deps[post] = append(deps[post], pre)
		rdeps[pre] = append(rdeps[pre], post)
	}

	// Find initially-ready steps.
	ready := lib.NewHeap(func(a, b interface{}) bool { return a.(string) < b.(string) })
	for k := range rdeps {
		if _, ok := deps[k]; !ok {
			ready.Insert(k)
		}
	}

	// Repeatedly consume ready steps in alphabetical order.
	var order string
	need := makeNeed(deps)
	for ready.Len() > 0 {
		step := ready.Pop().(string)
		order += step
		for _, s := range rdeps[step] {
			delete(need[s], step)
			if len(need[s]) == 0 {
				delete(need, s)
				ready.Insert(s)
			}
		}
	}
	fmt.Println(order)

	// Part 2: Multiple workers, and time needed for each step.
	const (
		nworkers = 5
		baseDur  = 60
	)

	// Find initially-ready steps.
	for k := range rdeps {
		if _, ok := deps[k]; !ok {
			ready.Insert(k)
		}
	}

	working := make(map[string]int) // map from step to remaining time
	need = makeNeed(deps)
	var t int
	for {
		lib.Assertf(len(working) > 0 || ready.Len() > 0, "Stuck at %v", t)

		// Decrement the remaining time for in-progress steps.
		for step := range working {
			working[step]--
			if working[step] == 0 {
				delete(working, step)
				for _, s := range rdeps[step] {
					delete(need[s], step)
					if len(need[s]) == 0 {
						delete(need, s)
						ready.Insert(s)
					}
				}
			}
		}

		// Start working on new steps.
		for ready.Len() > 0 && len(working) < nworkers {
			step := ready.Pop().(string)
			working[step] = baseDur + int(step[0]-'A'+1)
		}

		if len(need) == 0 && len(working) == 0 {
			break
		}
		t++
	}
	fmt.Println(t)
}

func makeNeed(deps map[string][]string) map[string]map[string]struct{} {
	need := make(map[string]map[string]struct{}, len(deps))
	for step, reqs := range deps {
		need[step] = make(map[string]struct{})
		for _, req := range reqs {
			need[step][req] = struct{}{}
		}
	}
	return need
}
