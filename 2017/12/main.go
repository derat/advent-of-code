package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	conns := make(map[int][]int)
	for _, ln := range lib.InputLines("2017/12") {
		var id int
		var rest string
		lib.Extract(ln, `^(\d+) <-> (.+)$`, &id, &rest)
		conns[id] = lib.ExtractInts(rest)
	}

	// Part 1: Count programs reachable from 0.
	fmt.Println(len(find(conns, 0)))

	// Part 2: Count groups.
	var cnt int
	for len(conns) > 0 {
		for id := range find(conns, lib.MapSomeKey(conns).(int)) {
			delete(conns, id)
		}
		cnt++
	}
	fmt.Println(cnt)
}

// find returns the set of nodes in conns reachable from start (including start itself).
func find(conns map[int][]int, start int) map[int]struct{} {
	todo := []int{start}
	seen := map[int]struct{}{}
	for len(todo) > 0 {
		id := todo[0]
		todo = todo[1:]
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		for _, o := range conns[id] {
			if _, ok := seen[o]; !ok {
				todo = append(todo, o)
			}
		}
	}
	return seen
}
