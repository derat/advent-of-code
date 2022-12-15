package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sensors []sensor
	for _, ln := range lib.InputLines("2022/15") {
		var s sensor
		lib.Extract(ln, `^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`,
			&s.p.x, &s.p.y, &s.b.x, &s.b.y)
		sensors = append(sensors, s)
	}

	// Part 1: Count positions that can't contain a beacon in row y=2000000.
	const ty = 2000000 // example input annoyingly uses 10 instead
	var cnt int
	for _, r := range getRanges(sensors, ty) {
		cnt += r.size()
	}
	var beacons map[int]struct{} // x coords of beacons with y=ty
	for _, s := range sensors {
		if s.b.y == ty {
			beacons = lib.AddSet(beacons, s.b.x)
		}
	}
	cnt -= len(beacons) // blegh, beacons apparently not included in count
	fmt.Println(cnt)

	// Part 2: Find position in (0, 0) to (4000000, 4000000) that contains missing beacon,
	// then compute 4000000*x + y.
	//
	// Since we're told that there's a single missing beacon, it'll need to be right next
	// to one or more intersections. Iterate over pairs of sensors and intersect the line
	// segments just outside their radii.
	const max = 4000000 // example input annoyingly uses 20 instead
	for i, s1 := range sensors {
		for j, s2 := range sensors {
			if i == j {
				continue
			}
		IntersectLoop:
			for _, p := range []*point{
				intersect(s1.left(), s1.bottom(), s2.left(), s2.top()),
				intersect(s1.left(), s1.bottom(), s2.bottom(), s2.right()),
				intersect(s1.top(), s1.right(), s2.left(), s2.top()),
				intersect(s1.top(), s1.right(), s2.bottom(), s2.right()),
			} {
				if p == nil || p.x < 0 || p.x > max || p.y < 0 || p.y > max {
					continue
				}
				// I was initially calling my getRanges function here, but it didn't work
				// for some reason; I didn't dig into it further. Testing the point against
				// each sensor's radius works fine.
				for _, ts := range sensors {
					if dist(ts.p, *p) <= ts.radius() {
						continue IntersectLoop
					}
				}
				fmt.Println(4000000*int64(p.x) + int64(p.y))
				return
			}
		}
	}
}

type point struct{ x, y int }

func dist(a, b point) int { return lib.Abs(a.x-b.x) + lib.Abs(a.y-b.y) }

// intersect returns the intersection of the line segment from a1 to a2 and the line segment
// from b1 to b2. If the segments do not intersect, nil is returned.
//
// The segments must be oriented as followed (X increases going right; Y increases going down):
//
//  a1  b2
//  *   *
//   * *
//    *
//   * *
//  *   *
//  b1  a2
func intersect(a1, a2, b1, b2 point) *point {
	lib.AssertLess(a1.x, a2.x)
	lib.AssertLess(a1.y, a2.y)
	lib.AssertEq(a2.x-a1.x, a2.y-a1.y)

	lib.AssertLess(b1.x, b2.x)
	lib.AssertLess(b2.y, b1.y)
	lib.AssertEq(b2.x-b1.x, b1.y-b2.y)

	// y = a1.y + (x - a1.x)
	// y = b1.y - (x - b1.x)

	// a1.y + (x - a1.x) = b1.y - (x - b1.x)
	// (x - a1.x) + (x - b1.x) = b1.y - a1.y
	// 2x - a1.x - b1.x = b1.y - a1.y
	// 2x = a1.x + b1.y + b1.y - a1.y
	// x = (a1.x + b1.x + b1.y - a1.y) / 2
	x := (a1.x + b1.x + b1.y - a1.y)
	if x%2 != 0 {
		return nil // only return intersections at integer coordinates
	}
	x /= 2
	y := a1.y - a1.x + x

	if x < a1.x || x > a2.x || x < b1.x || x > b2.x {
		return nil
	}
	return &point{x, y}
}

type sensor struct{ p, b point }

func (s *sensor) radius() int { return dist(s.p, s.b) }
func (s *sensor) xrange(y int) *rng {
	d := s.radius() - lib.Abs(y-s.p.y)
	if d < 0 {
		return nil
	}
	return &rng{s.p.x - d, s.p.x + d}
}

// These return points just outside the range reported by the sensor.
func (s *sensor) top() point    { return point{s.p.x, s.p.y - s.radius() - 1} }
func (s *sensor) bottom() point { return point{s.p.x, s.p.y + s.radius() + 1} }
func (s *sensor) left() point   { return point{s.p.x - s.radius() - 1, s.p.y} }
func (s *sensor) right() point  { return point{s.p.x + s.radius() + 1, s.p.y} }

type rng struct{ min, max int }

func (r *rng) size() int { return r.max - r.min + 1 }

// getRanges returns ascending and non-overlapping ranges of X coordinates that
// are visible to sensors at the given Y coordinate.
func getRanges(sensors []sensor, y int) []rng {
	var rngs []rng // non-overlapping and in ascending order
	for _, s := range sensors {
		sr := s.xrange(y)
		if sr == nil {
			continue
		}

		// Insert the range at the appropriate point based on its min.
		di := sort.Search(len(rngs), func(i int) bool { return rngs[i].min >= sr.min })
		if di >= 1 && rngs[di-1].max >= sr.min {
			rngs[di-1].max = lib.Max(rngs[di-1].max, sr.max)
			di--
		} else if di == len(rngs) {
			rngs = append(rngs, *sr)
		} else {
			rngs = append(rngs[:di+1], rngs[di:]...)
			rngs[di] = *sr
		}

		// Merge ranges that overlap with the new range.
		dst := &rngs[di]
		var merged int
		for i := di + 1; i < len(rngs); i++ {
			r := rngs[i]
			if r.min > dst.max {
				break
			}
			dst.max = lib.Max(r.max, dst.max)
			merged++
		}
		if merged > 0 {
			rngs = append(rngs[:di+1], rngs[di+merged+1:]...)
		}
	}
	return rngs
}
