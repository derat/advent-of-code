// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

// Dir represents a cardinal direction.
type Dir int

const (
	Up Dir = iota
	Left
	Down
	Right

	// Alternate names.
	North = Up
	West  = Left
	South = Down
	East  = Right
)

// Left returns the resulting direction after turning to the left.
func (d Dir) Left() Dir {
	return Dir((int(d) + 1) % 4)
}

// Right returns the resulting direction after turning to the right.
func (d Dir) Right() Dir {
	return Dir((int(d) + 3) % 4)
}

// Reverse returns the resulting direction after turning 180 degrees.
func (d Dir) Reverse() Dir {
	return Dir((int(d) + 2) % 4)
}

// DR returns the change in row after a single step.
func (d Dir) DR() int {
	switch d {
	case Up:
		return -1
	case Down:
		return 1
	case Left, Right:
		return 0
	default:
		panic("Invalid dir")
	}
}

// DC returns the change in column after a single step.
func (d Dir) DC() int {
	switch d {
	case Up, Down:
		return 0
	case Left:
		return -1
	case Right:
		return 1
	default:
		panic("Invalid dir")
	}
}
