package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	const total = 2020
	seen := make(map[int]struct{})
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			panic(err)
		}
		rem := total - n
		if _, ok := seen[rem]; ok {
			println(n * rem)
			os.Exit(0)
		}
		seen[n] = struct{}{}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
}
