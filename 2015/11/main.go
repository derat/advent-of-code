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

	// Remap to [0, ...].
	for i, ch := range pw {
		n, ok := charNum[ch]
		lib.Assert(ok) // input doesn't seem to have invalid chars
		pw[i] = n
	}

	// Remap back to the original chars.
	remap := func(pw []byte) string {
		b := make([]byte, len(pw))
		for i, ch := range pw {
			b[i] = numChar[ch]
		}
		return string(b)
	}

	// Part 1
	findNext(pw, max)
	fmt.Println(remap(pw))

	// Part 2
	findNext(pw, max)
	fmt.Println(remap(pw))
}

func findNext(pw []byte, max byte) {
	// Perform a simple increment first.
	increment(pw, max)

	// Now loop until we get both a straight and pairs.
	for {
		if !hasStraight(pw) {
			addStraight(pw, max)
			continue
		}
		if !hasPairs(pw) {
			for !hasPairs(pw) {
				// This is stupid, but I had enough trouble getting addStraight()
				// to work and don't want to do the same thing with an addPairs()
				// function.
				increment(pw, max)
			}
			continue
		}
		return
	}
}

// increment performs a simple increment of pw.
func increment(pw []byte, max byte) {
	for i := len(pw) - 1; i >= 0; i-- {
		if pw[i] < max {
			pw[i]++
			return
		}
		pw[i] = 0 // carry
	}
	panic("Overflow on increment")
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
	for {
		// Setting the first or second char in the loop could've produced a straight
		// further to the left.
		if hasStraight(pw) {
			return
		}

		// Consider sequences of three characters starting at the right.
		for i := len(pw) - 3; i >= 0; i-- {
			ch, ch1, ch2 := pw[i], pw[i+1], pw[i+2]

			// Try to set the third char to make the straight.
			if ch1 == ch+1 && ch2 < ch+2 && ch < max-1 {
				pw[i+2] = ch + 2
				for j := i + 3; j < len(pw); j++ {
					pw[j] = 0
				}
				return
			}

			// Set the second char.
			if ch1 < ch+1 && ch < max {
				pw[i+1] = ch + 1
				for j := i + 2; j < len(pw); j++ {
					pw[j] = 0
				}
				break
			}

			// Increment the first char.
			if ch1 >= ch+1 && ch < max {
				pw[i] = ch + 1
				for j := i + 1; j < len(pw); j++ {
					pw[j] = 0
				}
				break
			}

			// Otherwise, we'll move on to the left.
		}
	}
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
