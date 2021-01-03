package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var discs []disc
	for _, ln := range lib.InputLines("2016/15") {
		var d disc
		lib.Extract(ln, `^Disc #\d+ has (\d+) positions; at time=0, it is at position (\d+)\.$`, &d.npos, &d.initial)
		discs = append(discs, d)
	}

	for t := 0; true; t++ {
		var fail bool
		for i, d := range discs {
			if (d.initial+t+i+1)%d.npos != 0 {
				fail = true
				break
			}
		}
		if !fail {
			fmt.Println(t)
			break
		}
	}
}

type disc struct {
	npos, initial int
}
