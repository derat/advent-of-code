package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pw := lib.InputLinesBytes("2015/11")[0]

	charNum := make(map[byte]byte)
	numChar := make(map[byte]byte)
	var i, max byte
	for ch := byte('a'); ch <= 'z'; ch++ {
		if ch == 'i' || ch == 'l' || ch == 'o' {
			continue
		}
		charNum[ch] = i
		numChar[i] = ch
		max = i
		i++
	}

	for i, ch := range pw {
		n, ok := charNum[ch]
		lib.Assert(ok) // input doesn't seem to have invalid chars
		pw[i] = n
	}

	// Perform a simple increment first.
	increment(pw, max)

	for i := 0; true; i++ {
		if !hasStraight(pw) {
			addStraight(pw, max)
			continue
		}
		if !hasPairs(pw) {
			// TODO: Add pairs instead of just incrementing.
			increment(pw, max)
			continue
		}
		break
	}

	for i, ch := range pw {
		pw[i] = numChar[ch]
	}
	fmt.Println(string(pw))
}

// increment performs a simple increment of pw.
func increment(pw []byte, max byte) {
	for i := len(pw) - 1; i >= 0; i-- {
		if pw[i] < max {
			pw[i]++
			return
		}
		pw[i] = 0         // carry
		lib.Assert(i > 0) // overflow
	}
}

// hasStraight returns true if pw contains a 3-char ascending straight, e.g. "abc".
func hasStraight(pw []byte) bool {
	for i, ch := range pw[:len(pw)-2] {
		if pw[i+1] == ch+1 && pw[i+2] == ch+2 {
			return true
		}
	}
	return false
}

// addStraight inserts an ascending straight in the first possible position.
func addStraight(pw []byte, max byte) {
	for i := len(pw) - 3; i >= 0; i-- {
		ch, ch1, ch2 := pw[i], pw[i+1], pw[i+2]

		// If the two chars to the right are too big, increment this char.
		if ch1 > ch+1 || ch2 > ch+2 {
			ch++
		}

		// Assign the straight and set all the chars to the right to the lowest possible value.
		if ch <= max-2 {
			pw[i], pw[i+1], pw[i+2] = ch, ch+1, ch+2
			for j := i + 3; j < len(pw); j++ {
				pw[j] = 0
			}
			return
		}
	}

	panic(fmt.Sprintf("Failed to add straight to %v", pw))
}

// hasPairs returns true if pw contains at least two different, non-overlapping pairs of chars,
// e.g. "aa" and "zz".
func hasPairs(pw []byte) bool {
	pairs := make(map[byte]struct{})
	for i, ch := range pw[:len(pw)-1] {
		if pw[i+1] == ch {
			pairs[ch] = struct{}{}
			if len(pairs) >= 2 {
				return true
			}
		}
	}
	return false
}
