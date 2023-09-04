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

// GroupBy groups the input slice according to the key function.
func GroupBy[T any, R comparable](in []T, keyFn KeyFn[T, R]) map[R][]T {
	rtn := make(map[R][]T, len(in))

	for _, v := range in {
		k := keyFn(v)
		if _, ok := rtn[k]; !ok {
			rtn[k] = make([]T, 0, 1)
		}
		rtn[k] = append(rtn[k], v)
	}

	return rtn
}
