package lib

import "reflect"

// MapIntVals returns integer values from the provided map.
func MapIntVals(m interface{}) []int {
	mv := reflect.ValueOf(m)
	vals := make([]int, 0, mv.Len())
	it := mv.MapRange()
	for it.Next() {
		vals = append(vals, int(it.Value().Int()))
	}
	return vals
}

// MapStringKeys returns string keys from the provided map.
func MapStringKeys(m interface{}) []string {
	mv := reflect.ValueOf(m)
	keys := make([]string, 0, mv.Len())
	for _, kv := range mv.MapKeys() {
		keys = append(keys, kv.String())
	}
	return keys
}

// MapStringKeysWithVal returns string keys from map m with values equal to v.
func MapStringKeysWithVal(m, v interface{}) []string {
	mv := reflect.ValueOf(m)
	var keys []string
	it := mv.MapRange()
	for it.Next() {
		if it.Value().Interface() == v {
			keys = append(keys, it.Key().String())
		}
	}
	return keys
}

// MapSomeKey returns an arbitrary key from the supplied map.
func MapSomeKey(m interface{}) interface{} {
	mv := reflect.ValueOf(m)
	for _, kv := range mv.MapKeys() {
		return kv.Interface()
	}
	panic("Can't get key from empty map")
}

// AddSet adds the supplied values to the supplied string set.
// The set is returned (and should be used thereafter).
func AddSet(s map[string]struct{}, vals ...string) map[string]struct{} {
	if s == nil {
		s = make(map[string]struct{}, len(vals))
	}
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

// SetHas returns true if s contains val.
func SetHas(s map[string]struct{}, val string) bool {
	_, ok := s[val]
	return ok
}

// SetHas64 returns true if s contains val.
func SetHas64(s map[uint64]struct{}, val uint64) bool {
	_, ok := s[val]
	return ok
}
