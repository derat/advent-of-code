package main

import (
	"bufio"
	"fmt"
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

	mult := 1
	for _, sl := range []struct {
		dx, dy int
	}{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	} {
		trees := 0
		for x, y := 0, 0; y < len(lines); x, y = x+sl.dx, y+sl.dy {
			ln := lines[y]
			if ch := ln[x%len(ln)]; ch == '#' {
				trees++
			}
		}
		fmt.Printf("dx %d, dy %d: %d trees\n", sl.dx, sl.dy, trees)
		mult *= trees
	}
	fmt.Printf("mult: %d\n", mult)
}
