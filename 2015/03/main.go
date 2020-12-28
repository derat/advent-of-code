package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var x, y int
	houses := map[string]int{"0,0": 1} // keys are "X,Y"

	// Part 2: Santa and Robo-Santa
	var sx, sy, rx, ry int // part 2: Santa and Robo-Santa
	houses2 := map[string]int{"0,0": 2}

	for i, ch := range lib.InputLinesBytes("2015/3", '^', 'v', '>', '<')[0] {
		var dx, dy int
		switch ch {
		case '^':
			dy = 1
		case 'v':
			dy = -1
		case '>':
			dx = 1
		case '<':
			dx = -1
		}

		// Part 1:
		x += dx
		y += dy
		houses[key(x, y)]++

		// Part 2:
		if i%2 == 0 {
			sx += dx
			sy += dy
			houses2[key(sx, sy)]++
		} else {
			rx += dx
			ry += dy
			houses2[key(rx, ry)]++
		}
	}

	fmt.Println(len(houses))
	fmt.Println(len(houses2))
}

func key(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}
