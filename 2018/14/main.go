package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2018/14")[0]
	needed := input + 10

	const size = 19 // number of digits that can fit in uint64
	pows := make([]uint64, size)
	pows[0] = 1
	for i := 1; i < len(pows); i++ {
		pows[i] = 10 * pows[i-1]
	}

	var nscores int
	scores := make([]uint64, (needed+size-1)/size)

	add := func(score int) {
		lib.AssertLess(score, 10)
		idx := nscores / size
		pos := nscores % size
		scores[idx] += uint64(score) * pows[pos]
		nscores++
	}

	get := func(n int) int {
		lib.AssertLess(n, nscores)
		idx := n / size
		pos := n % size
		score := scores[idx]
		if pos < len(pows)-1 {
			score %= pows[pos+1]
		}
		return int(score / pows[pos])
	}

	// Initial scores.
	add(3)
	add(7)
	elf1, elf2 := 0, 1

	// Part 1: Print scores of ten recipes immediately after number of recipes in input.
	for nscores < needed {
		score1, score2 := get(elf1), get(elf2)
		sum := score1 + score2
		if sum > 9 {
			add(1)
		}
		add(sum % 10)

		elf1 = (elf1 + 1 + score1) % nscores
		elf2 = (elf2 + 1 + score2) % nscores
	}
	for i := 0; i < 10; i++ {
		fmt.Print(get(input + i))
	}
	fmt.Println()
}
