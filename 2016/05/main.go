package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

const animate = false

func main() {
	const (
		zeros = 5 // number of leading zeros
		pwLen = 8
	)

	pre := lib.InputLines("2016/5")[0]
	want := make([]byte, zeros/2) // full zero bytes

	var pw string
	for i := 0; len(pw) < pwLen; i++ {
		s := pre + strconv.Itoa(i)
		hash := md5.Sum([]byte(s))
		if bytes.HasPrefix(hash[:], want) && (zeros%2 == 0 || hash[zeros/2] < 16) {
			ch := fmt.Sprintf("%x", getVal(hash, zeros))
			pw += ch
			if animate {
				fmt.Print(ch) // print a byte at a time because it looks cooler
			}
		}
	}
	if animate {
		fmt.Println()
	} else {
		fmt.Println(pw)
	}

	// Part 2: Hash specifies position and then char
	// It'd be more efficient to put this in the above loop, but I want cool-looking
	// output and don't want to need to deal with printing to multiple lines.
	pw2 := bytes.Repeat([]byte{'.'}, pwLen)
	rem := pwLen
	if animate {
		fmt.Printf("%s", pw2)
	}
	for i := 0; rem > 0; i++ {
		s := pre + strconv.Itoa(i)
		hash := md5.Sum([]byte(s))
		if bytes.HasPrefix(hash[:], want) && (zeros%2 == 0 || hash[zeros/2] < 16) {
			pos := getVal(hash, zeros)
			if int(pos) >= len(pw2) || pw2[pos] != '.' {
				continue // skip invalid or already-filled position
			}

			pw2[pos] = fmt.Sprintf("%x", getVal(hash, zeros+1))[0]
			rem--
			if animate {
				fmt.Print(strings.Repeat("\b", pwLen)) // clear partial password
				fmt.Printf("%s", pw2)
			}
		}
	}
	if animate {
		fmt.Println()
	} else {
		fmt.Println(string(pw2))
	}
}

// Returns the pos-th 4-bit value from hash.
func getVal(hash [16]byte, pos int) byte {
	return lib.If(pos%2 == 0, hash[pos/2]>>4, hash[pos/2]) & 0xf
}
