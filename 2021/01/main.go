package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	depths := lib.InputInts("2021/1")

	// Part 1: Count increasing depths.
	incs := 0
	for i := 1; i < len(depths); i++ {
		if depths[i] > depths[i-1] {
			incs++
		}
	}
	fmt.Println(incs)

	// Part 2: Count increasing depths using three-measurement sliding window.
	var wnds []int
	for i := range depths[:len(depths)-2] {
		wnds = append(wnds, depths[i]+depths[i+1]+depths[i+2])
	}
	wndIncs := 0
	for i := 1; i < len(wnds); i++ {
		if wnds[i] > wnds[i-1] {
			wndIncs++
		}
	}
	fmt.Println(wndIncs)
}
