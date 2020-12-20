package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var dim int // width and height (tiles are square)
	tiles := make(map[int]tile)
	addTile := func(id int, rows []uint16) {
		if len(rows) != dim {
			log.Fatalf("Tile %d has %d row(s); want %d", id, len(rows), dim)
		}
		tiles[id] = newTile(rows)
	}

	headRegexp := regexp.MustCompile(`^Tile (\d+):$`)
	var id int
	var rows []uint16
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			continue
		}

		// Start of new tile.
		if m := headRegexp.FindStringSubmatch(ln); m != nil {
			if id > 0 {
				addTile(id, rows)
				rows = make([]uint16, 0, dim)
			}
			var err error
			if id, err = strconv.Atoi(m[1]); err != nil {
				log.Fatalf("Failed to parse ID in %q: %v", ln, err)
			}
			continue
		}

		// Row within current tile.
		if dim == 0 {
			if len(ln) > 16 {
				log.Fatalf("Tile dimension %d > 16", len(ln))
			}
			dim = len(ln)
		} else if len(ln) != dim {
			log.Fatalf("Tile %d row %q has width %d; want %d", id, ln, len(ln), dim)
		}

		var row uint16
		for _, ch := range ln {
			row <<= 1
			switch ch {
			case '.':
			case '#':
				row |= 1
			default:
				log.Fatalf("Tile %d row %q contains bad char %q", id, ln, ch)
			}
		}
		rows = append(rows, row)
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	if id > 0 {
		addTile(id, rows)
	}

	/*
		ti := tiles[2311]
		for _, r := range ti.rows {
			fmt.Println(dumpRow(r, dim))
		}
		fmt.Println()
		for _, edges := range ti.edges {
			fmt.Println(dumpRow(edges[top], dim), " t")
			fmt.Println(dumpRow(edges[bottom], dim), " b")
			fmt.Println(dumpRow(edges[left], dim), " l")
			fmt.Println(dumpRow(edges[right], dim), " r")
			fmt.Println()
		}
	*/

	// Returns a map from ID to transformations with a matching edge on the supplied side.
	find := func(edge uint16, side int) map[int][]int {
		matches := make(map[int][]int)
		for id, ti := range tiles {
			for tr, edges := range ti.edges {
				if edges[side] == edge {
					matches[id] = append(matches[id], tr)
				}
			}
		}
		return matches
	}

	var corners []int
	for id, ti := range tiles {
		for _, edges := range ti.edges {
			unique := func(edge uint16, side int) bool {
				m := find(edge, side)
				delete(m, id)
				return len(m) == 0
			}
			t := unique(edges[top], bottom)
			b := unique(edges[bottom], top)
			l := unique(edges[right], left)
			r := unique(edges[left], right)
			if (t || b) && (l || r) {
				corners = append(corners, id)
				break
			}
		}
	}
	if len(corners) != 4 {
		log.Fatal("Want 4 corners; got %v", corners)
	}
	var prod int64 = 1
	for _, id := range corners {
		prod *= int64(id)
	}
	fmt.Println(prod)
}

const (
	top    = 0
	bottom = 1
	left   = 2
	right  = 3
)

type tile struct {
	edges [8][4]uint16 // outer is transformation, inner is edge (LSB right/bottom)
	rows  []uint16     // LSB right
}

func newTile(rows []uint16) tile {
	dim := len(rows)
	ti := tile{rows: append([]uint16{}, rows...)}
	for tr := 0; tr < 8; tr++ {
		trows := append([]uint16{}, rows...)
		if tr&1 != 0 { // flip Y: https://stackoverflow.com/a/28058324/6882947
			for i, j := 0, len(trows)-1; i < j; i, j = i+1, j-1 {
				trows[i], trows[j] = trows[j], trows[i]
			}
		}
		if tr&2 != 0 { // flip X
			for i, r := range trows {
				trows[i] = bits.Reverse16(r) >> (16 - dim)
			}
		}
		if tr&4 != 0 { // rotate 90 degrees clockwise
			rrows := make([]uint16, dim)
			// The first row gets the leftmost column, bottom-to-top.
			// The second row gets the second-leftmost column, and so on.
			for di := 0; di < dim; di++ {
				for si := dim - 1; si >= 0; si-- {
					rrows[di] = (rrows[di] << 1) | ((trows[si] >> (dim - di - 1)) & 0x1)
				}
			}
			trows = rrows
		}

		ti.edges[tr][top] = trows[0]
		ti.edges[tr][bottom] = trows[len(trows)-1]
		for _, row := range trows {
			ti.edges[tr][left] = (ti.edges[tr][left] << 1) | (row >> (dim - 1))
			ti.edges[tr][right] = (ti.edges[tr][right] << 1) | (row & 0x1)
		}
	}
	return ti
}

func dumpRow(r uint16, dim int) string {
	var s string
	for i := dim - 1; i >= 0; i-- {
		if r&(1<<i) != 0 {
			s += "#"
		} else {
			s += "."
		}
	}
	return s
}
