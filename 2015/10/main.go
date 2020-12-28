package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var seq []uint8
	for _, v := range lib.ExtractDigits(lib.InputLines("2015/10")[0]) {
		seq = append(seq, uint8(v))
	}

	const (
		turns  = 40
		turns2 = 50
	)
	for turn := 1; turn <= turns2; turn++ {
		seq = say(seq)
		if turn == turns {
			fmt.Println(len(seq))
		} else if turn == turns2 {
			fmt.Println(len(seq))
			break
		}
	}
}

func say(seq []uint8) []uint8 {
	var newSeq []uint8
	var cur, cnt uint8
	for i, v := range seq {
		if i == 0 {
			cur = v
		}
		if v == cur {
			cnt++
		} else {
			newSeq = append(newSeq, cnt, cur)
			cur = v
			cnt = 1
		}
	}
	return append(newSeq, cnt, cur)
}
