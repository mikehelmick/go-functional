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

// Pair holds two values of potentially different types.
type Pair[A any, B any] struct {
	First  A
	Second B
}

// Zip combines two slices into a slice of Pairs. The result length equals
// the length of the shorter input slice.
func Zip[A any, B any](a []A, b []B) []Pair[A, B] {
	length := min(len(a), len(b))
	rtn := make([]Pair[A, B], length)
	for i := range length {
		rtn[i] = Pair[A, B]{First: a[i], Second: b[i]}
	}
	return rtn
}

// Unzip splits a slice of Pairs into two separate slices.
func Unzip[A any, B any](in []Pair[A, B]) ([]A, []B) {
	as := make([]A, len(in))
	bs := make([]B, len(in))
	for i, p := range in {
		as[i] = p.First
		bs[i] = p.Second
	}
	return as, bs
}
