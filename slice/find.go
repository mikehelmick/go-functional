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

// Find returns the first element in in for which matchFn returns true,
// and a boolean indicating whether a match was found.
// If no element matches, the zero value of T and false are returned.
func Find[T any](in []T, matchFn MatchFn[T]) (T, bool) {
	for _, e := range in {
		if matchFn(e) {
			return e, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element in in for which matchFn
// returns true, and a boolean indicating whether a match was found.
// If no element matches, -1 and false are returned.
func FindIndex[T any](in []T, matchFn MatchFn[T]) (int, bool) {
	for i, e := range in {
		if matchFn(e) {
			return i, true
		}
	}
	return -1, false
}
