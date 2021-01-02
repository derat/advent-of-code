package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		rows    = 6
		cols    = 50
		unlit   = ' '
		lit     = '#'
		animate = true
	)

	screen := lib.NewBytes(rows, cols, ' ')

	if animate {
		fmt.Println(lib.DumpBytes(screen))
	}
	for _, ln := range lib.InputLines("2016/8") {
		switch {
		case strings.Contains(ln, "rect"):
			var w, h int
			lib.Extract(ln, `^rect (\d+)x(\d+)$`, &w, &h)
			lib.SetBytes(screen, 0, 0, h-1, w-1, lit)
		case strings.Contains(ln, "rotate row"):
			var r, amt int
			lib.Extract(ln, `^rotate row y=(\d+) by (\d+)$`, &r, &amt)
			lib.AssertLess(amt, cols)
			// http://www.cplusplus.com/reference/algorithm/rotate/
			first, middle, last := 0, cols-amt, cols
			next := middle
			for first != next {
				screen[r][first], screen[r][next] = screen[r][next], screen[r][first]
				first, next = first+1, next+1
				if next == last {
					next = middle
				} else if first == middle {
					middle = next
				}
			}
		case strings.Contains(ln, "rotate column"):
			var c, amt int
			lib.Extract(ln, `^rotate column x=(\d+) by (\d+)$`, &c, &amt)
			lib.AssertLess(amt, rows)
			first, middle, last := 0, rows-amt, rows
			next := middle
			for first != next {
				screen[first][c], screen[next][c] = screen[next][c], screen[first][c]
				first, next = first+1, next+1
				if next == last {
					next = middle
				} else if first == middle {
					middle = next
				}
			}
		default:
			panic(fmt.Sprintf("Invalid command %q", ln))
		}

		if animate {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("\033[%dA", rows)
			fmt.Println(lib.DumpBytes(screen))
		}
	}

	if !animate {
		fmt.Println(lib.DumpBytes(screen))
	}

	var cnt int
	for _, row := range screen {
		cnt += bytes.Count(row[:], []byte{'#'})
	}
	fmt.Println(cnt)
	fmt.Println(lib.OCR(screen, ' '))
}
