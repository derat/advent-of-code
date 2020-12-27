package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	const width = 36

	var mask0, mask1 uint64
	var fbits []int // positions of 'X' bits in mask
	mem := make(map[uint64]uint64)
	mem2 := make(map[uint64]uint64)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		parts := strings.Split(sc.Text(), "=")
		if len(parts) != 2 {
			log.Fatalf("bad line %q", sc.Text())
		}
		lhs := strings.TrimSpace(parts[0])
		rhs := strings.TrimSpace(parts[1])
		switch {
		case lhs == "mask":
			if len(rhs) != width {
				log.Fatalf("invalid bitmask %q", rhs)
			}
			mask0, mask1 = 0, 0
			fbits = nil
			for i, ch := range rhs {
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
		case strings.HasPrefix(lhs, "mem[") && lhs[len(lhs)-1] == ']':
			addr, err := strconv.ParseUint(lhs[4:len(lhs)-1], 10, width)
			if err != nil {
				log.Fatal("bad address: ", err)
			}
			val, err := strconv.ParseUint(rhs, 10, width)
			if err != nil {
				log.Fatal("bad value: ", err)
			}

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
		default:
			log.Fatalf("invalid lhs %q", lhs)
		}

	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
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
