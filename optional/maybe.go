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

package optional

// Maybe represents an optional value of type T. A Maybe is either Some
// (holding a value) or None (holding nothing).
type Maybe[T any] struct {
	value   T
	present bool
}

// Some returns a Maybe containing the given value.
func Some[T any](v T) Maybe[T] {
	return Maybe[T]{value: v, present: true}
}

// None returns an empty Maybe of type T.
func None[T any]() Maybe[T] {
	return Maybe[T]{}
}

// IsPresent reports whether the Maybe holds a value.
func (m Maybe[T]) IsPresent() bool {
	return m.present
}

// IsEmpty reports whether the Maybe holds no value.
func (m Maybe[T]) IsEmpty() bool {
	return !m.present
}

// Get returns the value and true if present, or the zero value of T and
// false if empty.
func (m Maybe[T]) Get() (T, bool) {
	return m.value, m.present
}

// MustGet returns the value. It panics if the Maybe is empty.
func (m Maybe[T]) MustGet() T {
	if !m.present {
		panic("optional.Maybe: MustGet called on empty Maybe")
	}
	return m.value
}

// OrElse returns the value if present, otherwise returns the provided default.
func (m Maybe[T]) OrElse(defaultVal T) T {
	if m.present {
		return m.value
	}
	return defaultVal
}

// Map applies fn to the value inside m and returns a new Maybe containing
// the result. If m is empty, Map returns an empty Maybe of type R.
func Map[T any, R any](m Maybe[T], fn func(T) R) Maybe[R] {
	if !m.present {
		return None[R]()
	}
	return Some(fn(m.value))
}

// FlatMap applies fn to the value inside m and returns the resulting Maybe.
// If m is empty, FlatMap returns an empty Maybe of type R.
func FlatMap[T any, R any](m Maybe[T], fn func(T) Maybe[R]) Maybe[R] {
	if !m.present {
		return None[R]()
	}
	return fn(m.value)
}
