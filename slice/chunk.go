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

// ChunkEvery splits in into consecutive non-overlapping chunks of size n.
// The final chunk may be smaller than n if len(in) is not evenly divisible.
// Returns an empty slice if n <= 0 or in is empty.
func ChunkEvery[T any](in []T, n int) [][]T {
	if n <= 0 || len(in) == 0 {
		return [][]T{}
	}
	numChunks := (len(in) + n - 1) / n
	rtn := make([][]T, 0, numChunks)
	for i := 0; i < len(in); i += n {
		end := min(i+n, len(in))
		rtn = append(rtn, in[i:end])
	}
	return rtn
}

// ChunkBy splits in into consecutive chunks where all elements in a chunk
// return the same key according to keyFn. A new chunk begins each time the
// key changes.
func ChunkBy[T any, K comparable](in []T, keyFn KeyFn[T, K]) [][]T {
	if len(in) == 0 {
		return [][]T{}
	}
	rtn := make([][]T, 0)
	start := 0
	currentKey := keyFn(in[0])
	for i := 1; i < len(in); i++ {
		k := keyFn(in[i])
		if k != currentKey {
			rtn = append(rtn, in[start:i])
			start = i
			currentKey = k
		}
	}
	rtn = append(rtn, in[start:])
	return rtn
}
