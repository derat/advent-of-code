package lib

import (
	"reflect"
	"strings"
	"testing"
)

func TestUnion(t *testing.T) {
	for _, tc := range []struct {
		a, b, want string
	}{
		{"", "", ""},
		{"a", "", "a"},
		{"", "a", "a"},
		{"a", "a", "a"},
		{"a e b c", "x a e", "a b c e x"},
		{"x a e", "a e b c", "a b c e x"},
	} {
		got := Union(AddStringSet(nil, strings.Fields(tc.a)...),
			AddStringSet(nil, strings.Fields(tc.b)...))
		if want := AddStringSet(nil, strings.Fields(tc.want)...); !reflect.DeepEqual(got, want) {
			t.Errorf("Union(%q, %q) = %q; want %q",
				tc.a, tc.b, strings.Join(MapStringKeys(got), " "), tc.want)
		}
	}
}

func TestIntersect(t *testing.T) {
	for _, tc := range []struct {
		a, b, want string
	}{
		{"", "", ""},
		{"a", "", ""},
		{"", "a", ""},
		{"a", "a", "a"},
		{"a e b c", "x a e", "a e"},
		{"x a e", "a e b c", "a e"},
	} {
		got := Intersect(AddStringSet(nil, strings.Fields(tc.a)...),
			AddStringSet(nil, strings.Fields(tc.b)...))
		if want := AddStringSet(nil, strings.Fields(tc.want)...); !reflect.DeepEqual(got, want) {
			t.Errorf("Intersect(%q, %q) = %q; want %q",
				tc.a, tc.b, strings.Join(MapStringKeys(got), " "), tc.want)
		}
	}
}
