package lib

import (
	"reflect"
	"testing"
)

func TestExtractMaybe(t *testing.T) {
	for _, tc := range []struct {
		s, re string // input
		n     int    // desired match length
		ok    bool   // desired ok value
		want  []any  // desired extracted values
	}{
		{"abc 23 def", `^(\d)`, 0, false, []any{}},
		{"abc 23 def", `(\d+)`, 2, true, []any{23}},
		{"abc 23 def", `^abc (\d\d) def$`, 10, true, []any{23}},
		{"a -23.1  -45 -9223372036854775808   test  18446744073709551615",
			`^(a|b)\s+(-?\d+\.\d+)\s+(-?\d+)\s+(-?\d+)\s+([a-z]+)\s+(\d+)$`, 62, true,
			[]any{byte('a'), -23.1, -45, int64(-9223372036854775808), "test",
				uint64(18446744073709551615)}},
	} {
		// Create pointers matching the types of the desired data.
		dsts := make([]any, len(tc.want))
		for i, v := range tc.want {
			dsts[i] = reflect.New(reflect.TypeOf(v)).Interface()
		}

		n, ok := ExtractMaybe(tc.s, tc.re, dsts...)
		if n != tc.n || ok != tc.ok {
			t.Errorf("ExtractMaybe(%q, %q, ...) = %v, %v; want %v, %v",
				tc.s, tc.re, n, ok, tc.n, tc.ok)
			continue
		}

		// Dereference the pointers we created earlier so we can compare the concrete values.
		got := make([]any, len(dsts))
		for i, d := range dsts {
			got[i] = reflect.Indirect(reflect.ValueOf(d)).Interface()
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("ExtractMaybe(%q, %q, ...) extracted %v; want %v",
				tc.s, tc.re, got, tc.want)
		}
	}
}
