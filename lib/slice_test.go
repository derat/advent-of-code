package lib

import "testing"

func TestRotateSlice(t *testing.T) {
	for _, tc := range []struct {
		in   string
		amt  int
		want string
	}{
		{"abcd", 0, "abcd"},
		{"abcd", 1, "dabc"},
		{"abcd", 2, "cdab"},
		{"abcd", 3, "bcda"},
		{"abcd", 4, "abcd"},
		{"abcd", 5, "dabc"},
		{"abcd", -1, "bcda"},
		{"abcd", -2, "cdab"},
		{"abcd", -3, "dabc"},
		{"abcd", -4, "abcd"},
		{"abcd", -5, "bcda"},
		{"", 0, ""},
		{"", 1, ""},
		{"", -1, ""},
		{"a", 0, "a"},
		{"a", 1, "a"},
		{"a", -1, "a"},
	} {
		b := []byte(tc.in)
		if RotateSlice(b, tc.amt); string(b) != tc.want {
			t.Errorf("RotateBy(%q, %d) = %q; want %q", tc.in, tc.amt, b, tc.want)
		}
	}
}
