package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	w := newWorld()
	var start, lines int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			continue
		}
		if lines == 0 {
			start = len(ln) / 2 // center initial state
		}
		for i, ch := range ln {
			if ch == '#' {
				w.set(-start+i, start-lines, 0)
			}
		}
		lines++
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	//fmt.Println("Before any cycles:\n")
	//w.dump()

	count := func(x, y, z int) int {
		n := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				for dz := -1; dz <= 1; dz++ {
					if dx == 0 && dy == 0 && dz == 0 {
						continue
					}
					if w.active(x+dx, y+dy, z+dz) {
						n++
					}
				}
			}
		}
		return n
	}

	for i := 1; i <= 6; i++ {
		w2 := newWorld()
		w.forEach(func(x, y, z int) {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					for dz := -1; dz <= 1; dz++ {
						x0, y0, z0 := x+dx, y+dy, z+dz
						cnt := count(x0, y0, z0)
						active := w.active(x0, y0, z0)
						if active && (cnt == 2 || cnt == 3) {
							w2.set(x0, y0, z0)
						} else if !active && cnt == 3 {
							w2.set(x0, y0, z0)
						}
					}
				}
			}
		})

		w = w2
		//fmt.Printf("After %d cycle(s):\n\n", i)
		//w.dump()
	}

	cnt := 0
	w.forEach(func(x, y, z int) { cnt++ })
	fmt.Println(cnt)
}

type world struct {
	// Values represent 4x4x4 grids of active states for cubes.
	maps [8]map[uint64]uint64

	xrange, yrange, zrange [2]int
}

func newWorld() world {
	var w world
	for i := range w.maps {
		w.maps[i] = make(map[uint64]uint64)
	}
	return w
}

func (w *world) set(x, y, z int) {
	midx, mkey, mask := idx(x, y, z)
	w.maps[midx][mkey] |= mask

	if x < w.xrange[0] {
		w.xrange[0] = x
	} else if x > w.xrange[1] {
		w.xrange[1] = x
	}
	if y < w.yrange[0] {
		w.yrange[0] = y
	} else if y > w.yrange[1] {
		w.yrange[1] = y
	}
	if z < w.zrange[0] {
		w.zrange[0] = z
	} else if z > w.zrange[1] {
		w.zrange[1] = z
	}
}

func (w *world) active(x, y, z int) bool {
	midx, mkey, mask := idx(x, y, z)
	return w.maps[midx][mkey]&mask != 0
}

const (
	kbits = 21
	kmask = (1 << kbits) - 1
)

func (w *world) forEach(f func(x, y, z int)) {
	for mi, m := range w.maps {
		for mk, gr := range m {
			for bit := 0; bit < 64; bit++ {
				if gr&(1<<bit) != 0 {
					x, y, z := unIdx(mi, mk, bit)
					f(x, y, z)
				}
			}
		}
	}
}

func (w *world) dump() {
	for z := w.zrange[0]; z <= w.zrange[1]; z++ {
		fmt.Printf("z=%d\n", z)
		for y := w.yrange[1]; y >= w.yrange[0]; y-- {
			for x := w.xrange[0]; x <= w.xrange[1]; x++ {
				if w.active(x, y, z) {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func idx(x, y, z int) (mapIdx int, mapKey, mask uint64) {
	mapIdx = (neg(x) << 2) + (neg(y) << 1) + neg(z)
	ax, ay, az := uint64(abs(x)), uint64(abs(y)), uint64(abs(z))
	mapKey = ((ax / 4) << (2 * kbits)) | ((ay / 4) << kbits) | (az / 4)
	mask = 1 << ((ax%4)*16 + (ay%4)*4 + az%4)
	return mapIdx, mapKey, mask
}

func unIdx(mapIdx int, mapKey uint64, bit int) (x, y, z int) {
	x = int(((mapKey >> (2 * kbits)) & kmask) * 4)
	y = int(((mapKey >> kbits) & kmask) * 4)
	z = int((mapKey & kmask) * 4)

	x += bit / 16
	y += (bit % 16) / 4
	z += bit % 4

	if mapIdx&4 != 0 {
		x = -x
	}
	if mapIdx&2 != 0 {
		y = -y
	}
	if mapIdx&1 != 0 {
		z = -z
	}
	return x, y, z
}

// neg returns 1 if n is negative and 0 otherwise.
func neg(n int) int {
	if n < 0 {
		return 1
	}
	return 0
}

// abs returns the absolute value of n.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
