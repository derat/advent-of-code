package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var x, y int
	dx, dy := 0, 1
	for _, s := range strings.Split(lib.InputLines("2016/1")[0], ", ") {
		var dir string
		var dist int
		lib.Extract(s, `^([LR])(\d+)$`, &dir, &dist)
		mult := lib.If(dir == "L", 1, -1)
		if dx != 0 {
			dy = mult * dx
			dx = 0
		} else {
			dx = mult * -dy
			dy = 0
		}
		x, y = x+dist*dx, y+dist*dy
	}
	fmt.Println(lib.Abs(x) + lib.Abs(y))
}
