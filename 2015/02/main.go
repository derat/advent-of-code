package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	sum := 0
	for _, ln := range lib.InputLines("2015/2") {
		var l, w, h int
		lib.Extract(ln, `^(\d+)x(\d+)x(\d+)$`, &l, &w, &h)
		lw, wh, hl := l*w, w*h, h*l
		min := lib.Min(lw, wh, hl)
		sum += 2*lw + 2*wh + 2*hl + min
	}
	fmt.Println(sum)
}
