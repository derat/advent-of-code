package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	packages := lib.InputInts("2015/24")
	lib.AssertLessEq(len(packages), 32) // packed into 4 8-bit tables

	total := lib.Sum(packages...)
	lib.AssertEq(total%3, 0)
	each := total / 3

	// Precompute the summed weights, counts, and quantum entanglements of the 256
	// combinations of each group of 8 packages. sums[0][0] corresponds to no packages,
	// sums[0][1] corresponds to package[0], sums[0][3] corresponds to packages[0] and
	// packages[1], and so on.
	var sums [4][256]int
	var cnts [4][256]int8
	var qes [4][256]int64

	for i := 0; i < 4; i++ {
		for j := 0; j < 256; j++ {
			qes[i][j] = 1
			for k := 0; k < 8 && i*8+k < len(packages); k++ {
				if j&(1<<k) != 0 {
					w := packages[i*8+k]
					sums[i][j] += w
					qes[i][j] *= int64(w)
					cnts[i][j]++
				}
			}
		}
	}

	var bestCnt int8 = math.MaxInt8
	var bestQE int64 = math.MaxInt64
	var bestUsed uint64

	var packFirst func(int, int, int, int8, int64, uint64)
	var packOther func(int, int, int, uint64) bool

	// packFirst updates bestCnt, bestQE, and bestUsed with the combination of packages summing
	// to remain with the lowest total count and (for tie-breaking) lowest quantum entanglement.
	//
	// idx is the current index into sums/cnts/qes.
	// remain is the remaining weight that's needed.
	// repeat specifies number of containers beyond the original 3 that need to be packed.
	// cnt is the number of packages used so far.
	// qe is the quantum entanglement so far.
	// used is a bitfield specifying the packages that have been used (e.g. 0x1 is packages[0]).
	packFirst = func(idx, remain, repeat int, cnt int8, qe int64, used uint64) {
		if idx == len(sums) {
			// If this combination has the potential to be a new winner, see if we can pack
			// the second group.
			if remain == 0 && (cnt < bestCnt || (cnt == bestCnt && qe < bestQE)) {
				if packOther(0, each, repeat, used) {
					bestCnt = cnt
					bestQE = qe
					bestUsed = used
				}
			}
			return
		}

		for i, w := range sums[idx] {
			if w <= remain {
				packFirst(idx+1, remain-w, repeat, cnt+cnts[idx][i], qe*qes[idx][i], used|(uint64(i)<<(8*idx)))
			}
		}
	}

	// packOther returns true if there exists any combination of not-yet-used packages
	// with weights summing to remain.
	packOther = func(idx, remain, repeat int, used uint64) bool {
		if idx == len(sums) {
			if remain != 0 {
				return false
			}
			if repeat > 0 {
				return packOther(0, each, repeat-1, used)
			}
			return true
		}

		umask := used >> (8 * idx) & 0xff
		for i, w := range sums[idx] {
			if i&int(umask) != i {
				continue // combo uses packages that have already been used
			}
			if w <= remain {
				if packOther(idx+1, remain-w, repeat, used|(uint64(i)<<(8*idx))) {
					return true
				}
			}
		}
		return false
	}

	// Part 1: Pack 3 containers.
	packFirst(0, each, 0, 0, 1, 0)
	fmt.Println(bestQE)

	// Part 2: Pack 4 containers.
	lib.AssertEq(total%4, 0)
	each = total / 4
	bestCnt = math.MaxInt8
	bestQE = math.MaxInt64
	bestUsed = 0
	packFirst(0, each, 1, 0, 1, 0)
	fmt.Println(bestQE)
}
