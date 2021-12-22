package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Read all input, building up a list of (possibly-overlapping) 'on' rects.
	var rects []*rect
	for _, ln := range lib.InputLines("2021/22") {
		var s string
		var r rect
		lib.Extract(ln, `^(on|off) x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)$`,
			&s, &r.x.min, &r.x.max, &r.y.min, &r.y.max, &r.z.min, &r.z.max)

		if s == "on" {
			rects = append(rects, &r)
		} else {
			// Whenever we get an 'off' instruction, subtract the rect from all of the
			// existing 'on' rects.
			nrects := make([]*rect, 0, len(rects))
			for _, e := range rects {
				nrects = append(nrects, e.sub(r)...)
			}
			rects = nrects
		}
	}

	// Part 1: "Execute the reboot steps. Afterward, considering only cubes in the region
	// x=-50..50,y=-50..50,z=-50..50, how many cubes are on?"
	//
	// Start with the initialization procedure region and subtract all of the 'on' cubes.
	// Then subtract the 'off' cubes from the whole region to figure out how many are on.
	init := rect{rng{-50, 50}, rng{-50, 50}, rng{-50, 50}}
	off := suball(&init, rects)
	var offcnt int
	for _, o := range off {
		offcnt += o.vol()
	}
	fmt.Println(init.vol() - offcnt)

	// Part 2: "Starting again with all cubes off, execute all reboot steps. Afterward,
	// considering all cubes, how many cubes are on?"
	//
	// Sum up all the 'on' rects, subtracting all previous rects each time to avoid
	// double-counting.
	var oncnt int
	for i, r := range rects {
		for _, on := range suball(r, rects[:i]) {
			oncnt += on.vol()
		}
	}
	fmt.Println(oncnt)
}

// suball subtracts others from start and returns the remaining rects.
func suball(start *rect, others []*rect) []*rect {
	rem := []*rect{start}
	for _, o := range others {
		nrem := make([]*rect, 0, len(rem))
		for _, r := range rem {
			nrem = append(nrem, r.sub(*o)...)
		}
		rem = nrem
	}
	return rem
}

type rng struct{ min, max int }

func (r *rng) size() int { return lib.If(r.max < r.min, 0, r.max-r.min+1) }

func (r *rng) overlaps(o rng) bool {
	return r.size() > 0 && o.size() > 0 && r.min <= o.max && r.max >= o.min
}

// min is the portion of r that precedes o.
// max is the portion of r that intersects o.
// max is the portion of r that comes after o.
func (r *rng) segs(o rng) (min, mid, max *rng) {
	if r.min < o.min && r.max >= o.min-1 {
		min = &rng{r.min, o.min - 1}
	}
	if r.min <= o.max && r.max >= o.min {
		mid = &rng{lib.Max(o.min, r.min), lib.Min(o.max, r.max)}
	}
	if r.max > o.max && r.min <= o.max+1 {
		max = &rng{o.max + 1, r.max}
	}
	return min, mid, max
}

type rect struct{ x, y, z rng }

func (r *rect) vol() int { return r.x.size() * r.y.size() * r.z.size() }

func (r *rect) overlaps(o rect) bool {
	return r.x.overlaps(o.x) && r.y.overlaps(o.y) && r.z.overlaps(o.z)
}

// sub subtracts o from r and returns the (zero to six) remaining rects..
func (r *rect) sub(o rect) []*rect {
	if !r.overlaps(o) {
		return []*rect{r}
	}

	var rs []*rect
	ymin, ymid, ymax := r.y.segs(o.y)

	// Single rect where r extends above o.
	if ymin != nil {
		rs = append(rs, &rect{r.x, *ymin, r.z})
	}

	// Up to four rects where r and o have overlapping y values.
	if ymid != nil {
		xmin, xmid, xmax := r.x.segs(o.x)
		if xmin != nil {
			rs = append(rs, &rect{*xmin, *ymid, r.z})
		}
		if xmid != nil {
			zmin, _, zmax := r.z.segs(o.z)
			if zmin != nil {
				rs = append(rs, &rect{*xmid, *ymid, *zmin})
			}
			if zmax != nil {
				rs = append(rs, &rect{*xmid, *ymid, *zmax})
			}
		}
		if xmax != nil {
			rs = append(rs, &rect{*xmax, *ymid, r.z})
		}
	}

	// Single rect where r extends below o.
	if ymax != nil {
		rs = append(rs, &rect{r.x, *ymax, r.z})
	}

	return rs
}
