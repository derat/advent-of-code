package lib

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

// ByteLines returns newline-separated lines from s. Blank lines are skipped.
// If valid is non-empty, panics if any unlisted bytes are encountered.
func ByteLines(s string, valid ...byte) [][]byte {
	var lines [][]byte
	for i, ln := range strings.Split(s, "\n") {
		if len(ln) == 0 {
			continue
		}
		if len(valid) > 0 {
			for j, ch := range ln {
				if bytes.IndexByte(valid, byte(ch)) == -1 {
					panic(fmt.Sprintf("Invalid byte %v (%q) at position %d of line %d", ch, ch, j, i))
				}
			}
		}
		lines = append(lines, []byte(ln))
	}
	return lines
}

// NewBytes returns a 2-dimensional array of ch with r rows and c columns.
func NewBytes(r, c int, ch byte) [][]byte {
	g := make([][]byte, r)
	for i := range g {
		g[i] = bytes.Repeat([]byte{ch}, c)
	}
	return g
}

// CopyBytes returns a deep copy of the supplied 2-dimensional array.
func CopyBytes(b [][]byte) [][]byte {
	return CopyBytesRegion(b, 0, 0, len(b)-1, len(b[0])-1)
}

// CopyBytesRegion returns a copy of the region bounded by (r0, c0) and (r1, c1), inclusive.
func CopyBytesRegion(b [][]byte, r0, c0, r1, c1 int) [][]byte {
	r0, r1 = clampAndSwap(r0, r1, 0, len(b)-1)
	c0, c1 = clampAndSwap(c0, c1, 0, len(b[0])-1)
	n := make([][]byte, r1-r0+1)
	for r := range n {
		n[r] = make([]byte, c1-c0+1)
		copy(n[r], b[r0+r][c0:])
	}
	return n
}

// SetBytes sets the rectangle bounded by (r0, c0) and (r1, c1), inclusive, to ch.
func SetBytes(b [][]byte, r0, c0, r1, c1 int, ch byte) {
	IterBytesRect(b, r0, c0, r1, c1, func(r, c int) { b[r][c] = ch })
}

// CountBytes returns the number of occurrences of chars in the rectangle bounded
// by (r0, c0) and (r1, c1), inclusive. The supplied bounds are clamped.
func CountBytes(b [][]byte, r0, c0, r1, c1 int, chars ...byte) int {
	var cnt int
	IterBytesRect(b, r0, c0, r1, c1, func(r, c int) {
		for _, ch := range chars {
			if b[r][c] == ch {
				cnt++
				return
			}
		}
	})
	return cnt
}

// CountBytesFull calls CountBytes for the full bounds of b.
func CountBytesFull(b [][]byte, chars ...byte) int {
	return CountBytes(b, 0, 0, len(b)-1, len(b[0])-1, chars...)
}

// IterBytesRect calls f for each coordinate in the rectangle bounded by (r0, c0)
// and (r1, c1), inclusive. The supplied bounds are clamped and swapped if needed.
func IterBytesRect(b [][]byte, r0, c0, r1, c1 int, f func(r, c int)) {
	r0, r1 = clampAndSwap(r0, r1, 0, len(b)-1)
	c0, c1 = clampAndSwap(c0, c1, 0, len(b[0])-1)
	for r := r0; r <= r1; r++ {
		for c := c0; c <= c1; c++ {
			f(r, c)
		}
	}
}

// clampAndSwap clamps a and b within [min, max] and swaps them if needed to
// be in ascending order.
func clampAndSwap(a, b, min, max int) (c, d int) {
	c = Max(Min(a, b), min)
	d = Min(Max(a, b), max)
	return c, d
}

// IterBytesLine calls f for each coordinate in the line from (r0, c0) to (r1, c1).
// The supplied points may be outside of b's bounds, but f will only be called for
// in-bounds coordinates.
func IterBytesLine(b [][]byte, r0, c0, r1, c1 int, f func(r, c int)) {
	// This uses the "line drawing on a grid" algorithm that Amit Patel
	// describes at https://www.redblobgames.com/grids/line-drawing.html.
	diagDist := Max(Abs(r1-r0), Abs(c1-c0))
	for step := 0; step <= diagDist; step++ {
		var t float64
		if diagDist > 0 {
			t = float64(step) / float64(diagDist)
		}
		r := int(math.Round(float64(r0) + t*float64(r1-r0)))
		c := int(math.Round(float64(c0) + t*float64(c1-c0)))
		if r >= 0 && r < len(b) && c >= 0 && c < len(b[0]) {
			f(r, c)
		}
	}
}

// DumpBytes returns b as a newline-separated string.
func DumpBytes(b [][]byte) string {
	return string(bytes.Join(b, []byte{'\n'}))
}

// FlipBytesX flips the supplied grid across the X axis.
func FlipBytesX(b [][]byte) {
	for _, r := range b {
		ReverseBytes(r)
	}
}

// FlipBytesY flips the supplied grid across the Y axis.
func FlipBytesY(b [][]byte) {
	Reverse(b)
}

// RotateBytes returns a copy of the supplied grid rotated 90 degrees clockwise.
func RotateBytes(b [][]byte) [][]byte {
	rb := make([][]byte, len(b[0]))
	for r := range rb {
		rb[r] = make([]byte, len(b))
		for c := range rb[r] {
			rb[r][c] = b[len(b)-c-1][r]
		}
	}
	return rb
}
