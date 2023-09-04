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

// MapFn maps an individual instance of a type T into an
// individual instance of type R.
type MapFn[T any, R any] func(T) R

// Map takes in a slice of type T, and a map function that
// can turn elements of type T into elements of type R.
// A new slice of type R is returned, withe the map function
// applied to each element in the input.
func Map[T any, R any](in []T, fn MapFn[T, R]) []R {
	rtn := make([]R, len(in))
	for i, t := range in {
		rtn[i] = fn(t)
	}
	return rtn
}
