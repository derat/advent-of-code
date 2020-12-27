package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var rows [][]byte
	for sc.Scan() {
		if len(sc.Bytes()) == 0 {
			continue
		}
		rows = append(rows, append([]byte(nil), sc.Bytes()...)) // !@#!@#@!
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	occ := func(rows [][]byte, row, col int) int {
		if row < 0 || col < 0 || row >= len(rows) || col >= len(rows[row]) {
			return 0
		}
		if rows[row][col] != '#' {
			return 0
		}
		return 1
	}

	look := func(rows [][]byte, row, col, dr, dc int) int {
		for r, c := row+dr, col+dc; r >= 0 && r < len(rows) && c >= 0 && c < len(rows[r]); r, c = r+dr, c+dc {
			switch rows[r][c] {
			case '#':
				return 1
			case 'L':
				return 0
			}
		}
		return 0
	}

	// Part 1
	println(run(rows, func(rows [][]byte, row, col int) int {
		return occ(rows, row-1, col-1) + occ(rows, row-1, col) + occ(rows, row-1, col+1) +
			occ(rows, row, col-1) + occ(rows, row, col+1) +
			occ(rows, row+1, col-1) + occ(rows, row+1, col) + occ(rows, row+1, col+1)
	}, 4))

	// Part 2
	println(run(rows, func(rows [][]byte, row, col int) int {
		return look(rows, row, col, -1, -1) + look(rows, row, col, -1, 0) + look(rows, row, col, -1, 1) +
			look(rows, row, col, 0, -1) + look(rows, row, col, 0, 1) +
			look(rows, row, col, 1, -1) + look(rows, row, col, 1, 0) + look(rows, row, col, 1, 1)
	}, 5))
}

type countFunc func(rows [][]byte, row, col int) int

func run(rows [][]byte, f countFunc, flip int) int {
	for {
		newRows := make([][]byte, len(rows))
		for i := range rows {
			newRows[i] = bytes.Repeat([]byte{0}, len(rows[i]))
			for j := range rows[i] {
				ch := rows[i][j]
				switch ch {
				case '.': // floor
				case 'L': // empty
					if n := f(rows, i, j); n == 0 {
						ch = '#'
					}
				case '#': // occupied
					if n := f(rows, i, j); n >= flip {
						ch = 'L'
					}
				default:
					log.Fatalf("invalid state %q", ch)
				}
				newRows[i][j] = ch
			}
		}

		//fmt.Println(dump(rows) + "\n")

		if bytes.Equal(bytes.Join(newRows, nil), bytes.Join(rows, nil)) {
			cnt := 0
			for _, row := range rows {
				for _, seat := range row {
					if seat == '#' {
						cnt++
					}
				}
			}
			return cnt
		}

		rows = newRows
	}
}

func dump(rows [][]byte) string {
	return string(bytes.Join(rows, []byte{'\n'}))
}
