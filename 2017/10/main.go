package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const nitems = 256
	items := make([]int, nitems)
	for i := range items {
		items[i] = i
	}
	var pos, skip int
	for _, length := range lib.InputInts("2017/10") {
		lib.RotateSlice(items, -pos)
		lib.Reverse(items[:length])
		lib.RotateSlice(items, pos)
		pos = (pos + length + skip) % nitems
		skip++
	}
	fmt.Println(items[0] * items[1])
}
