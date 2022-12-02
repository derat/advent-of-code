// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

// MapKeys returns keys from the provided map in an arbitrary order.
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapVals returns values from the provided map in an arbitrary order.
func MapVals[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// MapKeysWithVal returns keys from the provided map that have want as their value.
func MapKeysWithVal[K, V comparable](m map[K]V, want V) []K {
	var keys []K
	for k, v := range m {
		if v == want {
			keys = append(keys, k)
		}
	}
	return keys
}

// MapSomeKey returns an arbitrary key from the supplied map.
func MapSomeKey[K comparable, V any](m map[K]V) K {
	for k := range m {
		return k
	}
	panic("Can't get key from empty map")
}

// MapHasKey returns true if map m contains key k.
func MapHasKey[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

// MapHasValue returns true if map m contains the specified value.
func MapHasValue[K, V comparable](m map[K]V, want V) bool {
	for _, v := range m {
		if v == want {
			return true
		}
	}
	return false
}

// InvertMap returns a new map in which m's values map to its keys.
func InvertMap[K, V comparable](m map[K]V) map[V]K {
	r := make(map[V]K, len(m))
	for k, v := range m {
		r[v] = k
	}
	return r
}

// Set returns a map-to-empty-struct containing keys from m, a map.
func Set[K comparable, V any](m map[K]V) map[K]struct{} {
	s := make(map[K]struct{}, len(m))
	for k := range m {
		s[k] = struct{}{}
	}
	return s
}

// AddSet adds keys to the supplied (possibly-nil) map to struct{}.
// The set is returned (and should be used thereafter).
func AddSet[K comparable](m map[K]struct{}, keys ...K) map[K]struct{} {
	if m == nil {
		m = make(map[K]struct{}, len(keys))
	}
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

// Union returns a new map with the union of keys from a and b.
// If a key is present in both maps, the value from a will be used.
func Union[K comparable, V any](a, b map[K]V) map[K]V {
	c := make(map[K]V, len(a)+len(b))
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

// Intersect returns a new map with the intersection of keys from a and b.
// Values from a will be used.
func Intersect[K comparable, V any](a, b map[K]V) map[K]V {
	c := make(map[K]V, Max(len(a), len(b)))
	for k, v := range a {
		if _, ok := b[k]; ok {
			c[k] = v
		}
	}
	return c
}
