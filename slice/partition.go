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

// Partition splits in into two slices in a single pass: the first return value
// contains elements for which matchFn returns true, the second contains the rest.
// It is equivalent to calling Filter and Reject together without iterating twice.
func Partition[T any](in []T, matchFn MatchFn[T]) ([]T, []T) {
	matches := make([]T, 0, len(in))
	rejects := make([]T, 0, len(in))
	for _, e := range in {
		if matchFn(e) {
			matches = append(matches, e)
		} else {
			rejects = append(rejects, e)
		}
	}
	return matches, rejects
}
