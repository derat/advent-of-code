package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/derat/advent-of-code/lib"
)

const (
	animate = false
	delay   = 10 * time.Millisecond
)

func main() {
	input := lib.InputInt64s("2019/13")

	// Part 1: Print count of block tiles (2) on screen when game exits.
	// Also grab the screen dimensions for part 2.
	tiles := make(map[uint64]int64)
	var rmax, cmax int
	run(input, false, func(r, c int, v int64) {
		tiles[lib.PackInts(r, c)] = v
		rmax, cmax = lib.Max(rmax, r), lib.Max(cmax, c)
	}, nil)
	var cnt int
	for _, v := range tiles {
		if v == 2 {
			cnt++
		}
	}
	fmt.Println(cnt)

	// Part 2: Beat the game by breaking all tiles and print final score.
	nrows, ncols := rmax+1, cmax+1
	screen := lib.NewByteGrid(nrows+1, ncols, ' ') // extra row for score
	clear := func() { fmt.Printf("\033[%dA", len(screen)) }

	var drawn bool         // have drawn at least one frame
	var ball, pball [2]int // cur and prev position of ball
	var pad [2]int         // cur position of paddle
	var score int64        // last score

	run(input, true, func(r, c int, v int64) {
		switch {
		case r == 0 && c == -1:
			score = v
			str := strconv.FormatInt(score, 10)
			for i := range screen[nrows] {
				if i < len(str) {
					screen[nrows][i] = str[i]
				} else {
					screen[nrows][i] = ' '
				}
			}
		default:
			screen[r][c] = tile(v)
			switch v {
			case 3:
				pad = [2]int{r, c}
			case 4:
				pball = ball
				ball = [2]int{r, c}
			}
		}
	}, func() int64 {
		if animate {
			if drawn {
				clear()
			} else {
				drawn = true
			}
			fmt.Println(screen.Dump())
			time.Sleep(delay)
		}

		if pball == [2]int{0, 0} {
			return 0 // don't know where ball is going yet
		}

		// This code doesn't even handle predicting collisions
		// with walls, but it still manages to win.
		dc := ball[1] - pball[1]   // ball horizontal movement
		nr := pad[0] - ball[0] - 1 // num rows to get to paddle
		tc := ball[1] + dc*nr      // col when ball reaches paddle
		switch {
		case tc < pad[1]:
			return -1 // left
		case tc > pad[1]:
			return 1 // right
		default:
			return 0
		}
	})

	if animate {
		clear()
		fmt.Println(screen.Dump())
	} else {
		fmt.Println(score)
	}
}

type drawFunc func(r, c int, v int64)
type inputFunc func() int64

func run(prog []int64, play bool, draw drawFunc, input inputFunc) {
	vm := lib.NewIntcode(prog)
	if play {
		vm.Mem[0] = 2
	}

	var nout int // num out calls
	var r, c int // current draw row/col

	vm.InFunc = input
	vm.OutFunc = func(v int64) {
		switch nout % 3 {
		case 0:
			c = int(v)
		case 1:
			r = int(v)
		case 2:
			draw(r, c, v)
		}
		nout++
	}
	lib.Assert(vm.Run())
}

func tile(v int64) byte {
	switch v {
	case 0:
		return ' '
	case 1:
		return 'X'
	case 2:
		return '#'
	case 3:
		return '='
	case 4:
		return 'o'
	}
	panic("Invalid tile")
}
