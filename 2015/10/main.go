package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	seq := lib.ExtractDigits(lib.InputLines("2015/10")[0])
	const turns = 40
	for turn := 0; turn < turns; turn++ {
		var newSeq []int
		var cur, cnt int
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
		newSeq = append(newSeq, cnt, cur)
		seq = newSeq
	}
	fmt.Println(len(seq))
}
