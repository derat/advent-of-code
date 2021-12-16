package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var cnt1, cnt2 int64
	for _, state := range lib.InputInts("2021/6") {
		cnt1 += 1 + numProd(state, 80)
		cnt2 += 1 + numProd(state, 256)
	}
	fmt.Println(cnt1)
	fmt.Println(cnt2)
}

type key struct{ state, days int }

var numProdCache = make(map[key]int64)

func numProd(state, days int) int64 {
	k := key{state, days}
	if v, ok := numProdCache[k]; ok {
		return v
	}

	var cnt int64
	for days = days - state - 1; days >= 0; days -= 7 {
		cnt += 1 + numProd(8, days)
	}
	numProdCache[k] = cnt
	return cnt
}
