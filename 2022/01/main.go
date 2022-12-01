package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sums []int
	for _, pg := range lib.InputParagraphs("2022/1") {
		sums = append(sums, lib.Sum(lib.ExtractInts(strings.Join(pg, " "))...))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sums)))

	// Part 1: Print calories of top elf.
	fmt.Println(sums[0])

	// Part 2: Print summed calories of top three elves.
	fmt.Println(sums[0] + sums[1] + sums[2])
}
