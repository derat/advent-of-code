package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum int
	for _, ln := range lib.InputLines("2016/4") {
		var name, checksum string
		var id int
		lib.Extract(ln, `^([-a-z]+)-(\d+)\[([a-z]+)\]$`, &name, &id, &checksum)
		var counts [256]int
		for _, ch := range name {
			counts[byte(ch)]++
		}
		letters := make([]byte, 26)
		for i := range letters {
			letters[i] = 'a' + byte(i)
		}
		sort.Slice(letters, func(i, j int) bool {
			li, lj := letters[i], letters[j]
			return counts[li] > counts[lj] || (counts[li] == counts[lj] && li < lj)
		})
		if string(letters[:5]) == checksum {
			sum += id
		}
	}
	fmt.Println(sum)
}
