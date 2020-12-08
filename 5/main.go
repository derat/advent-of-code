package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	max := -1
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s := sc.Text()
		if len(s) != 10 {
			panic(fmt.Sprintf("bad line %q", s))
		}
		row := find(s[:7])
		col := find(s[7:])
		if id := row*8 + col; id > max {
			max = id
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	fmt.Println(max)
}

func find(s string) int {
	min, max := 0, int(math.Pow(2, float64(len(s)))-1)
	for _, ch := range s {
		half := (max-min)/2 + 1
		switch ch {
		case 'F', 'L':
			max -= half
		case 'B', 'R':
			min += half
		default:
			panic(fmt.Sprintf("bad line %q", s))
		}
	}
	if min != max {
		panic(fmt.Sprintf("didn't find row for %q: [%d, %d]\n", s, min, max))
	}
	return min
}
