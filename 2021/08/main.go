package main

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var cnt, sum int
	for _, ln := range lib.InputLines("2021/8") {
		fields := strings.Fields(ln)
		lib.AssertEq(len(fields), 15) // 10 unique signal patterns, '|', 4 output values
		pats := fields[:10]
		outs := fields[11:]

		// Part 1: Count easy output digits:
		//  1 uses 2 segments
		//  4 uses 4 segments
		//  7 uses 3 segments
		//  8 uses 7 segments
		for _, o := range outs {
			if segs := len(o); segs == 2 || segs == 4 || segs == 3 || segs == 7 {
				cnt++
			}
		}

		// Part 2: Decode and sum all output numbers (not digits).
		//
		//    0:      1:      2:      3:      4:
		//   aaaa    ....    aaaa    aaaa    ....
		//  b    c  .    c  .    c  .    c  b    c
		//  b    c  .    c  .    c  .    c  b    c
		//   ....    ....    dddd    dddd    dddd
		//  e    f  .    f  e    .  .    f  .    f
		//  e    f  .    f  e    .  .    f  .    f
		//   gggg    ....    gggg    gggg    ....
		//
		//    5:      6:      7:      8:      9:
		//   aaaa    aaaa    aaaa    aaaa    aaaa
		//  b    .  b    .  .    c  b    c  b    c
		//  b    .  b    .  .    c  b    c  b    c
		//   dddd    dddd    ....    dddd    dddd
		//  .    f  e    f  .    f  e    f  .    f
		//  .    f  e    f  .    f  e    f  .    f
		//   gggg    gggg    ....    gggg    gggg
		//
		// Across all digits:
		//  a is used 8 times
		//  b is used 6 times
		//  c is used 8 times
		//  d is used 7 times
		//  e is used 4 times
		//  f is used 9 times
		//  g is used 7 times
		//
		// a is the signal that appears in 7 (3 segments) but not 4 (2 segments)
		//
		// So it's easy to get a, b, c, e, and f.
		// That leaves d and g, both of which appear in 7 digits.
		// Between d and g, g is the signal that doesn't appear in 4 (which we know).

		toBits := func(pat string) uint8 {
			var b uint8
			for _, r := range pat {
				b |= (1 << (r - 'a'))
			}
			return b
		}

		var four, seven uint8
		seen := make(map[rune]int) // input signal counts across 0-9
		for _, p := range pats {
			for _, r := range p {
				seen[r]++
			}
			// The 2 and 7 cases correspond to one and eight, respectively, but we don't need them.
			switch len(p) {
			case 3:
				seven = toBits(p)
			case 4:
				four = toBits(p)
			}
		}

		unfix := make(map[rune]rune, 10) // correct signals -> input signals

		a := 'a' + rune(bits.TrailingZeros8(seven & ^four))
		unfix['a'] = a

		for r, n := range seen {
			switch n {
			case 4:
				unfix['e'] = r
			case 6:
				unfix['b'] = r
			case 7:
				if toBits(string(r))&four == 0 {
					unfix['g'] = r
				} else {
					unfix['d'] = r
				}
			case 8:
				if r != a {
					unfix['c'] = r
				}
			case 9:
				unfix['f'] = r
			}

		}

		lookup := make(map[uint8]int)
		for p, n := range map[string]int{
			"abcefg":  0,
			"cf":      1,
			"acdeg":   2,
			"acdfg":   3,
			"bcdf":    4,
			"abdfg":   5,
			"abdefg":  6,
			"acf":     7,
			"abcdefg": 8,
			"abcdfg":  9,
		} {
			var ip string
			for _, r := range p {
				ip += string(unfix[r])
			}
			lookup[toBits(ip)] = n
		}

		var n int
		for _, o := range outs {
			n = 10*n + lookup[toBits(o)]
		}
		sum += n
	}
	fmt.Println(cnt)
	fmt.Println(sum)
}
