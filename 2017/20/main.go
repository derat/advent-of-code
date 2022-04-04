package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	parts := make(map[int]*particle)
	for _, ln := range lib.InputLines("2017/20") {
		vals := lib.ExtractInts(ln)
		lib.AssertEq(len(vals), 9)
		var p particle
		copy(p.pos[:], vals[0:3])
		copy(p.vel[:], vals[3:6])
		copy(p.acc[:], vals[6:9])
		parts[len(parts)] = &p
	}

	// Part 1: Find particle that stays closest to origin.
	minID := -1
	var minPart *particle
	for id, p := range parts {
		if minPart == nil || p.less(minPart) {
			minID = id
			minPart = p
		}
	}
	fmt.Println(minID)

	// Part 2: Find remaining particles after all collisions are resolved.
	for {
		// Compute each particle's updated position.
		seen := make(map[uint64][]int) // packed (x,y,z) to ID
		for id, p := range parts {
			p.update()
			k := lib.PackInts(p.pos[:]...)
			seen[k] = append(seen[k], id)
		}
		// Delete all particles that collided.
		for _, ids := range seen {
			if len(ids) > 1 {
				for _, id := range ids {
					delete(parts, id)
				}
			}
		}
		// Once the particles are in their final ordering in each dimension, we're done.
		if ordered(parts) {
			break
		}
	}
	fmt.Println(len(parts))
}

type particle struct {
	pos, vel, acc [3]int
}

// less returns true if p will remain closer to the origin than o over the long term.
func (p *particle) less(o *particle) bool {
	if pa, oa := mag(p.acc), mag(o.acc); pa < oa {
		return true
	} else if pa > oa {
		return false
	}
	// I'm not convinced that this tie-breaking logic is correct.
	// In my input, there are two particles with the same minimum acceleration.
	// I'm then ordering them by the initial velocity's contribution to the
	// acceleration, before looking at the initial position.
	if pa, oa := mag(add(p.acc, p.vel)), mag(add(o.acc, o.vel)); pa < oa {
		return true
	} else if pa > oa {
		return false
	}
	return mag(p.pos) < mag(o.pos)
}

// update updates p's velocite and acceleration for one tick.
func (p *particle) update() {
	p.vel = add(p.vel, p.acc)
	p.pos = add(p.pos, p.vel)
}

// mag returns the squared magnitude of the supplied vector.
func mag(v [3]int) int {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

// add adds the two supplied vectors.
func add(a, b [3]int) [3]int {
	return [3]int{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

// ordered returns true if the supplied particles' positions have reached their final orderings
// in all dimensions. We sort by ascending position in each dimension and then check than the
// corresponding velocities and accelerations are also ordered. Once they are, the particles won't
// able to overtake each other.
func ordered(parts map[int]*particle) bool {
	ids := lib.MapKeys(parts)
	for dim := 0; dim < 3; dim++ {
		sort.Slice(ids, func(a, b int) bool { return parts[ids[a]].pos[dim] < parts[ids[b]].pos[dim] })
		for i, id := range ids[:len(ids)-1] {
			if p, n := parts[id], parts[ids[i+1]]; p.vel[dim] > n.vel[dim] || p.acc[dim] > n.acc[dim] {
				return false
			}
		}
	}
	return true
}
