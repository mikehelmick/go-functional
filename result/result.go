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

package result

// Result holds either a successful value of type T or an error.
type Result[T any] struct {
	value T
	err   error
}

// OK returns a successful Result containing v.
func OK[T any](v T) Result[T] {
	return Result[T]{value: v}
}

// Err returns a failed Result containing err. Panics if err is nil.
func Err[T any](err error) Result[T] {
	if err == nil {
		panic("result.Err: err must not be nil")
	}
	return Result[T]{err: err}
}

// Of wraps the idiomatic Go (value, error) pair into a Result.
func Of[T any](v T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return OK(v)
}

// IsOK reports whether the Result holds a value.
func (r Result[T]) IsOK() bool {
	return r.err == nil
}

// IsErr reports whether the Result holds an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Get returns the value and nil error on success, or the zero value of T
// and the error on failure.
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// MustGet returns the value. It panics if the Result holds an error.
func (r Result[T]) MustGet() T {
	if r.err != nil {
		panic("result.Result: MustGet called on error Result: " + r.err.Error())
	}
	return r.value
}

// Error returns the error, or nil if the Result is successful.
func (r Result[T]) Error() error {
	return r.err
}

// OrElse returns the value if successful, otherwise returns defaultVal.
func (r Result[T]) OrElse(defaultVal T) T {
	if r.err == nil {
		return r.value
	}
	return defaultVal
}

// Map applies fn to the value inside r and returns a new Result containing
// the result. If r is an error, Map returns an error Result of type R.
func Map[T any, R any](r Result[T], fn func(T) R) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return OK(fn(r.value))
}

// FlatMap applies fn to the value inside r and returns the resulting Result.
// If r is an error, FlatMap returns an error Result of type R.
func FlatMap[T any, R any](r Result[T], fn func(T) Result[R]) Result[R] {
	if r.err != nil {
		return Err[R](r.err)
	}
	return fn(r.value)
}
