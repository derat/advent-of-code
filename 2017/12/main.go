package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var progs []*prog
	conns := make(map[int][]int)
	for _, ln := range lib.InputLines("2017/12") {
		var id int
		var rest string
		lib.Extract(ln, `^(\d+) <-> (.+)$`, &id, &rest)
		lib.AssertEq(id, len(progs))
		progs = append(progs, &prog{make(map[int]*prog)})
		conns[id] = lib.ExtractInts(rest)
	}
	// Maybe this won't be necessary...
	for id, others := range conns {
		p := progs[id]
		for _, o := range others {
			p.conns[o] = progs[o]
		}
	}

	// Part 1: Count programs reachable from 0.
	todo := []int{0}
	seen := map[int]struct{}{}
	for len(todo) > 0 {
		id := todo[0]
		todo = todo[1:]
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		for o := range progs[id].conns {
			if _, ok := seen[o]; !ok {
				todo = append(todo, o)
			}
		}
	}
	fmt.Println(len(seen))
}

type prog struct {
	conns map[int]*prog
}
