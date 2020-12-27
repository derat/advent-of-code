package main

import (
	"bufio"
	"os"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var sum, sum2 int
	var grp []string
	for sc.Scan() {
		s := sc.Text()
		if s != "" {
			grp = append(grp, s)
		} else {
			if len(grp) > 0 {
				sum += count(grp, false)
				sum2 += count(grp, true)
			}
			grp = nil
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	if len(grp) > 0 {
		sum += count(grp, false)
		sum2 += count(grp, true)
	}
	println(sum, sum2)
}

func count(grp []string, all bool) int {
	seen := make(map[rune]int)
	for _, s := range grp {
		for _, ch := range s {
			seen[ch] += 1
		}
	}
	if all {
		sum := 0
		for _, n := range seen {
			if n == len(grp) {
				sum++
			}
		}
		return sum
	}
	return len(seen)
}
