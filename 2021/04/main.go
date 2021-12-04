package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputParagraphs("2021/4")
	lib.AssertEq(len(input[0]), 1)
	drawn := lib.ExtractInts(input[0][0])
	var boards []board
	for _, lines := range input[1:] {
		b := board{
			seen:   lib.NewBytes(len(lines), len(lines), ' '),
			lookup: make(map[int][2]int),
		}
		for _, ln := range lines {
			r := lib.ExtractInts(ln)
			for c, v := range r {
				lib.Assert(!lib.MapHasKey(b.lookup, v))
				b.lookup[v] = [2]int{len(b.rows), c}
			}
			b.rows = append(b.rows, r)

		}
		boards = append(boards, b)
	}

	var sawWin bool
Loop:
	for _, d := range drawn {
		keptBoards := make([]board, 0, len(boards))
		for _, b := range boards {
			if won := b.set(d); won {
				// Part 1: Print score (sum of unmarked times drawn number) of first winner.
				if !sawWin {
					fmt.Println(b.score(d))
					sawWin = true
				}
				// Part 2: Print score of last winner.
				if len(boards) == 1 {
					fmt.Println(b.score(d))
					break Loop
				}
			} else {
				keptBoards = append(keptBoards, b)
			}
		}
		boards = keptBoards
	}
}

type board struct {
	rows   [][]int
	seen   [][]byte
	lookup map[int][2]int // vals are [row, col]
}

func (b *board) set(v int) bool {
	coords, ok := b.lookup[v]
	if !ok {
		return false
	}
	r, c := coords[0], coords[1]
	b.seen[r][c] = 'x'
	nr, nc := len(b.rows), len(b.rows[r])
	return lib.CountBytes(b.seen, r, 0, r, nc-1, 'x') == nc ||
		lib.CountBytes(b.seen, 0, c, nr-1, c, 'x') == nr
}

func (b *board) score(drawn int) int {
	var sum int
	for r, row := range b.rows {
		for c, val := range row {
			if b.seen[r][c] == ' ' {
				sum += val
			}
		}
	}
	return sum * drawn
}
