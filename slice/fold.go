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

// AccFn is the type of the accumulator function,
// used on various fold operations.
type AccFn[T any, A any] func(T, A) A

// FoldL folds (reduces) the input slice from the left.
// Requires an initial accumulator and an accumulator function.
func FoldL[T any, A any](in []T, acc A, fn AccFn[T, A]) A {
	for _, v := range in {
		acc = fn(v, acc)
	}
	return acc
}

// FoldR folds (reduces) the input slice from the right.
// Requires an initial accumulator and an accumulator function.
func FoldR[T any, A any](in []T, acc A, fn AccFn[T, A]) A {
	for i := len(in) - 1; i >= 0; i-- {
		acc = fn(in[i], acc)
	}
	return acc
}
