package lib

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	want := []string{"a", "b", "c"}
	got := MapKeys(m)
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapKeys(%v) = %v; want %v", m, got, want)
	}
}

func TestMapVals(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	want := []int{1, 2, 3}
	got := MapVals(m)
	sort.Ints(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapVals(%v) = %v; want %v", m, got, want)
	}
}

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
		got := Union(AddSet(nil, strings.Fields(tc.a)...),
			AddSet(nil, strings.Fields(tc.b)...))
		if want := AddSet(nil, strings.Fields(tc.want)...); !reflect.DeepEqual(got, want) {
			t.Errorf("Union(%q, %q) = %q; want %q",
				tc.a, tc.b, strings.Join(MapKeys(got), " "), tc.want)
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
		got := Intersect(AddSet(nil, strings.Fields(tc.a)...), AddSet(nil, strings.Fields(tc.b)...))
		if want := AddSet(nil, strings.Fields(tc.want)...); !reflect.DeepEqual(got, want) {
			t.Errorf("Intersect(%q, %q) = %q; want %q",
				tc.a, tc.b, strings.Join(MapKeys(got), " "), tc.want)
		}
	}
}
