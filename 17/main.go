package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	s1 := newSpace()
	s2 := newSpace()
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
				s1.set(-start+i, start-lines, 0, 0)
				s2.set(-start+i, start-lines, 0, 0)
			}
		}
		lines++
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	//fmt.Println("Before any cycles:\n")
	//s2.dump()

	count := func(s space, x, y, z, w int) int {
		n := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				for dz := -1; dz <= 1; dz++ {
					for dw := -1; dw <= 1; dw++ {
						if dx == 0 && dy == 0 && dz == 0 && dw == 0 {
							continue
						}
						if s.active(x+dx, y+dy, z+dz, w+dw) {
							n++
						}
					}
				}
			}
		}
		return n
	}

	for i := 1; i <= 6; i++ {
		ns1 := newSpace()
		ns2 := newSpace()

		s1.forEach(func(x, y, z, _ int) {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					for dz := -1; dz <= 1; dz++ {
						x0, y0, z0 := x+dx, y+dy, z+dz
						cnt := count(s1, x0, y0, z0, 0)
						active := s1.active(x0, y0, z0, 0)
						if active && (cnt == 2 || cnt == 3) {
							ns1.set(x0, y0, z0, 0)
						} else if !active && cnt == 3 {
							ns1.set(x0, y0, z0, 0)
						}
					}
				}
			}
		})

		s2.forEach(func(x, y, z, w int) {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					for dz := -1; dz <= 1; dz++ {
						for dw := -1; dw <= 1; dw++ {
							x0, y0, z0, w0 := x+dx, y+dy, z+dz, w+dw
							cnt := count(s2, x0, y0, z0, w0)
							active := s2.active(x0, y0, z0, w0)
							if active && (cnt == 2 || cnt == 3) {
								ns2.set(x0, y0, z0, w0)
							} else if !active && cnt == 3 {
								ns2.set(x0, y0, z0, w0)
							}
						}
					}
				}
			}
		})

		s1, s2 = ns1, ns2
		//fmt.Printf("After %d cycle(s):\n\n", i)
		//s2.dump()
	}

	cnt := 0
	s1.forEach(func(x, y, z, _ int) { cnt++ })
	fmt.Println(cnt)

	cnt = 0
	s2.forEach(func(x, y, z, w int) { cnt++ })
	fmt.Println(cnt)
}

type space struct {
	// Values represent 2x2x2x2 grids of active states for cubes.
	maps [16]map[uint64]uint16

	xrange, yrange, zrange, wrange [2]int
}

func newSpace() space {
	var s space
	for i := range s.maps {
		s.maps[i] = make(map[uint64]uint16)
	}
	return s
}

func (s *space) set(x, y, z, w int) {
	midx, mkey, mask := idx(x, y, z, w)
	s.maps[midx][mkey] |= mask

	if x < s.xrange[0] {
		s.xrange[0] = x
	} else if x > s.xrange[1] {
		s.xrange[1] = x
	}
	if y < s.yrange[0] {
		s.yrange[0] = y
	} else if y > s.yrange[1] {
		s.yrange[1] = y
	}
	if z < s.zrange[0] {
		s.zrange[0] = z
	} else if z > s.zrange[1] {
		s.zrange[1] = z
	}
	if w < s.zrange[0] {
		s.wrange[0] = w
	} else if z > s.wrange[1] {
		s.wrange[1] = w
	}
}

func (s *space) active(x, y, z, w int) bool {
	midx, mkey, mask := idx(x, y, z, w)
	return s.maps[midx][mkey]&mask != 0
}

const (
	kbits = 16
	kmask = (1 << kbits) - 1
)

func (s *space) forEach(f func(x, y, z, w int)) {
	for mi, m := range s.maps {
		for mk, gr := range m {
			for bit := 0; bit < 16; bit++ {
				if gr&(1<<bit) != 0 {
					x, y, z, w := unIdx(mi, mk, bit)
					f(x, y, z, w)
				}
			}
		}
	}
}

func (s *space) dump() {
	for w := s.wrange[0]; w <= s.wrange[1]; w++ {
		for z := s.zrange[0]; z <= s.zrange[1]; z++ {
			fmt.Printf("z=%d, w=%d\n", z, w)
			for y := s.yrange[1]; y >= s.yrange[0]; y-- {
				for x := s.xrange[0]; x <= s.xrange[1]; x++ {
					if s.active(x, y, z, w) {
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
}

func idx(x, y, z, w int) (mapIdx int, mapKey uint64, mask uint16) {
	mapIdx = (neg(x) << 3) + (neg(y) << 2) + (neg(z) << 1) + neg(w)
	ax, ay, az, aw := uint16(abs(x)), uint16(abs(y)), uint16(abs(z)), uint16(abs(w))
	mapKey = ((uint64(ax) / 2) << (3 * kbits)) | ((uint64(ay) / 2) << (2 * kbits)) |
		((uint64(az) / 2) << kbits) | (uint64(aw) / 2)
	mask = 1 << ((ax%2)*8 + (ay%2)*4 + (az%2)*2 + (aw % 2))
	return mapIdx, mapKey, mask
}

func unIdx(mapIdx int, mapKey uint64, bit int) (x, y, z, w int) {
	x = int(((mapKey >> (3 * kbits)) & kmask) * 2)
	y = int(((mapKey >> (2 * kbits)) & kmask) * 2)
	z = int(((mapKey >> kbits) & kmask) * 2)
	w = int((mapKey & kmask) * 2)

	x += (bit / 8) % 2
	y += (bit / 4) % 2
	z += (bit / 2) % 2
	w += bit % 2

	if mapIdx&8 != 0 {
		x = -x
	}
	if mapIdx&4 != 0 {
		y = -y
	}
	if mapIdx&2 != 0 {
		z = -z
	}
	if mapIdx&1 != 0 {
		w = -w
	}
	return x, y, z, w
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
