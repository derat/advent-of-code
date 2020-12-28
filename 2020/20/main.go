package main

import (
	"fmt"
	"log"
	"math"
	"math/bits"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var dim int // width and height of each tile (tiles are square)
	tiles := make(map[int]tile)

	for _, pg := range lib.ReadParagraphs() {
		var id int
		lib.Parse(pg[0], `^Tile (\d+):$`, &id)

		var rows []uint16
		for _, ln := range pg[1:] {
			if dim == 0 {
				lib.AssertLessEq(len(ln), 16)
				dim = len(ln)
			}
			lib.AssertEq(len(ln), dim)

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
		lib.AssertEq(len(rows), dim)
		tiles[id] = newTile(rows)
	}

	// Map from edge to tile IDs containing it.
	edgeIds := make(map[uint16]map[int]struct{})
	for id, ti := range tiles {
		for _, edges := range ti.edges {
			for _, edge := range edges {
				m := edgeIds[edge]
				if m == nil {
					m = make(map[int]struct{})
				}
				m[id] = struct{}{}
				edgeIds[edge] = m
			}
		}
	}

	sideLen := int(math.Sqrt(float64(len(tiles))))
	type tileInfo struct{ id, tr int }
	sol := make([][]tileInfo, sideLen)
	for i := range sol {
		sol[i] = make([]tileInfo, sideLen)
	}

	remain := make(map[int]struct{}, len(tiles))
	for id := range tiles {
		remain[id] = struct{}{}
	}

	var solve func(r, c int) bool
	solve = func(r, c int) bool {
		//fmt.Printf("Solving for %v, %v\n", r, c)
		for id := range remain {
			delete(remain, id) // use the tile
			for tr, edges := range tiles[id].edges {
				// The first/last row need unique top/bottom edges.
				if r == 0 && len(edgeIds[edges[top]]) > 1 {
					continue
				} else if r == sideLen-1 && len(edgeIds[edges[bottom]]) > 1 {
					continue
				}
				// The first/last column need unique left/right edges.
				if c == 0 && len(edgeIds[edges[left]]) > 1 {
					continue
				} else if c == sideLen-1 && len(edgeIds[edges[right]]) > 1 {
					continue
				}
				// If we've already filled in a tile above or to the left, make sure
				// that our top/left edge matches its bottom/right edge.
				if r > 0 {
					if t := sol[r-1][c]; edges[top] != tiles[t.id].edges[t.tr][bottom] {
						continue
					}
				}
				if c > 0 {
					if l := sol[r][c-1]; edges[left] != tiles[l.id].edges[l.tr][right] {
						continue
					}
				}

				// Fill in the tile.
				sol[r][c] = tileInfo{id, tr}
				var nr, nc int
				if c < sideLen-1 {
					nr, nc = r, c+1
				} else if r < sideLen-1 {
					nr, nc = r+1, 0
				} else {
					return true // all done!
				}
				if done := solve(nr, nc); done {
					return true
				}
			}
			remain[id] = struct{}{} // didn't work; put the tile back
		}
		return false
	}

	if !solve(0, 0) {
		log.Fatal("Didn't find solution")
	}

	// Print solution:
	/*
		for srow := range sol {
			for trow := 0; trow < dim; trow++ { // skip top and bottom rows
				for scol := range sol[srow] {
					info := sol[srow][scol]
					row := tiles[info.id].rows[info.tr][trow]
					fmt.Print(dumpRow(row, dim), " ")
				}
				fmt.Println()
			}
			fmt.Println()
		}
	*/

	// Part 1:
	var prod int64 = int64(sol[0][0].id) * int64(sol[0][sideLen-1].id) *
		int64(sol[sideLen-1][0].id) * int64(sol[sideLen-1][sideLen-1].id)
	fmt.Println(prod)

	// Part 2:
	tdim := dim - 2
	idim := sideLen * tdim
	img := make([][]byte, idim)
	for i := range img {
		img[i] = make([]byte, idim)
	}
	for srow := range sol {
		for scol := range sol[srow] {
			info := sol[srow][scol]
			rows := tiles[info.id].rows[info.tr]
			for trow := 0; trow < tdim; trow++ {
				row := rows[trow+1] // skip first row
				for tcol := 0; tcol < tdim; tcol++ {
					var val byte = '.'
					if row&(1<<(tdim-tcol)) != 0 { // skip first col
						val = '#'
					}
					irow := srow*tdim + trow
					icol := scol*tdim + tcol
					img[irow][icol] = val
				}
			}
		}
	}

	monster := [][]byte{
		[]byte("                  # "),
		[]byte("#    ##    ##    ###"),
		[]byte(" #  #  #  #  #  #   "),
	}
	mwidth := len(monster[0])
	mheight := len(monster)

	for tr := 0; tr < 8; tr++ {
		timg := make([][]byte, idim)
		for i, r := range img {
			timg[i] = append([]byte{}, r...)
		}

		if tr&1 == 1 { // flip Y
			for i, j := 0, idim-1; i < j; i, j = i+1, j-1 {
				timg[i], timg[j] = timg[j], timg[i]
			}
		}
		if tr&2 != 0 { // flip X
			for _, r := range timg {
				for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
					r[i], r[j] = r[j], r[i]
				}
			}
		}
		if tr&4 != 0 { // rotate 90 degrees clockwise
			rimg := make([][]byte, idim)
			for i := 0; i < idim; i++ { // iterate over dest rows
				rimg[i] = make([]byte, idim)
				for j := 0; j < idim; j++ { // iterate over dest cols
					rimg[i][j] = timg[j][i]
				}
			}
			timg = rimg
		}

		/*
			// Print transformed image.
			for _, r := range timg {
				fmt.Println(string(r))
			}
			fmt.Println()
		*/

		nmon := 0
		for row := 0; row < idim-mheight; row++ {
			for col := 0; col < idim-mwidth; col++ {
				match := true
			Check:
				for i := 0; i < mheight; i++ {
					for j := 0; j < mwidth; j++ {
						if monster[i][j] == '#' && timg[row+i][col+j] != '#' {
							match = false
							break Check
						}
					}
				}
				if match {
					nmon++
					for i := 0; i < mheight; i++ {
						for j := 0; j < mwidth; j++ {
							if monster[i][j] == '#' {
								timg[row+i][col+j] = '.'
							}
						}
					}
				}
			}
		}
		if nmon > 0 {
			roughness := 0
			for _, r := range timg {
				for _, c := range r {
					if c == '#' {
						roughness++
					}
				}
			}
			fmt.Println(roughness)
			break
		}
	}
}

const (
	top    = 0
	bottom = 1
	left   = 2
	right  = 3
)

type tile struct {
	edges [8][4]uint16 // outer is transformation, inner is edge (LSB right/bottom)
	rows  [8][]uint16  // outer is transformation, inner is rows (LSB right)
}

func newTile(rows []uint16) tile {
	dim := len(rows)
	var ti tile

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
		ti.rows[tr] = trows
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
