package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var bots []bot
	for _, ln := range lib.InputLines("2018/23") {
		var b bot
		lib.Extract(ln, `^pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)$`, &b.c.x, &b.c.y, &b.c.z, &b.r)
		bots = append(bots, b)
	}

	// Part 1: Find bot with largest signal radius and count bots in range of its signals.
	var maxb bot
	for _, b := range bots {
		if b.r > maxb.r {
			maxb = b
		}
	}
	var cnt int
	for _, b := range bots {
		if maxb.check(b.c) {
			cnt++
		}
	}
	fmt.Println(cnt)

	// Part 2: Find the shortest Manhattan distance to the origin of a point that's in
	// range of the most bots possible.
	//
	// I had a huge amount of trouble with this. I spent a lot of time thinking about this
	// problem in two dimensions, where there seem to be some properties that may not hold
	// in three dimensions:
	//
	// - If two bots have overlapping ranges, then at least one corner of one bot's range
	//   should be in range of the other bot.
	// - For a clique of bots with overlapping ranges, at least one corner of one bot's
	//   range should be in range of all the bots in the clique.
	//
	// If this were true in three dimensions, I think that this problem would be as simple
	// as looping through all of the corners and choosing the one with the maximum number
	// of matches with other bots, and then somehow exploring that corner's region to find
	// the point with the shortest Manhattan distance to the origin.
	//
	// When I tried that (just submitting the corner's distance), I ended up with an
	// answer that was too low, suggesting that this approach doesn't actually work.
	//
	// The first thing that I tried was exploring the space with large steps to try to find
	// the maximum number of in-range bots, but I was just testing individual points instead
	// of regions (because I didn't know how to test whether any portion of a region is in
	// range of a given bot).
	//
	// I also read a bit about clique-finding and used Bron-Kerbosch to find the maximum
	// clique. Even with that, I'd still need to find the points within any of those bots'
	// ranges that were in range of all of them, though. I spent a lot of time spinning my
	// wheels with the abovementioned misconception that one of the ranges' corners would
	// be in range of all the bots.
	//
	// After checking Reddit to see that there wasn't any trick that I was missing (I'd
	// already tried e.g. removing entirely-overlapping regions), I finally realized how
	// easy it is to test entire region, sos I combined that with the clique approach to
	// get the answer. I think that finding a clique is unnecessary; one could just follow
	// the regions with the maximum number of in-range bots when performing the search.

	// Find other bots with overlapping signal ranges for each bot.
	overlaps := make(map[int]map[int]struct{}, len(bots)) // keys are indexes into bots
	for i := range bots {
		overlaps[i] = make(map[int]struct{})
	}
	for i := range bots {
		for j := i + 1; j < len(bots); j++ {
			if dist(bots[i].c, bots[j].c) <= bots[i].r+bots[j].r {
				overlaps[i][j] = struct{}{}
				overlaps[j][i] = struct{}{}
			}
		}
	}

	// Find the clique containing the maximum number of bots (NP-hard).
	// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm
	// This usually runs quickly, but I've seen it be slow occasionally, presumably because
	// it doesn't find the huge clique early on and wastes its time exploring a bunch of
	// smaller ones. It'd probably be possible to bail out early on small cliques to avoid
	// this (making the assumption that all inputs have a big clique).
	var inc map[int]struct{}
	var bronKerbosch func(r, p, x map[int]struct{})
	bronKerbosch = func(r, p, x map[int]struct{}) {
		if len(p) == 0 && len(x) == 0 {
			if len(r) > len(inc) {
				inc = r
			}
			return
		}
		for v := range p {
			bronKerbosch(
				lib.Union(r, map[int]struct{}{v: struct{}{}}).(map[int]struct{}),
				lib.Intersect(p, overlaps[v]).(map[int]struct{}),
				lib.Intersect(x, overlaps[v]).(map[int]struct{}))
			// Early exit once we've found the maximum clique.
			if len(inc) > len(bots)/2 {
				return
			}
			delete(p, v)
			x = lib.Union(x, map[int]struct{}{v: struct{}{}}).(map[int]struct{})
		}
	}

	rem := make(map[int]struct{}, len(bots))
	for i := range bots {
		rem[i] = struct{}{}
	}
	bronKerbosch(nil, rem, nil)
	lib.Assertf(len(inc) > 0, "Didn't find any cliques")

	// Assertion: each bot in the clique should overlap with every other bot.
	for i := range inc {
		for j := range inc {
			if i != j {
				lib.Assert(lib.MapHasKey(overlaps[i], j))
				lib.Assert(lib.MapHasKey(overlaps[j], i))
			}
		}
	}

	// Get the bounds of the centers of all of the bots in the clique.
	pts := make([]point, len(inc))
	for i := range inc {
		pts = append(pts, bots[i].c)
	}
	min, max := bounds(pts)

	bestDist := math.MaxInt32
	var bestPoint point

	queue := [][2]point{{min, max}}
	for len(queue) > 0 {
		var next [][2]point
		seen := make(map[string]struct{})

		for _, region := range queue {
			min, max := region[0], region[1]
			mid := point{
				min.x + (max.x-min.x)/2,
				min.y + (max.y-min.y)/2,
				min.z + (max.z-min.z)/2,
			}

			for o := 0; o < 8; o++ {
				var omin, omax point
				if o&4 == 0 || min.x == max.x {
					omin.x, omax.x = min.x, mid.x
				} else {
					omin.x, omax.x = mid.x+1, max.x
				}
				if o&2 == 0 || min.y == max.y {
					omin.y, omax.y = min.y, mid.y
				} else {
					omin.y, omax.y = mid.y+1, max.y
				}
				if o&1 == 0 || min.z == max.z {
					omin.z, omax.z = min.z, mid.z
				} else {
					omin.z, omax.z = mid.z+1, max.z
				}

				good := true
				for i := range inc {
					if b := bots[i]; !b.check(point{
						lib.Min(lib.Max(b.c.x, omin.x), omax.x),
						lib.Min(lib.Max(b.c.y, omin.y), omax.y),
						lib.Min(lib.Max(b.c.z, omin.z), omax.z),
					}) {
						good = false
						break
					}
				}
				if !good {
					continue
				}

				key := fmt.Sprint(omin, omax)
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}

				if omin == omax {
					d := dist(omin, point{})
					if d < bestDist {
						bestDist = d
						bestPoint = omin
					}
				} else {
					next = append(next, [2]point{omin, omax})
				}
			}
		}
		queue = next
	}

	for i := range inc {
		lib.Assertf(bots[i].check(bestPoint), "%v not in range of bot %d (%v)", bestPoint, i)
	}
	fmt.Println(bestDist)
}

type point struct {
	x, y, z int
}

// dist returns the Manhattan distance between a and b.
func dist(a, b point) int {
	return lib.Abs(a.x-b.x) + lib.Abs(a.y-b.y) + lib.Abs(a.z-b.z)
}

type bot struct {
	c point
	r int
}

// check returns true if p is within b's signal range.
func (b *bot) check(p point) bool {
	return dist(p, b.c) <= b.r
}

// bounds returns a pair of points with the min and max coords from the supplied points.
func bounds(points []point) (point, point) {
	min := point{math.MaxInt32, math.MaxInt32, math.MaxInt32}
	max := point{math.MinInt32, math.MinInt32, math.MinInt32}
	for _, p := range points {
		min.x, max.x = lib.Min(min.x, p.x), lib.Max(max.x, p.x)
		min.y, max.y = lib.Min(min.y, p.y), lib.Max(max.y, p.y)
		min.z, max.z = lib.Min(min.z, p.z), lib.Max(max.z, p.z)
	}
	return min, max
}
