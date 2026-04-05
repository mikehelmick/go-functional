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

// Dedup removes consecutive duplicate elements from in, keeping the first
// occurrence of each run. Only adjacent duplicates are removed; use
// Frequencies or a sort before Dedup to remove all duplicates.
func Dedup[T comparable](in []T) []T {
	if len(in) == 0 {
		return []T{}
	}
	rtn := make([]T, 0, len(in))
	rtn = append(rtn, in[0])
	for i := 1; i < len(in); i++ {
		if in[i] != in[i-1] {
			rtn = append(rtn, in[i])
		}
	}
	return rtn
}

// DedupBy removes consecutive elements from in that produce the same key
// according to keyFn, keeping the first element of each run.
func DedupBy[T any, K comparable](in []T, keyFn KeyFn[T, K]) []T {
	if len(in) == 0 {
		return []T{}
	}
	rtn := make([]T, 0, len(in))
	rtn = append(rtn, in[0])
	currentKey := keyFn(in[0])
	for i := 1; i < len(in); i++ {
		k := keyFn(in[i])
		if k != currentKey {
			rtn = append(rtn, in[i])
			currentKey = k
		}
	}
	return rtn
}
