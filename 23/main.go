package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	ln, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	var cups []int
	for _, ch := range strings.TrimSpace(ln) {
		cups = append(cups, int(ch-'0'))
	}

	// Assertion for later code: cup values should be [1, ..., n].
	chk := append([]int{}, cups...)
	sort.Ints(chk)
	for i, v := range chk {
		if v != i+1 {
			log.Fatalf("Cups %v aren't ascending from 1", cups)
		}
	}

	// Part 1:
	ncups := len(cups)
	const nmoves = 100
	for move := 0; move < nmoves; move++ {
		// The destination is the value of the current cup minus one.
		dv := cups[0] - 1
		if dv == 0 {
			dv = ncups
		}

		// Remove the three cups after the current one and rotate so the
		// current cup is last.
		rem := cups[1:4]
		cups = append(cups[4:], cups[0])

		// If the destination value is on one of the removed cups,
		// count down and wrap around to top until we find a valid dest.
		for dv == rem[0] || dv == rem[1] || dv == rem[2] {
			if dv--; dv == 0 {
				dv = ncups
			}
		}

		// Find index of cup with destination value.
		di := -1
		for i, v := range cups {
			if v == dv {
				di = i
				break
			}
		}
		if di < 0 {
			log.Fatal("Didn't find dest cup ", dv)
		}

		// Place removed cups immediately after dest cup.
		nc := make([]int, 0, ncups)
		nc = append(nc, cups[:di+1]...)
		nc = append(nc, rem...)
		nc = append(nc, cups[di+1:]...)
		cups = nc
	}

	// Print all the cups after the one labeled 1.
	for i, v := range cups {
		if v == 1 {
			for j := 1; j < len(cups); j++ {
				v := cups[(i+j)%len(cups)]
				fmt.Print(v)
			}
			fmt.Println()
			break
		}
	}
}
