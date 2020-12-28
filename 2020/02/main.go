package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var valid1, valid2 int
	for _, ln := range lib.InputLines("2020/2") {
		var min, max int
		var ch, pw string
		lib.Parse(ln, `^(\d+)-(\d+) ([a-z]): ([a-z]+)$`, &min, &max, &ch, &pw)

		if n := len(pw) - len(strings.ReplaceAll(pw, ch, "")); n >= min && n <= max {
			valid1++
		}
		if min >= 1 && min <= len(pw) && max >= 1 && max <= len(pw) &&
			((pw[min-1] == ch[0]) != (pw[max-1] == ch[0])) {
			valid2++
		}
	}
	fmt.Println(valid1)
	fmt.Println(valid2)
}
