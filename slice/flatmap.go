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

// FlatMapFn maps an individual instance of type T into a slice of type R.
type FlatMapFn[T any, R any] func(T) []R

// FlatMap applies fn to each element of in and concatenates the resulting
// slices into a single slice of type R.
func FlatMap[T any, R any](in []T, fn FlatMapFn[T, R]) []R {
	rtn := make([]R, 0, len(in))
	for _, e := range in {
		rtn = append(rtn, fn(e)...)
	}
	return rtn
}

// Flatten takes a slice of slices and concatenates them into a single slice.
func Flatten[T any](in [][]T) []T {
	total := 0
	for _, s := range in {
		total += len(s)
	}
	rtn := make([]T, 0, total)
	for _, s := range in {
		rtn = append(rtn, s...)
	}
	return rtn
}
