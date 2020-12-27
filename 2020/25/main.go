package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	vals := lib.ReadIntsN(2)
	cpub, dpub := vals[0], vals[1]

	const (
		subj = 7
		div  = 20201227
	)

	val := 1
	var cloop int
	for val != cpub {
		cloop++
		val = (val * subj) % div
	}

	key := 1
	for i := 0; i < cloop; i++ {
		key = (key * dpub) % div
	}
	fmt.Println(key)
}
