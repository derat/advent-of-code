package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/derat/advent-of-code/lib"
)

const animate = false

func main() {
	const (
		rows  = 6
		cols  = 50
		unlit = ' '
		lit   = '#'
	)

	screen := lib.NewByteGrid(rows, cols, ' ')

	if animate {
		fmt.Println(screen.Dump())
	}
	for _, ln := range lib.InputLines("2016/8") {
		switch {
		case strings.Contains(ln, "rect"):
			var w, h int
			lib.Extract(ln, `^rect (\d+)x(\d+)$`, &w, &h)
			screen.SetRect(0, 0, h-1, w-1, lit)
		case strings.Contains(ln, "rotate row"):
			var r, amt int
			lib.Extract(ln, `^rotate row y=(\d+) by (\d+)$`, &r, &amt)
			lib.RotateBy(cols, amt, func(i, j int) {
				screen[r][i], screen[r][j] = screen[r][j], screen[r][i]
			})
		case strings.Contains(ln, "rotate column"):
			var c, amt int
			lib.Extract(ln, `^rotate column x=(\d+) by (\d+)$`, &c, &amt)
			lib.RotateBy(rows, amt, func(i, j int) {
				screen[i][c], screen[j][c] = screen[j][c], screen[i][c]
			})
		default:
			panic(fmt.Sprintf("Invalid command %q", ln))
		}

		if animate {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("\033[%dA", rows)
			fmt.Println(screen.Dump())
		}
	}

	var cnt int
	for _, row := range screen {
		cnt += bytes.Count(row[:], []byte{'#'})
	}
	fmt.Println(cnt)
	fmt.Println(lib.OCR(screen, ' '))
}
