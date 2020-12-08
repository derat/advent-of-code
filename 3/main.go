package main

import (
	"bufio"
	"os"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	if len(lines) == 0 {
		panic("no input")
	}

	trees := 0
	for x, y := 0, 0; y < len(lines); x, y = x+3, y+1 {
		ln := lines[y]
		if ch := ln[x%len(ln)]; ch == '#' {
			trees++
		}
	}
	println(trees)
}
