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
	wndIncs := 0
	prevWnd := depths[0] + depths[1] + depths[2]
	for i := 1; i < len(depths)-2; i++ {
		wnd := prevWnd - depths[i-1] + depths[i+2]
		if wnd > prevWnd {
			wndIncs++
		}
		prevWnd = wnd
	}
	fmt.Println(wndIncs)
}
