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

// Take returns an amount of elements from the beginning of a slice.
// If the amount provides is negative, it returns that amount of elements
// from the end of the slice. If the amount is zero, an empty slice
// is returned.
func Take[T any](in []T, amount int) []T {
	if amount == 0 {
		return make([]T, 0)
	} else if amount > 0 {
		if amount >= len(in) {
			return in
		}
		return in[0:amount]
	}
	// negative amount, take from the end.
	if amount*-1 >= len(in) {
		return in
	}
	return in[len(in)+amount:]
}

// TakeEvery returns a new slice with every nth element
// from the original slice. The first element is always included
// unless n is 0.
func TakeEvery[T any](in []T, n uint) []T {
	if n == 0 {
		return make([]T, 0)
	}
	rtn := make([]T, 0, len(in)/int(n))
	nInt := int(n)
	for i := 0; i < len(in); i = i + nInt {
		rtn = append(rtn, in[i])
	}
	return rtn
}

// TakeWhile returns elements from the beginning of the
// input slice as long as the provided MatchFn returns true.
func TakeWhile[T any](in []T, fn MatchFn[T]) []T {
	idx := 0
	//revive:disable:empty-block
	for ; idx < len(in) && fn(in[idx]); idx++ {
	}
	//revive:enable:empty-block
	return in[0:idx]
}
