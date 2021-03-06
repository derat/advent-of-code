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

// GenericSet returns a map[interface{}]struct{} containing keys from m, a map.
func GenericSet(m interface{}) map[interface{}]struct{} {
	mv := reflect.ValueOf(m)
	s := make(map[interface{}]struct{}, mv.Len())
	it := reflect.ValueOf(m).MapRange()
	for it.Next() {
		s[it.Key().Interface()] = struct{}{}
	}
	return s
}

// AddSet adds the supplied as keys in the supplied (possibly-nil) map to struct{}.
// The set is returned (and should be used thereafter).
func AddSet(s interface{}, vals ...interface{}) interface{} {
	sv := reflect.ValueOf(s)
	if sv.IsNil() {
		sv = reflect.MakeMapWithSize(sv.Type(), len(vals))
	}
	for _, v := range vals {
		sv.SetMapIndex(reflect.ValueOf(v), reflect.ValueOf(struct{}{}))
	}
	return sv.Interface()
}

// AddStringSet is a wrapper around AddSet for string keys.
func AddStringSet(s map[string]struct{}, vals ...string) map[string]struct{} {
	ivals := make([]interface{}, len(vals))
	for i := range vals {
		ivals[i] = vals[i]
	}
	return AddSet(s, ivals...).(map[string]struct{})
}

// Union returns a new map with the union of keys from a and b.
func Union(a, b interface{}) interface{} {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	AssertEq(av.Type(), bv.Type())
	cv := reflect.MakeMapWithSize(av.Type(), Max(av.Len(), bv.Len()))

	ai := av.MapRange()
	for ai.Next() {
		cv.SetMapIndex(ai.Key(), ai.Value())
	}
	bi := bv.MapRange()
	for bi.Next() {
		cv.SetMapIndex(bi.Key(), bi.Value())
	}

	return cv.Interface()
}

// Intersect returns a new map with the intersection of keys from a and b.
func Intersect(a, b interface{}) interface{} {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	AssertEq(av.Type(), bv.Type())
	cv := reflect.MakeMap(av.Type())

	var it *reflect.MapIter
	var ov reflect.Value
	if av.Len() < bv.Len() {
		it = av.MapRange()
		ov = bv
	} else {
		it = bv.MapRange()
		ov = av
	}
	for it.Next() {
		if ov.MapIndex(it.Key()).IsValid() {
			cv.SetMapIndex(it.Key(), it.Value())
		}
	}

	return cv.Interface()
}
