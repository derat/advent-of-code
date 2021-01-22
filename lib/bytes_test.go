package lib

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestCopyBytesRegion(t *testing.T) {
	for _, tc := range []struct {
		full   string
		r0, c0 int
		r1, c1 int
		want   string
	}{
		{"123\n456\n789", 0, 0, 2, 2, "123\n456\n789"},
		{"123\n456\n789", 0, 0, 0, 0, "1"},
		{"123\n456\n789", -1, -1, 0, 0, "1"},
		{"123\n456\n789", 2, 2, 3, 3, "9"},
		{"123\n456\n789", 0, 0, 0, 1, "12"},
		{"123\n456\n789", 0, 0, 1, 0, "1\n4"},
		{"123\n456\n789", 0, 0, 1, 1, "12\n45"},
		{"123\n456\n789", 1, 0, 2, 2, "456\n789"},
	} {
		out := CopyBytesRegion(ByteLines(tc.full), tc.r0, tc.c0, tc.r1, tc.c1)
		if got := string(bytes.Join(out, []byte("\n"))); got != tc.want {
			t.Errorf("CopyBytesRegion(%q, %d, %d, %d, %d) = %q; want %q",
				tc.full, tc.r0, tc.c0, tc.r1, tc.c1, got, tc.want)
		}
	}
}

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
