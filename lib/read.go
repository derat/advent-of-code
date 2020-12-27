package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

// ReadIntsN is similar to ReadInts() but panics unless exactly n ints are read.
func ReadIntsN(n int) []int {
	vals := ReadInts()
	if len(vals) != n {
		panic(fmt.Sprintf("Got %v int(s); want %v", len(vals), n))
	}
	return vals
}
