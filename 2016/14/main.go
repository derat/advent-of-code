package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/derat/advent-of-code/lib"
)

const (
	noThree   = 0xff
	lookahead = 1000
	nkeys     = 64
)

func main() {
	salt := lib.InputLines("2016/14")[0]
	fmt.Println(find(salt, 0)[nkeys-1])
	fmt.Println(find(salt, 2016)[nkeys-1])
}

type num struct {
	idx   int
	three byte
	fives []byte
}

func find(salt string, stretch int) []int {
	nums := make([]num, 0, lookahead)
	numsNext := 0 // next index into nums

	var numsFives [16]int // counts of 5-repeated chars in nums
	var keys []int

	for idx := 0; len(keys) < nkeys; idx++ {
		s := salt + strconv.Itoa(idx)
		hash := md5.Sum([]byte(s))
		for i := 0; i < stretch; i++ {
			hash = md5.Sum([]byte(hex.EncodeToString(hash[:])))
		}

		n := num{idx: idx, three: three(hash[:])}
		if n.three != noThree {
			n.fives = fives(hash[:])
		}

		if len(nums) < lookahead {
			nums = append(nums, n)
		} else {
			old := nums[numsNext]
			for _, v := range old.fives {
				numsFives[v]--
			}

			nums[numsNext] = n
			for _, v := range n.fives {
				numsFives[v]++
			}

			if old.three != noThree && numsFives[old.three] > 0 {
				keys = append(keys, old.idx)
			}
		}
		numsNext = (numsNext + 1) % lookahead
	}

	return keys
}

func three(b []byte) byte {
	for i := range b[:len(b)-1] {
		if (lib.HiIsLo(b[i]) && lib.Hi(b[i]) == lib.Hi(b[i+1])) ||
			(lib.HiIsLo(b[i+1]) && lib.Lo(b[i]) == lib.Lo(b[i+1])) {
			return lib.Lo(b[i])
		}
	}
	return noThree
}

func fives(b []byte) []byte {
	var vals []byte
	for i := range b[:len(b)-2] {
		if (lib.HiIsLo(b[i]) && b[i] == b[i+1] && lib.Hi(b[i]) == lib.Hi(b[i+2])) ||
			(lib.HiIsLo(b[i+1]) && b[i+1] == b[i+2] && lib.Lo(b[i]) == lib.Lo(b[i+1])) {
			vals = append(vals, lib.Lo(b[i+1]))
		}
	}
	return vals
}
