package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nplayers, maxMarble int
	lib.Extract(lib.InputLines("2018/9")[0],
		`^(\d+) players; last marble is worth (\d+) points$`, &nplayers, &maxMarble)

	cur := &marble{val: 0}
	cur.left = cur
	cur.right = cur

	scores := make([]int, nplayers)
	var idx int // current player

	for val := 1; val <= maxMarble; val, idx = val+1, (idx+1)%len(scores) {
		if val%23 != 0 {
			// Insert marble between marbles 1 and 2 clockwise (left) of current.
			mar := &marble{
				val:   val,
				left:  cur.left.left,
				right: cur.left,
			}
			cur.left.left.right = mar
			cur.left.left = mar
			cur = mar
		} else {
			scores[idx] += val

			// Remove marble 7 counter-clockwise (right) of current and add its
			// value to the current player's score. The marble clockwise (left)
			// of it becomes current.
			rem := cur
			for i := 0; i < 7; i++ {
				rem = rem.right
			}
			scores[idx] += rem.val
			rem.left.right = rem.right
			rem.right.left = rem.left
			cur = rem.left
		}
	}

	// Part 1: Print winning player's score.
	fmt.Println(lib.Max(scores...))
}

type marble struct {
	val         int
	left, right *marble // left is cw, right is ccw
}
