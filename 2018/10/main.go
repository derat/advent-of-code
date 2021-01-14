package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var points []point
	for _, ln := range lib.InputLines("2018/10") {
		var p point
		lib.Extract(ln, `^position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>$`,
			&p.x, &p.y, &p.dx, &p.dy)
		points = append(points, p)
	}

	lastArea := int64(math.MaxInt64)
	t := 1
	for ; true; t++ {
		for i := range points {
			p := &points[i]
			p.x += p.dx
			p.y += p.dy
		}
		// When the area increases from last time, assume that we just passed
		// the proper arrangement.
		_, _, width, height := bounds(points)
		area := int64(width * height)
		if area > lastArea {
			for i := range points {
				p := &points[i]
				p.x -= p.dx
				p.y -= p.dy
			}
			t--
			break
		}
		lastArea = area
	}

	xmin, ymin, width, height := bounds(points)
	grid := lib.NewBytes(height+2, width+2, ' ')
	for _, p := range points {
		r := 1 + p.y - ymin
		c := 1 + p.x - xmin
		grid[r][c] = '#'
	}
	fmt.Println(lib.OCR(grid, ' '))
	fmt.Println(t)
}

type point struct {
	x, y   int
	dx, dy int
}

func bounds(points []point) (xmin, ymin, width, height int) {
	xmin, xmax := math.MaxInt32, math.MinInt32
	ymin, ymax := math.MaxInt32, math.MinInt32
	for i := range points {
		p := &points[i]
		xmin, xmax = lib.Min(xmin, p.x), lib.Max(xmax, p.x)
		ymin, ymax = lib.Min(ymin, p.y), lib.Max(ymax, p.y)
	}
	return xmin, ymin, xmax - xmin + 1, ymax - ymin + 1
}
