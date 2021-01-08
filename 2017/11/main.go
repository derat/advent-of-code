package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// 2020/24 also used hexagonal grids, but "pointy topped" ones rather than "flat topped",
	// and without needing to calculate distances. I used offset coordinates there, but cube
	// or axial coordinates seem more appropriate here. See
	// https://www.redblobgames.com/grids/hexagons/#distances.
	//
	//            r dec
	//             ___
	//  q dec  ___/   \___  q inc
	//        /   \___/   \ r dec
	//        \___/q,r\___/
	//        /   \___/   \
	//  q dec \___/   \___/
	//  r inc     \___/     q inc
	//
	//            r inc
	var q, r int
	for _, dir := range strings.Split(lib.InputLines("2017/11")[0], ",") {
		switch dir {
		case "n":
			r--
		case "s":
			r++
		case "nw":
			q--
		case "ne":
			q++
			r--
		case "sw":
			q--
			r++
		case "se":
			q++
		default:
			lib.Panicf("Invalid direction %q", dir)
		}
	}

	dist := (lib.Abs(q) + lib.Abs(q+r) + lib.Abs(r)) / 2
	fmt.Println(dist)
}
