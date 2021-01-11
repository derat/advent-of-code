package lib

import (
	"reflect"
	"strings"
	"testing"
)

func TestFlipBytesX(t *testing.T) {
	for _, tc := range []struct {
		input, want string
	}{
		{
			`
##.
..#`,
			`
.##
#..`,
		},
		{
			`
##.#
#..#
##.#
..#.`,
			`
#.##
#..#
#.##
.#..`,
		},
	} {
		got := ByteLines(strings.TrimLeft(tc.input, "\n"))
		FlipBytesX(got)
		if want := ByteLines(strings.TrimLeft(tc.want, "\n")); !reflect.DeepEqual(got, want) {
			t.Errorf("FlipBytesX produced %q; want %q", got, want)
		}
	}
}

func TestFlipBytesY(t *testing.T) {
	for _, tc := range []struct {
		input, want string
	}{
		{
			`
##.
..#`,
			`
..#
##.`,
		},
		{
			`
##.#
#..#
##.#
..#.`,
			`
..#.
##.#
#..#
##.#`,
		},
	} {
		got := ByteLines(strings.TrimLeft(tc.input, "\n"))
		FlipBytesY(got)
		if want := ByteLines(strings.TrimLeft(tc.want, "\n")); !reflect.DeepEqual(got, want) {
			t.Errorf("FlipBytesY produced %q; want %q", got, want)
		}
	}
}

func TestRotateBytes(t *testing.T) {
	for _, tc := range []struct {
		input, want string
	}{
		{
			`
##.
..#`,
			`
.#
.#
#.`,
		},
		{
			`
##.#
#..#
##.#
..#.`,
			`
.###
.#.#
#...
.###`,
		},
	} {
		got := RotateBytes(ByteLines(strings.TrimLeft(tc.input, "\n")))
		if want := ByteLines(strings.TrimLeft(tc.want, "\n")); !reflect.DeepEqual(got, want) {
			t.Errorf("RotateBytes produced %q; want %q", got, want)
		}
	}
}
