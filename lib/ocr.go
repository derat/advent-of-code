package lib

import (
	"fmt"
	"math"
)

var canonGlyphs map[byte][]glyph

// OCR does a terrible job of trying to recognize glyphs in b.
// Glyphs must be in a single row with one or more blank columns between them.
func OCR(b [][]byte, blank byte) string {
	if canonGlyphs == nil {
		canonGlyphs = make(map[byte][]glyph, len(canonGlyphData))
		for ch, s := range canonGlyphData {
			canonGlyphs[ch] = extractGlyphs(NewByteGridString(s, '#', ' '), ' ')
		}
	}

	var rec string
	for _, gl := range extractGlyphs(b, blank) {
		var bestCh byte
		var bestScore float64
		for ch, cgs := range canonGlyphs {
			for _, cg := range cgs {
				if sc := gl.compare(&cg); sc > bestScore {
					bestCh = ch
					bestScore = sc
				}
			}
		}
		rec += string(rune(bestCh))
	}
	return rec
}

// extractGlyphs extracts glyphs from b.
func extractGlyphs(b ByteGrid, blank byte) []glyph {
	rows, cols := len(b), len(b[0])
	for r := range b[1:] {
		AssertEq(len(b[r]), cols)
	}

	// Identify left and right bounds of glyphs by looking for empty columns.
	var glyphs []glyph
	var inGlyph bool
	for c := range b[0] {
		empty := b.CountRect(0, c, rows-1, c, blank) == rows
		if inGlyph && empty {
			inGlyph = false
		} else if inGlyph && !empty {
			glyphs[len(glyphs)-1].right = c
		} else if !inGlyph && !empty {
			inGlyph = true
			glyphs = append(glyphs, glyph{left: c, right: c, features: make(map[string]float64)})
		}
	}

	// Identify top and bottom bounds by looking for nonempty rows.
	for i := range glyphs {
		gl := &glyphs[i]
		for r := 0; r < rows; r++ {
			if b.CountRect(r, gl.left, r, gl.right, blank) != gl.width() {
				gl.top = r
				break
			}
		}
		for r := rows - 1; r >= 0; r-- {
			if b.CountRect(r, gl.left, r, gl.right, blank) != gl.width() {
				gl.bottom = r
				break
			}
		}
		// Allow narrow widths to handle e.g. single-column 'I' glyphs.
		if gl.height() < 3 {
			panic(fmt.Sprintf("Glyph %d has invalid vertical bounds [%d, %d] with horizontal bounds [%d, %d]",
				i, gl.top, gl.bottom, gl.left, gl.right))
		}
		gl.addFeatures(b, blank)
	}

	return glyphs
}

type glyph struct {
	left, right int
	top, bottom int
	features    map[string]float64
}

func (g *glyph) width() int {
	return g.right - g.left + 1
}

func (g *glyph) height() int {
	return g.bottom - g.top + 1
}

func (g *glyph) addFeatures(b ByteGrid, blank byte) {
	// Returns the fraction of non-blank bytes in the supplied inclusive bounds.
	filled := func(r0, c0, r1, c1 int) float64 {
		empty := b.CountRect(r0, c0, r1, c1, blank)
		total := (r1 - r0 + 1) * (c1 - c0 + 1)
		return float64(total-empty) / float64(total)
	}

	g.features["corner-0,0"] = filled(g.top, g.left, g.top, g.left)
	g.features["corner-0,1"] = filled(g.top, g.right, g.top, g.right)
	g.features["corner-1,0"] = filled(g.bottom, g.left, g.bottom, g.left)
	g.features["corner-1,1"] = filled(g.bottom, g.right, g.bottom, g.right)

	g.features["edge-left"] = filled(g.top, g.left, g.bottom, g.left)
	g.features["edge-right"] = filled(g.top, g.right, g.bottom, g.right)
	g.features["edge-top"] = filled(g.top, g.left, g.top, g.right)
	g.features["edge-bottom"] = filled(g.bottom, g.left, g.bottom, g.right)

	// Splits [a, b] into three inclusive [start, end] pairs (possibly nil).
	thirds := func(a, b int) [][]int {
		n := (b - a + 1) / 3
		if n == 0 {
			return [][]int{nil, {a, b}, nil}
		}
		return [][]int{{a, a + n - 1}, {a + n, b - n}, {b - n + 1, b}}
	}
	for i, n := range thirds(g.top, g.bottom) {
		for j, o := range thirds(g.left, g.right) {
			key := fmt.Sprintf("density-%d,%d", i, j)
			if n == nil || o == nil {
				g.features[key] = 0
			} else {
				g.features[key] = filled(n[0], o[0], n[1], o[1])
			}
		}
	}
}

func (g *glyph) compare(o *glyph) float64 {
	var sum float64
	for k, v := range g.features {
		// TODO: Weight different features differently?
		diff := math.Abs(v - o.features[k])
		sum += (1 - diff)
	}
	return sum / float64(len(g.features))
}
