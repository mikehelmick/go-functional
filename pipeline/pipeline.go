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

package pipeline

// Pipe composes zero or more same-type functions left-to-right, returning a
// single function that applies each in sequence. Pipe(f, g, h)(x) == h(g(f(x))).
// If no functions are provided, Pipe returns the identity function.
func Pipe[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for _, fn := range fns {
			v = fn(v)
		}
		return v
	}
}

// Pipe2 composes two functions left-to-right where types may differ.
// Pipe2(f, g)(x) == g(f(x)).
func Pipe2[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Pipe3 composes three functions left-to-right where types may differ.
// Pipe3(f, g, h)(x) == h(g(f(x))).
func Pipe3[A, B, C, D any](f func(A) B, g func(B) C, h func(C) D) func(A) D {
	return func(a A) D {
		return h(g(f(a)))
	}
}

// Pipe4 composes four functions left-to-right where types may differ.
// Pipe4(f, g, h, i)(x) == i(h(g(f(x)))).
func Pipe4[A, B, C, D, E any](f func(A) B, g func(B) C, h func(C) D, i func(D) E) func(A) E {
	return func(a A) E {
		return i(h(g(f(a))))
	}
}
