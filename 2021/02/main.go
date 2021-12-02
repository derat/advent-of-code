package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var horiz, depth int
	var horiz2, depth2, aim2 int
	for _, ln := range lib.InputLines("2021/2") {
		var s string
		var n int
		lib.Extract(ln, `^(forward|up|down)\s+(\d+)$`, &s, &n)
		switch s {
		case "forward":
			horiz += n
			horiz2 += n
			depth2 += aim2 * n
		case "up":
			depth -= n
			aim2 -= n
		case "down":
			depth += n
			aim2 += n
		}
	}
	fmt.Println(horiz * depth)
	fmt.Println(horiz2 * depth2)
}
