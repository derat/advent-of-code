package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadAll reads and returns all data from stdin.
func ReadAll() []byte {
	var b bytes.Buffer
	_, err := io.Copy(&b, os.Stdin)
	if err != nil {
		panic(fmt.Sprint("Failed to read input: ", err))
	}
	return b.Bytes()
}

// ReadInts reads positive integers from stdin, ignoring non-digits.
// See ExtractInts.
func ReadInts() []int {
	return ExtractInts(string(ReadAll()))
}

// ReadIntsN is similar to ReadInts but panics unless exactly n ints are read.
func ReadIntsN(n int) []int {
	vals := ReadInts()
	if len(vals) != n {
		panic(fmt.Sprintf("Got %v int(s); want %v", len(vals), n))
	}
	return vals
}

// ReadLines reads and returns newline-separated lines of input from stdin.
func ReadLines() []string {
	var lines []string
	for _, ln := range strings.Split(string(ReadAll()), "\n") {
		if len(ln) > 0 {
			lines = append(lines, ln)
		}
	}
	return lines
}

// ReadLinesBytes reads and returns newline-separated lines of input from stdin.
// If valid is non-empty, panics if any unlisted bytes are encountered.
func ReadLinesBytes(valid ...byte) [][]byte {
	var lines [][]byte
	for i, ln := range bytes.Split(ReadAll(), []byte{'\n'}) {
		if len(ln) == 0 {
			continue
		}
		if len(valid) > 0 {
			for j, ch := range ln {
				if bytes.IndexByte(valid, ch) == -1 {
					panic(fmt.Sprintf("Invalid byte %v (%q) at position %d of line %d", ch, ch, j, i))
				}
			}
		}
		lines = append(lines, ln)
	}
	return lines
}

// ReadLineBytes is similar to ReadLinesBytes but panics unless exactly one line is read.
func ReadLineBytes(valid ...byte) []byte {
	lines := ReadLinesBytes(valid...)
	if len(lines) != 1 {
		panic(fmt.Sprintf("Got %v lines; want 1", len(lines)))
	}
	return lines[0]
}
