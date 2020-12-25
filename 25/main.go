package main

import (
	"fmt"
	"log"
)

func main() {
	var cpub, dpub int
	if _, err := fmt.Scan(&cpub, &dpub); err != nil {
		log.Fatal(err)
	}

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
