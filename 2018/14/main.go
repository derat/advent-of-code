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
	var scores []uint64

	var seq uint64
	nseq := len(fmt.Sprint(input))

	add := func(score int) {
		lib.AssertLess(score, 10)

		idx := nscores / size
		pos := nscores % size
		if idx == len(scores) {
			scores = append(scores, 0)
		} else if idx > len(scores) {
			lib.Panicf("Invalid index %d with scores array of len %d", idx, len(scores))
		}
		scores[idx] += uint64(score) * pows[pos]
		nscores++

		seq = (10*seq + uint64(score)) % pows[nseq]
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

	done := false // printed part 1 answer
	for {
		score1, score2 := get(elf1), get(elf2)
		sum := score1 + score2
		if sum > 9 {
			add(1)
			if seq == uint64(input) {
				break
			}
		}
		add(sum % 10)
		if seq == uint64(input) {
			break
		}

		// Part 1: Print scores of ten recipes immediately after number of recipes in input.
		if !done && nscores >= needed {
			for i := 0; i < 10; i++ {
				fmt.Print(get(input + i))
			}
			fmt.Println()
			done = true
		}

		elf1 = (elf1 + 1 + score1) % nscores
		elf2 = (elf2 + 1 + score2) % nscores
	}

	// Part 2: Print number of recipes appearing on scoreboard to left of score sequence in input.
	lib.Assertf(done, "Found part 2 sequence after %d scores; need %d for part 1", nscores, needed)
	fmt.Println(nscores - nseq)
}
