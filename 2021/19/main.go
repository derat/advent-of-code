package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

const thresh = 12 // match threshold (given in problem)

func main() {
	var scanners []*scanner
	for _, pg := range lib.InputParagraphs("2021/19") {
		lib.AssertLessEq(2, len(pg))
		var id int
		lib.Extract(pg[0], `^--- scanner (\d+) ---$`, &id)
		var pts []point
		for _, ln := range pg[1:] {
			p := lib.ExtractInts(ln)
			lib.AssertEq(len(p), 3)
			pts = append(pts, point{p[0], p[1], p[2]})
		}
		scanners = append(scanners, newScanner(id, pts))
	}

	// Part 1: How many beacons are there?
	// Start with the untransformed view from the first scanner.
	glob := make(pointMap)
	for _, p := range scanners[0].beacons {
		glob[p] = struct{}{}
	}

	offs := []point{point{}}                            // scanner offsets from scanner 0 for part 2
	rem := make(map[*scanner]struct{}, len(scanners)-1) // unaligned scanners
	for _, s := range scanners[1:] {
		rem[s] = struct{}{}
	}

	for len(rem) > 0 {
		// For each scanner that we haven't yet lined up with the global map,
		// iterate over all of its orientations.
		for s := range rem {
		ViewLoop:
			for _, v := range s.views {
				// Try offsetting each of the scanner's beacons to match each
				// beacon in the global map and check if enough of the scanner's
				// beacons match known beacons.
				//
				// This is slow and I'm sure there's a better approach. However, my
				// input still just takes 3-4 seconds, which is good enough for me.
				//
				// I'm also not sure if this is totally correct, since the problem
				// just describes comparing one scanner's beacons to another scanner's.
				// I'm comparing the scanner's beacons against the global view, which
				// contains multiple scanner's beacons. It'd probably be better to do
				// some more bookkeeping and check separately against each known scanner,
				// but that feels a bit tedious so I'm not bothering with it for now.
				//
				// Update: I tried comparing against known scanner views instead of the
				// global view, and my naive implementation was even slower than this.
				// *shrug*
				for p := range v {
					for gp := range glob {
						off := sub(gp, p)
						if overlap(v, glob, off) {
							for p := range v {
								glob[add(p, off)] = struct{}{}
							}
							offs = append(offs, off)
							delete(rem, s)
							break ViewLoop
						}
					}
				}
			}
		}
	}
	fmt.Println(len(glob))

	// Part 2: Find the maximum Manhattan distance between two scanners (not beacons --
	// I misread this at first and was confused about why I was getting an answer that
	// didn't match the example answer for something that seemed so simple).
	max := -1
	for _, a := range offs {
		for _, b := range offs {
			max = lib.Max(max, lib.Abs(a.x-b.x)+lib.Abs(a.y-b.y)+lib.Abs(a.z-b.z))
		}
	}
	fmt.Println(max)
}

type point struct{ x, y, z int }
type pointMap map[point]struct{}

func add(a, b point) point { return point{a.x + b.x, a.y + b.y, a.z + b.z} }
func sub(a, b point) point { return point{a.x - b.x, a.y - b.y, a.z - b.z} }

// overlap returns true if ma and mb have at least thresh matching points
// after offsetting ma by off.
func overlap(ma, mb pointMap, off point) bool {
	var cnt, seen int
	for pa := range ma {
		if _, ok := mb[add(pa, off)]; ok {
			// After we've reached the threshold, we don't need to keep counting.
			if cnt++; cnt >= thresh {
				return true
			}
		}
		// If there aren't enough remaining points to reach the threshold, give up.
		if seen++; thresh-cnt > len(ma)-seen {
			break
		}
	}
	return false
}

type scanner struct {
	id      int
	beacons []point    // original positions
	views   []pointMap // positions in different orientations
}

func newScanner(id int, beacons []point) *scanner {
	sc := &scanner{id, beacons, nil}
	for i := 0; i < 24; i++ {
		sc.views = append(sc.views, make(pointMap))
	}

	// I originally started thinking about rotation matrices, e.g.
	// https://en.wikipedia.org/wiki/Rotation_matrix#Basic_rotations:
	//
	//          |    1      0       0    |
	//  Rx(a) = |    0    cos(a) -sin(a) |
	//          |    0    sin(a)  cos(a) |
	//
	//          |  cos(a)    0    sin(a) |
	//  Ry(a) = |    0       1      0    |
	//          | -sin(a)    0    cos(a) |
	//
	//          |  cos(a) -sin(a)   0    |
	//  Rz(a) = |  sin(a)  cos(a)   0    |
	//          |    0       0      1    |
	//
	// cos(90) = 0, cos(180) = -1, cos(270) =  0
	// sin(90) = 1, sin(180) =  0, sin(270) = -1
	//
	// But then I got lazy and Googled [3d rotation axis 24] and just used the approach described in
	// https://stackoverflow.com/a/16467849 to figure out where the beacons would be with all 24
	// different potential orientations.
	for _, p := range beacons {
		for i := 0; i < 2; i++ {
			for j := 0; j < 3; j++ {
				p = roll(p)
				sc.views[12*i+4*j][p] = struct{}{}
				for k := 1; k < 4; k++ {
					p = turn(p)
					sc.views[12*i+4*j+k][p] = struct{}{}
				}
			}
			p = roll(turn(roll(p)))
		}
	}

	return sc
}

func roll(p point) point { return point{p.x, p.z, -p.y} }
func turn(p point) point { return point{-p.y, p.x, p.z} }
