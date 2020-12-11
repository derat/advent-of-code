package main

import (
	"bufio"
	"bytes"
	"fmt"
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

	const flip = 4 // number of adj occ seats to switch to empty

	occ := func(row, col int) int {
		if row < 0 || col < 0 || row >= len(rows) || col >= len(rows[row]) {
			return 0
		}
		if rows[row][col] != '#' {
			return 0
		}
		return 1
	}

	for {
		newRows := make([][]byte, len(rows))
		for i := range rows {
			newRows[i] = bytes.Repeat([]byte{0}, len(rows[i]))
			for j := range rows[i] {
				ch := rows[i][j]
				if ch != '.' {
					adj := occ(i-1, j-1) + occ(i-1, j) + occ(i-1, j+1) +
						occ(i, j-1) + occ(i, j+1) +
						occ(i+1, j-1) + occ(i+1, j) + occ(i+1, j+1)
					if ch == 'L' {
						if adj == 0 {
							ch = '#'
						}
					} else if ch == '#' {
						if adj >= flip {
							ch = 'L'
						}
					} else {
						log.Fatalf("invalid state %q", ch)
					}
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
			fmt.Println(cnt)
			break
		}

		rows = newRows
	}
}

func dump(rows [][]byte) string {
	return string(bytes.Join(rows, []byte{'\n'}))
}
