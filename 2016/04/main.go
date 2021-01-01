package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum, id2 int
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

			// Part 2: Decrypt room names.
			dec := []byte(name)
			for i, ch := range dec {
				if ch == '-' {
					dec[i] = ' '
				} else {
					dec[i] = byte('a' + (int(ch-'a')+id)%26)
				}
			}

			// This felt ugly. The question just asks "What is the sector ID of the room where North
			// Pole objects are stored?", so it seems like you're supposed to print the decoded
			// names first and then manually scan through them for one that seems similar to the
			// description. For my input, it's "northpole object storage".
			if s := strings.ToLower(string(dec)); strings.Contains(s, "north") && strings.Contains(s, "pole") {
				id2 = id
			}
		}
	}
	fmt.Println(sum)
	fmt.Println(id2)
}
