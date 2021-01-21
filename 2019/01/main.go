package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	reqs := make(map[int]int)
	var req func(int) int
	req = func(m int) int {
		if m <= 0 {
			return 0
		}
		if r, ok := reqs[m]; ok {
			return r
		}
		r := lib.Max(m/3-2, 0)
		r += req(r)
		reqs[m] = r
		return r
	}

	var sum, sum2 int
	for _, m := range lib.InputInts("2019/1") {
		sum += m/3 - 2
		sum2 += req(m)
	}
	fmt.Println(sum)
	fmt.Println(sum2)
}
