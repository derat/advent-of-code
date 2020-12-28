package main

import (
	"fmt"
	"log"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const width = 36

	var mask0, mask1 uint64
	var fbits []int // positions of 'X' bits in mask
	mem := make(map[uint64]uint64)
	mem2 := make(map[uint64]uint64)
	for _, ln := range lib.InputLines("2020/14") {
		var mask string
		var addr, val uint64
		lib.Parse(ln, `^mask = (.+)|mem\[(\d+)\] = (\d+)$`, &mask, &addr, &val)

		switch {
		case mask != "":
			if len(mask) != width {
				log.Fatalf("invalid bitmask %q", mask)
			}
			mask0, mask1 = 0, 0
			fbits = nil
			for i, ch := range mask {
				if i > 0 {
					mask0 <<= 1
					mask1 <<= 1
				}
				switch ch {
				case '0':
					mask0 |= 1
				case '1':
					mask1 |= 1
				case 'X':
					fbits = append(fbits, width-i-1)
				default:
					log.Fatalf("invalid bit %q", ch)
				}
			}
		default:
			// Part 1: mask applies to value: 0 unset, 1 set, X unchanged.
			mem[addr] = (val | mask1) & ^mask0

			// Part 2: mask applies to address: 0 unchanged, 1 set, X floating.
			// Recursion would be simpler, but I wanted to write an interative solution.
			for state := 0; state < 1<<len(fbits); state++ {
				maddr := addr | mask1
				for i, b := range fbits {
					if (1<<i)&state != 0 {
						maddr |= (1 << b)
					} else {
						maddr &= ^(1 << b)
					}
				}
				mem2[maddr] = val
			}
		}

	}

	var sum uint64
	for _, val := range mem {
		sum += val
	}
	fmt.Println(sum)

	var sum2 uint64
	for _, val := range mem2 {
		sum2 += val
	}
	fmt.Println(sum2)
}
