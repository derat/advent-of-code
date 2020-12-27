package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Line reads and returns a single line of input from stdin.
// Leading and trailing whitespace is stripped.
// If valid is non-empty, only the supplied bytes are permitted.
func Line(valid ...byte) []byte {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, os.Stdin)
	if err != nil {
		panic(fmt.Sprint("Failed to read input: ", err))
	}
	b := bytes.TrimSpace(buf.Bytes())
	if len(valid) > 0 {
		for i, ch := range b {
			if bytes.IndexByte(valid, ch) == -1 {
				panic(fmt.Sprintf("Invalid byte %v (%q) at position %d", ch, ch, i))
			}
		}
	}
	return b
}
