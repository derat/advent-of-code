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

// CountBytes returns the number of occurrences of ch in the rectangle bounded
// by (r0, c0) and (r1, c1), inclusive. The supplied bounds are clamped.
func CountBytes(b [][]byte, r0, c0, r1, c1 int, ch byte) int {
	cnt := 0
	for r := Max(r0, 0); r <= Min(r1, len(b)-1); r++ {
		for c := Max(c0, 0); c <= Min(c1, len(b[r])-1); c++ {
			if b[r][c] == ch {
				cnt++
			}
		}
	}
	return cnt
}

// CountBytesFull calls CountBytes for the full bounds of b.
func CountBytesFull(b [][]byte, ch byte) int {
	return CountBytes(b, 0, 0, len(b)-1, len(b[0])-1, ch)
}
