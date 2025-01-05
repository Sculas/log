// Package sortedmap provides a sorted map implementation.
// It is kindly taken from https://github.com/gobs/sortedmap under the MIT license.
package sortedmap

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
)

// KeyValuePair describes an entry in SortedMap
type KeyValuePair[K constraints.Ordered, T any] struct {
	Key   K
	Value T
}

// String implements the Stringer interface for KeyValuePair
func (e KeyValuePair[K, T]) String() string {
	if key, ok := any(e.Key).(string); ok {
		return fmt.Sprintf("%q: %v", key, e.Value)
	}

	return fmt.Sprintf("%v: %v", e.Key, e.Value)
}

// SortedMap is a structure that can sort a map[string]type by key
type SortedMap[K constraints.Ordered, T any] []KeyValuePair[K, T]

func (s SortedMap[K, T]) Len() int           { return len(s) }
func (s SortedMap[K, T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortedMap[K, T]) Less(i, j int) bool { return s[i].Key < s[j].Key }

// Sort sorts a SortedMap
func (s SortedMap[K, T]) Sort() { sort.Sort(s) }

// Add adds an entry to a SortedMap (this requires re-sorting the SortedMap when ready to display).
// Note that a SortedMap is internally a slice, so you need to do something like:
//
//	s := NewSortedMap()
//	s = s.Add(key1, value1)
//	s = s.Add(key2, value2)
func (s SortedMap[K, T]) Add(key K, value T) SortedMap[K, T] {
	return append(s, KeyValuePair[K, T]{key, value})
}

// Keys returns the list of keys for the entries in this SortedMap
func (s SortedMap[K, T]) Keys() []K {
	keys := make([]K, 0, len(s))
	for _, kv := range s {
		keys = append(keys, kv.Key)
	}

	return keys
}

// Values returns the list of values for the entries in this SortedMap
func (s SortedMap[K, T]) Values() []T {
	values := make([]T, 0, len(s))
	for _, kv := range s {
		values = append(values, kv.Value)
	}

	return values
}

// AsSortedMap return a SortedMap from a map[string]type.
// Note that this will panic if the input object is not a map
func AsSortedMap[K constraints.Ordered, T any](m map[K]T) SortedMap[K, T] {
	s := make(SortedMap[K, T], 0, len(m))
	for k, v := range m {
		s = append(s, KeyValuePair[K, T]{k, v})
	}
	s.Sort()

	return s
}
