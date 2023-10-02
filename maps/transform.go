// Copyright 2023 the go-functional authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maps

import "errors"

var (
	ErrDuplicateKeys = errors.New("input slice has duplicate keys according to key function")
)

type KeyFunc[K comparable, T any] func(v T) K

// ToMap converts a slice that is map indexed into a map according to the provided key
// function. Duplicates are not allowed, so an error is returned in the case of a duplicate.
func ToMap[K comparable, T any](in []T, keyFunc KeyFunc[K, T]) (map[K]T, error) {
	rtn := make(map[K]T, len(in))
	for _, v := range in {
		k := keyFunc(v)
		if _, dupe := rtn[k]; dupe {
			return nil, ErrDuplicateKeys
		}
		rtn[k] = v
	}
	return rtn, nil
}

// ToSlice takes in a map and returns a slice of the values from the map.
// Order of the slice is not guaranteed in any way.
func ToSlice[K comparable, T any](in map[K]T) []T {
	rtn := make([]T, 0, len(in))
	for _, v := range in {
		rtn = append(rtn, v)
	}
	return rtn
}
