package lib

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

// NewByteLines returns newline-separated lines from s. Blank lines are skipped.
// If valid is non-empty, panics if any unlisted bytes are encountered.
func NewByteLines(s string, valid ...byte) [][]byte {
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

// ByteGrid holds a two-dimensional grid of bytes.
type ByteGrid [][]byte

// NewByteGrid returns a 2-dimensional array of ch with r rows and c columns.
func NewByteGrid(r, c int, ch byte) ByteGrid {
	AssertLess(0, r)
	AssertLess(0, c)
	g := make([][]byte, r)
	for i := range g {
		g[i] = bytes.Repeat([]byte{ch}, c)
	}
	return ByteGrid(g)
}

// NewByteGrid creates a ByteGrid containing only the supplied row.
func NewByteGridRow(r []byte) ByteGrid {
	AssertLess(0, len(r))
	return ByteGrid([][]byte{r})
}

// NewByteGridString splits s into rows per NewByteLines and returns a ByteGrid.
func NewByteGridString(s string, valid ...byte) ByteGrid {
	lines := NewByteLines(s, valid...)
	AssertLess(0, len(lines))
	AssertLess(0, len(lines[0]))
	for i := 0; i < len(lines); i++ {
		AssertEq(len(lines[0]), len(lines[1]))
	}
	return ByteGrid(lines)
}

func (b ByteGrid) Rows() int { return len(b) }
func (b ByteGrid) Cols() int { return len(b[0]) }

func (b ByteGrid) MaxRow() int { return b.Rows() - 1 }
func (b ByteGrid) MaxCol() int { return b.Cols() - 1 }

func (b ByteGrid) InBounds(r, c int) bool {
	return r >= 0 && c >= 0 && r <= b.MaxRow() && c <= b.MaxCol()
}

// Copy returns a deep copy of b.
func (b ByteGrid) Copy() ByteGrid {
	return b.CopyRect(0, 0, b.MaxRow(), b.MaxCol())
}

// CopyRect returns a copy of the rectangle bounded by (r0, c0) and (r1, c1), inclusive.
func (b ByteGrid) CopyRect(r0, c0, r1, c1 int) ByteGrid {
	r0, r1 = clampAndSwap(r0, r1, 0, b.MaxRow())
	c0, c1 = clampAndSwap(c0, c1, 0, b.MaxCol())
	n := make([][]byte, r1-r0+1)
	for r := range n {
		n[r] = make([]byte, c1-c0+1)
		copy(n[r], b[r0+r][c0:])
	}
	return n
}

// Get returns the byte at (r, c). If the coordinates are outside b's bounds,
// def is returned instead.
func (b ByteGrid) Get(r, c int, def byte) byte {
	if r < 0 || r > b.MaxRow() || c < 0 || c > b.MaxCol() {
		return def
	}
	return b[r][c]
}

// SetRect sets the rectangle bounded by (r0, c0) and (r1, c1), inclusive, to ch.
func (b ByteGrid) SetRect(r0, c0, r1, c1 int, ch byte) {
	b.IterRect(r0, c0, r1, c1, func(r, c int) { b[r][c] = ch })
}

// CountRect returns the number of occurrences of chars in the rectangle bounded
// by (r0, c0) and (r1, c1), inclusive. The supplied bounds are clamped.
func (b ByteGrid) CountRect(r0, c0, r1, c1 int, chars ...byte) int {
	var cnt int
	b.IterRect(r0, c0, r1, c1, func(r, c int) {
		for _, ch := range chars {
			if b[r][c] == ch {
				cnt++
				return
			}
		}
	})
	return cnt
}

// Count returns the number of occurrences of chars in b.
func (b ByteGrid) Count(chars ...byte) int {
	return b.CountRect(0, 0, b.MaxRow(), b.MaxCol(), chars...)
}

// Iter calls f for each coordinate in b.
func (b ByteGrid) Iter(f func(r, c int)) {
	b.IterRect(0, 0, b.Rows(), b.Cols(), f)
}

// IterRect calls f for each coordinate in the rectangle bounded by (r0, c0)
// and (r1, c1), inclusive. The supplied bounds are clamped and swapped if needed.
func (b ByteGrid) IterRect(r0, c0, r1, c1 int, f func(r, c int)) {
	r0, r1 = clampAndSwap(r0, r1, 0, b.MaxRow())
	c0, c1 = clampAndSwap(c0, c1, 0, b.MaxCol())
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

// IterLine calls f for each coordinate in the line from (r0, c0) to (r1, c1).
// The supplied points may be outside of b's bounds, but f will only be called for
// in-bounds coordinates.
func (b ByteGrid) IterLine(r0, c0, r1, c1 int, f func(r, c int)) {
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
		if r >= 0 && r < b.Rows() && c >= 0 && c < b.Cols() {
			f(r, c)
		}
	}
}

// Dump returns b as a newline-separated string.
func (b ByteGrid) Dump() string {
	return string(bytes.Join(b, []byte{'\n'}))
}

// FlipHoriz returns a copy of b reflected across the Y axis.
func (b ByteGrid) FlipHoriz() ByteGrid {
	n := b.Copy()
	for _, r := range n {
		ReverseBytes(r)
	}
	return n
}

// FlipVert returns a copy of b reflected across the X axis.
func (b ByteGrid) FlipVert() ByteGrid {
	n := b.Copy()
	Reverse(n)
	return n
}

// RotateCW returns a copy of b rotated 90 degrees clockwise.
func (b ByteGrid) RotateCW() ByteGrid {
	rb := make([][]byte, b.Cols())
	for r := range rb {
		rb[r] = make([]byte, b.Rows())
		for c := range rb[r] {
			rb[r][c] = b[b.Rows()-c-1][r]
		}
	}
	return ByteGrid(rb)
}
