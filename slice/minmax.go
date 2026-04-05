// Copyright 2026 the go-functional authors
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

// Parts of this file were written with the assistance of Claude Code (claude.ai/code).

package slice

import "cmp"

// MinBy returns the element of in with the smallest key according to keyFn,
// and true. Returns the zero value of T and false if in is empty.
// When multiple elements share the minimum key, the first is returned.
func MinBy[T any, K cmp.Ordered](in []T, keyFn func(T) K) (T, bool) {
	if len(in) == 0 {
		var zero T
		return zero, false
	}
	minElem := in[0]
	minKey := keyFn(in[0])
	for i := 1; i < len(in); i++ {
		if k := keyFn(in[i]); k < minKey {
			minKey = k
			minElem = in[i]
		}
	}
	return minElem, true
}

// MaxBy returns the element of in with the largest key according to keyFn,
// and true. Returns the zero value of T and false if in is empty.
// When multiple elements share the maximum key, the first is returned.
func MaxBy[T any, K cmp.Ordered](in []T, keyFn func(T) K) (T, bool) {
	if len(in) == 0 {
		var zero T
		return zero, false
	}
	maxElem := in[0]
	maxKey := keyFn(in[0])
	for i := 1; i < len(in); i++ {
		if k := keyFn(in[i]); k > maxKey {
			maxKey = k
			maxElem = in[i]
		}
	}
	return maxElem, true
}
