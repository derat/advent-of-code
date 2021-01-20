package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var bots []bot
	for _, ln := range lib.InputLines("2018/23") {
		var b bot
		lib.Extract(ln, `^pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)$`, &b.x, &b.y, &b.z, &b.r)
		bots = append(bots, b)
	}

	// Part 1: Find bot with largest signal radius and count bots in range of its signals.
	maxi, maxr := -1, 0
	for i, b := range bots {
		if b.r > maxr {
			maxi, maxr = i, b.r
		}
	}
	maxb := bots[maxi]
	var cnt int
	for _, b := range bots {
		if lib.Abs(b.x-maxb.x)+lib.Abs(b.y-maxb.y)+lib.Abs(b.z-maxb.z) <= maxr {
			cnt++
		}
	}
	fmt.Println(cnt)
}

type bot struct {
	x, y, z, r int
}
