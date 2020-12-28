package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var paper, ribbon int
	for _, ln := range lib.InputLines("2015/2") {
		var l, w, h int
		lib.Extract(ln, `^(\d+)x(\d+)x(\d+)$`, &l, &w, &h)

		lw, wh, hl := l*w, w*h, h*l
		min := lib.Min(lw, wh, hl)
		paper += 2*lw + 2*wh + 2*hl + min

		dist := lib.Min(2*l+2*w, 2*w+2*h, 2*h+2*l)
		vol := l * w * h
		ribbon += dist + vol
	}
	fmt.Println(paper)
	fmt.Println(ribbon)
}
