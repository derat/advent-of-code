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
	mem := make(map[uint64]uint64) // input addresses seem sparse
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
				case 'X': // no-op
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
			val |= mask1
			val &= ^mask0
			mem[addr] = val
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
}
