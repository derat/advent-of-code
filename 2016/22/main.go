package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nodes [][]node                                 // indexes are x, y
	for _, ln := range lib.InputLines("2016/22")[2:] { // skip dumb header
		var x, y int
		var n node
		lib.Extract(ln, `^/dev/grid/node-x(\d+)-y(\d+)\s+(\d+)T\s+(\d+)T\s+(\d+)T\s+(\d+)%$`,
			&x, &y, &n.size, &n.used, &n.avail, &n.pct)
		for x >= len(nodes) {
			nodes = append(nodes, nil)
		}
		for y >= len(nodes[x]) {
			nodes[x] = append(nodes[x], node{})
		}
		nodes[x][y] = n
	}

	var cnt int
	for x, ns := range nodes {
		for y, a := range ns {
			for x2, ns2 := range nodes {
				for y2, b := range ns2 {
					if x == x2 && y == y2 {
						continue // don't compare nodes against themselves
					}
					if a.used == 0 {
						continue // skip empty source nodes
					}
					if a.used <= b.avail {
						cnt++
					}
				}
			}
		}
	}
	fmt.Println(cnt)
}

type node struct {
	size, used, avail, pct int
}
