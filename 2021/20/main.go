package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2021/20")
	lib.AssertEq(len(pgs), 2)

	// I initially had an assert that algo[0] was '.', which saved me
	// from a bug when running on my actual input (where it's '#').
	lib.AssertEq(len(pgs[0]), 1)
	algo := []byte(pgs[0][0])

	img := makeImage()
	for r, ln := range pgs[1] {
		for c, ch := range ln {
			switch ch {
			case '#':
				img.set(point{r, c})
			case '.':
			default:
				lib.Panicf("Invalid char %q at %d,%d", ch, r, c)
			}
		}
	}

	// Part 1: How many pixels are lit after enhancing the image twice?
	for i := 1; i <= 50; i++ {
		img = img.enhance(algo)
		if i == 2 {
			fmt.Println(len(img.m))
		}
	}

	// Part 2: How many pixels are lit after 50 enhancements?
	// This was a letdown. I was expecting to need to do something clever
	// here, but my approach from part 1 runs in about half a second, and
	// the bounds of the map in the final image are only -36,-50 to 149,149.
	fmt.Println(len(img.m))

	/*
		// I tried printing the final image in the hopes that it'd display
		// something cute, but it doesn't seem like it.
		grid := lib.NewByteGrid(img.max.r-img.min.r+1, img.max.c-img.min.c+1, ' ')
		for p := range img.m {
			grid[p.r-img.min.r][p.c-img.min.c] = '#'
		}
		fmt.Println(grid.Dump())
	*/
}

type point struct{ r, c int }

type image struct {
	m        map[point]struct{} // light pixels
	min, max point              // min and max coords in m
	fill     bool               // value for pixels not in m
}

func makeImage() *image {
	return &image{
		m:   make(map[point]struct{}),
		min: point{math.MaxInt64, math.MaxInt64},
		max: point{math.MinInt64, math.MinInt64},
	}
}

func (img *image) get(p point) bool {
	if p.r < img.min.r || p.r > img.max.r || p.c < img.min.c || p.c > img.max.c {
		return img.fill
	}
	_, ok := img.m[p]
	return ok
}

func (img *image) set(p point) {
	img.m[p] = struct{}{}
	img.min.r = lib.Min(img.min.r, p.r)
	img.min.c = lib.Min(img.min.c, p.c)
	img.max.r = lib.Max(img.max.r, p.r)
	img.max.c = lib.Max(img.max.c, p.c)
}

func (img *image) enhance(algo []byte) *image {
	n := makeImage()
	for r := img.min.r - 1; r <= img.max.r+1; r++ {
		for c := img.min.c - 1; c <= img.max.c+1; c++ {
			var v int
			for dr := r - 1; dr <= r+1; dr++ {
				for dc := c - 1; dc <= c+1; dc++ {
					v <<= 1
					if img.get(point{dr, dc}) {
						v |= 1
					}
				}
			}
			if ch := algo[v]; ch == '#' {
				n.set(point{r, c})
			}
		}
	}

	if img.fill {
		n.fill = algo[511] == '#'
	} else {
		n.fill = algo[0] == '#'
	}

	return n
}
