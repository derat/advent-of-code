package main

import (
	"bufio"
	"os"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sum := 0
	var grp []string
	for sc.Scan() {
		s := sc.Text()
		if s != "" {
			grp = append(grp, s)
		} else {
			if len(grp) > 0 {
				sum += count(grp)
			}
			grp = nil
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	if len(grp) > 0 {
		sum += count(grp)
	}
	println(sum)
}

func count(grp []string) int {
	seen := make(map[rune]struct{})
	for _, s := range grp {
		for _, ch := range s {
			seen[ch] = struct{}{}
		}
	}
	return len(seen)
}
