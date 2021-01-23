package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/13")

	// Part 1: Print count of block tiles (2) on screen when game exits.
	vm := lib.NewIntcode(input)
	vm.Start()
	var nout int
	var r, c int
	tiles := make(map[uint64]int)
	for v := range vm.Out {
		switch nout % 3 {
		case 0:
			c = int(v)
		case 1:
			r = int(v)
		case 2:
			tiles[lib.PackInts(r, c)] = int(v)
		}
		nout++
	}
	lib.Assert(vm.Wait())

	var cnt int
	for _, v := range tiles {
		if v == 2 {
			cnt++
		}
	}
	fmt.Println(cnt)
}
