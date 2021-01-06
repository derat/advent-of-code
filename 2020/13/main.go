package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lns := lib.InputLines("2020/13")
	lib.AssertEq(len(lns), 2)
	est := lib.ExtractInt64s(lns[0])[0]
	var ids []int64
	for _, s := range strings.Split(lns[1], ",") {
		if s == "x" {
			ids = append(ids, 0)
		} else {
			ids = append(ids, lib.ExtractInt64s(s)[0])
		}
	}

Loop:
	for i := est; true; i++ {
		for _, id := range ids {
			if id > 0 && i%id == 0 {
				wait := i - est
				fmt.Println(id * wait)
				break Loop
			}
		}
	}

	// Find t such that (t+n) % ids[n] = 0 for all n where ids[n] != 0.
	// IDs seem to all be coprime to each other.

	// Finds the first x >= t such that (t+off)%mod == 0.
	// The supplied step is used as an increment.
	find := func(t, mod, off, step int64) int64 {
		for {
			if (t+off)%mod == 0 {
				return t
			}
			t += step
		}
	}

	// There's probably a simpler solution but I've spent too much time on this already. :-(
	t := ids[0]
	step := ids[0]
	for i := 1; i < len(ids); i++ {
		mod := ids[i]
		if mod == 0 {
			continue
		}
		t = find(t, mod, int64(i), step)
		// Since the IDs are coprime, we can multiply them together to find the period.
		step *= mod
	}
	fmt.Println(t)
}
