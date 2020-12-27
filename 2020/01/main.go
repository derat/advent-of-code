package main

import (
	"bufio"
	"fmt"
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
			fmt.Printf("%d + %d = %d, %d * %d = %d\n", n, rem, total, n, rem, n*rem)
		}
		seen[n] = struct{}{}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}

	for n := range seen {
		r1 := total - n
		for m := range seen {
			r2 := r1 - m
			if _, ok := seen[r2]; ok {
				fmt.Printf("%d + %d + %d = %d, %d * %d * %d = %d\n", n, m, r2, total, n, m, r2, n*m*r2)
			}
		}
	}
}
