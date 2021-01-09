package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const nprogs = 16
	progs := make([]byte, nprogs)
	for i := range progs {
		progs[i] = byte('a' + i)
	}

	for _, s := range strings.Split(lib.InputLines("2017/16")[0], ",") {
		var a, b int
		var x, y string
		switch {
		case lib.ExtractMaybe(s, `^s(\d+)$`, &a):
			lib.RotateSlice(progs, a)
		case lib.ExtractMaybe(s, `^x(\d+)/(\d+)$`, &a, &b):
			progs[a], progs[b] = progs[b], progs[a]
		case lib.ExtractMaybe(s, `^p([a-p])/([a-p]+)$`, &x, &y):
			for i, ch := range progs {
				if ch == x[0] {
					progs[i] = y[0]
				} else if ch == y[0] {
					progs[i] = x[0]
				}
			}
		}
	}
	fmt.Println(string(progs))
}
