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

package slice

// MatchFn indicates if an element of a slice matches
// the inclusion function.
type MatchFn[T any] func(T) bool

// All returns true if all elements in the input slice
// are true according to the provided matchFn.
func All[T any](in []T, matchFn MatchFn[T]) bool {
	for _, v := range in {
		if !matchFn(v) {
			return false
		}
	}
	return true
}

// Any returns true if any of the elements in the
// input slice return true when the provided matchFn
// is called.
func Any[T any](in []T, matchFn MatchFn[T]) bool {
	for _, v := range in {
		if matchFn(v) {
			return true
		}
	}
	return false
}
