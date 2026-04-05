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

// Scan is like FoldL but returns a slice of all intermediate accumulator
// values, including the initial value. The result always has len(in)+1
// elements, with the first element being the initial accumulator.
func Scan[T any, A any](in []T, acc A, fn AccFn[T, A]) []A {
	rtn := make([]A, len(in)+1)
	rtn[0] = acc
	for i, v := range in {
		acc = fn(v, acc)
		rtn[i+1] = acc
	}
	return rtn
}
