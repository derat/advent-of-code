package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Given in puzzle:
	const mod = 2147483647
	factors := [2]int64{16807, 48271}

	init := lib.InputInts("2017/15")
	lib.AssertEq(len(init), 2)

	// Part 1: Number of pairs with matching bottom 16 bits.
	// Just using a naive approach here.
	const rounds = 40_000_000
	a, b := int64(init[0]), int64(init[1])
	var match int
	for i := 0; i < rounds; i++ {
		a = (a * factors[0]) % mod
		b = (b * factors[1]) % mod
		if a&0xffff == b&0xffff {
			match++
		}
	}
	fmt.Println(match)

	// Part 2: Again take naive approach of iterating in parallel.
	const rounds2 = 5_000_000

	cha := make(chan int64, 1)
	go gen(int64(init[0]), factors[0], mod, 4, rounds2, cha)
	chb := make(chan int64, 1)
	go gen(int64(init[1]), factors[1], mod, 8, rounds2, chb)

	var match2 int
	for i := 0; i < rounds2; i++ {
		a := <-cha
		b := <-chb
		if a&0xffff == b&0xffff {
			match2++
		}
	}
	fmt.Println(match2)
}

// gen repeatedly multiplies v by fact and takes the resulting value modulo mod.
// Values divisible by filter are sent to ch. Returns after cnt values have been sent.
func gen(v, fact, mod, filter int64, cnt int, ch chan<- int64) {
	for n := 0; n < cnt; {
		v = (v * fact) % mod
		if v%filter == 0 {
			ch <- v
			n++
		}
	}
}
