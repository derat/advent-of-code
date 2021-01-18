package lib

import (
	"bytes"
	"fmt"
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
	n := make([][]byte, len(b))
	for r := range b {
		n[r] = append([]byte(nil), b[r]...)
	}
	return n
}

// SetBytes the rectangle bounded by (r0, c0) and (r1, c1), inclusive, to ch.
func SetBytes(b [][]byte, r0, c0, r1, c1 int, ch byte) {
	for r := Max(r0, 0); r <= Min(r1, len(b)-1); r++ {
		for c := Max(c0, 0); c <= Min(c1, len(b[r])-1); c++ {
			b[r][c] = ch
		}
	}
}

// CountBytes returns the number of occurrences of chars in the rectangle bounded
// by (r0, c0) and (r1, c1), inclusive. The supplied bounds are clamped.
func CountBytes(b [][]byte, r0, c0, r1, c1 int, chars ...byte) int {
	cnt := 0
	for r := Max(r0, 0); r <= Min(r1, len(b)-1); r++ {
		for c := Max(c0, 0); c <= Min(c1, len(b[r])-1); c++ {
			for _, ch := range chars {
				if b[r][c] == ch {
					cnt++
					break
				}
			}
		}
	}
	return cnt
}

// CountBytesFull calls CountBytes for the full bounds of b.
func CountBytesFull(b [][]byte, ch byte) int {
	return CountBytes(b, 0, 0, len(b)-1, len(b[0])-1, ch)
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
