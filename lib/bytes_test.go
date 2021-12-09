package lib

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestByteGrid_CopyRect(t *testing.T) {
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
		out := NewByteGridString(tc.full).CopyRect(tc.r0, tc.c0, tc.r1, tc.c1)
		if got := string(bytes.Join(out, []byte("\n"))); got != tc.want {
			t.Errorf("%q CopyRect(%d, %d, %d, %d) = %q; want %q",
				tc.full, tc.r0, tc.c0, tc.r1, tc.c1, got, tc.want)
		}
	}
}

func TestByteGrid_SetRect(t *testing.T) {
	const (
		sr = 3
		sc = 4
	)
	for _, tc := range []struct {
		r0, c0 int
		r1, c1 int
		want   string
	}{
		{0, 0, 0, 0, "x...\n....\n...."},
		{2, 3, 2, 3, "....\n....\n...x"},
		{0, 0, 2, 3, "xxxx\nxxxx\nxxxx"},
		{0, 0, 1, 2, "xxx.\nxxx.\n...."},
		{2, 3, 1, 1, "....\n.xxx\n.xxx"},   // swap
		{-1, -1, 3, 4, "xxxx\nxxxx\nxxxx"}, // clamp
	} {
		b := NewByteGrid(sr, sc, '.')
		b.SetRect(tc.r0, tc.c0, tc.r1, tc.c1, 'x')
		if got := b.Dump(); got != tc.want {
			t.Errorf("SetRect(%d, %d, %d, %d, 'x') = %q; want %q",
				tc.r0, tc.c0, tc.r1, tc.c1, got, tc.want)
		}
	}
}

func TestByteGrid_IterLine(t *testing.T) {
	const (
		sr = 3
		sc = 4
	)
	for _, tc := range []struct {
		r0, c0 int
		r1, c1 int
		want   string
	}{
		{0, 0, 0, 0, "x...\n....\n...."},
		{2, 3, 2, 3, "....\n....\n...x"},
		{1, 1, 1, 2, "....\n.xx.\n...."},   // horizontal
		{0, 2, 1, 2, "..x.\n..x.\n...."},   // vertical
		{0, 0, 1, 1, "x...\n.x..\n...."},   // diagonal
		{1, 3, 2, 2, "....\n...x\n..x."},   // diagonal
		{2, 2, 0, 0, "x...\n.x..\n..x."},   // diagonal
		{0, 0, 2, 3, "x...\n.xx.\n...x"},   // diagonal, not 1:1
		{0, 3, 2, 0, "...x\n.xx.\nx..."},   // diagonal, not 1:1
		{-1, -1, 4, 4, "x...\n.x..\n..x."}, // out of bounds
	} {
		b := NewByteGrid(sr, sc, '.')
		b.IterLine(tc.r0, tc.c0, tc.r1, tc.c1, func(r, c int) { b[r][c] = 'x' })
		if got := b.Dump(); got != tc.want {
			t.Errorf("IterLine(%d, %d, %d, %d, ...) = %q; want %q",
				tc.r0, tc.c0, tc.r1, tc.c1, got, tc.want)
		}
	}
}

func TestByteGrid_FlipHoriz(t *testing.T) {
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
		got := NewByteGridString(strings.TrimLeft(tc.input, "\n")).FlipHoriz()
		if want := NewByteGridString(strings.TrimLeft(tc.want, "\n")); !reflect.DeepEqual(got, want) {
			t.Errorf("FlipHoriz produced %q; want %q", got, want)
		}
	}
}

func TestByteGrid_FlipVert(t *testing.T) {
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
		got := NewByteGridString(strings.TrimLeft(tc.input, "\n")).FlipVert()
		if want := NewByteGridString(strings.TrimLeft(tc.want, "\n")); !reflect.DeepEqual(got, want) {
			t.Errorf("FlipVert produced %q; want %q", got, want)
		}
	}
}

func TestByteGrid_RotateCW(t *testing.T) {
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
		got := NewByteGridString(strings.TrimLeft(tc.input, "\n")).RotateCW()
		if want := NewByteGridString(strings.TrimLeft(tc.want, "\n")); !reflect.DeepEqual(got, want) {
			t.Errorf("RotateCW produced %q; want %q", got, want)
		}
	}
}
