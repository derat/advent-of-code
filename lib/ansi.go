// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import "fmt"

// Color returns the ANSI escape sequence for setting the foreground color to c.
// See https://stackoverflow.com/a/33206814 for available colors.
func Color(c uint8) string {
	return fmt.Sprintf("\033[38;5;%dm", c)
}

const (
	// ClearScreen clears the screen.
	ClearScreen = "\033[2J"
	// MoveHome moves the cursor to the top-left corner of the screen.
	// Writing to stdout will overwrite the existing contents of the screen.
	MoveHome = "\033[H"
)
