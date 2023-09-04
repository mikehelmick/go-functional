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

// KeyFn determines how the input value should be
// used when determining frequencies.
type KeyFn[T any, R comparable] func(T) R

// Frequencies returns a map with a count of the number of time
// each element occurs in the input slice.
func Frequencies[T comparable](in []T) map[T]int {
	id := func(v T) T { return v }
	return FrequenciesBy[T, T](in, id)
}

// FrequenciesBy returns a map with a count of the number of times each
// key occurred in the input sequence. The KeyFn provided can manipulate
// the input elements in any way.
func FrequenciesBy[T any, R comparable](in []T, keyFn KeyFn[T, R]) map[R]int {
	rtn := make(map[R]int, len(in))
	for _, v := range in {
		key := keyFn(v)
		if cur, ok := rtn[key]; ok {
			rtn[key] = cur + 1
		} else {
			rtn[key] = 1
		}
	}
	return rtn
}
