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

// MapIntKeys returns integer keys from the provided map.
func MapIntKeys(m interface{}) []int {
	mv := reflect.ValueOf(m)
	keys := make([]int, 0, mv.Len())
	for _, kv := range mv.MapKeys() {
		keys = append(keys, int(kv.Int()))
	}
	return keys
}

// MapUint64Keys returns uint64 keys from the provided map.
func MapUint64Keys(m interface{}) []uint64 {
	mv := reflect.ValueOf(m)
	keys := make([]uint64, 0, mv.Len())
	for _, kv := range mv.MapKeys() {
		keys = append(keys, kv.Uint())
	}
	return keys
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

// MapHasKey returns true if map m contains key k.
func MapHasKey(m, k interface{}) bool {
	return reflect.ValueOf(m).MapIndex(reflect.ValueOf(k)).IsValid()
}

// MapHasValue returns true if map m contains value v.
func MapHasValue(m, v interface{}) bool {
	it := reflect.ValueOf(m).MapRange()
	for it.Next() {
		if it.Value().Interface() == v {
			return true
		}
	}
	return false
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
