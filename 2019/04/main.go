package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.ExtractUints(lib.Input("2019/4"))
	min, max := input[0], input[1]

	// This seems similar to a different problem I did (I forget when)
	// that also involved finding numbers between two ranges that had
	// per-digit requirements. I remember that one being pretty hard for
	// me (I coded the logic to skip over invalid numbers). For this one,
	// I just tried the stupid approach of iterating over all numbers
	// and manually inspecting each. Shockingly, it works for both parts!
	var cnt, cnt2 int
Loop:
	for n := min; n <= max; n++ {
		b := []byte(fmt.Sprint(n))
		var double, double2 bool
		for i, ch := range b {
			if i > 0 && ch < b[i-1] {
				continue Loop // digits can't decrease going left to right
			}

			var pre, pre2, post byte
			if i > 0 {
				pre = b[i-1]
			}
			if i > 1 {
				pre2 = b[i-2]
			}
			if i < 5 {
				post = b[i+1]
			}

			// Part 1: Need at least one span repeating a digit.
			if ch == pre {
				double = true
			}
			// Part 2: Need at least one span of length 2 repeating a digit.
			if ch == pre && ch != pre2 && ch != post {
				double2 = true
			}
		}
		if double {
			cnt++
		}
		if double2 {
			cnt2++
		}
	}
	fmt.Println(cnt)
	fmt.Println(cnt2)
}
