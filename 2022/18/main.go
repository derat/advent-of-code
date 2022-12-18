package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	cubes := make(map[point]struct{})
	min := point{math.MaxInt, math.MaxInt, math.MaxInt}
	max := point{math.MinInt, math.MinInt, math.MinInt}
	for _, ln := range lib.InputLines("2022/18") {
		vals := lib.ExtractInts(ln)
		lib.AssertEq(len(vals), 3)
		p := point{vals[0], vals[1], vals[2]}
		cubes[p] = struct{}{}
		min.x, max.x = lib.Min(min.x, p.x), lib.Max(max.x, p.x)
		min.y, max.y = lib.Min(min.y, p.y), lib.Max(max.y, p.y)
		min.z, max.z = lib.Min(min.z, p.z), lib.Max(max.z, p.z)
	}

	// Part 1: "What is the surface area of your scanned lava droplet?"
	var area int
	sides := []point{{-1, 0, 0}, {1, 0, 0}, {0, -1, 0}, {0, 1, 0}, {0, 0, -1}, {0, 0, 1}}
	for c := range cubes {
		for _, o := range sides {
			if !lib.MapHasKey(cubes, c.add(o)) {
				area++
			}
		}
	}
	fmt.Println(area)

	// Part 2: "What is the exterior surface area of your scanned lava droplet?"
	var surface int
	bmin, bmax := min.add(point{-1, -1, -1}), max.add(point{1, 1, 1})
	vis := make(map[point]struct{})
	var grow func(point)
	grow = func(start point) {
		if lib.MapHasKey(vis, start) {
			return // already visited
		}
		vis[start] = struct{}{}
		for _, o := range sides {
			p := start.add(o)
			if p.x < bmin.x || p.x > bmax.x || p.y < bmin.y || p.y > bmax.y || p.z < bmin.z || p.z > bmax.z {
				continue // outside the bounds
			}
			if lib.MapHasKey(cubes, p) {
				surface++ // occupied; just count the side
			} else {
				grow(p) // visit neighbor
			}
		}
	}
	grow(bmin)
	fmt.Println(surface)
}

type point struct{ x, y, z int }

func (p *point) add(o point) point { return point{p.x + o.x, p.y + o.y, p.z + o.z} }
