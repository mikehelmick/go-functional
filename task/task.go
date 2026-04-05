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

package task

import "sync"

type taskResult[T any] struct {
	value T
	err   error
}

// Task represents an asynchronous computation of type T. A Task is safe to
// Await from multiple goroutines; all callers receive the same result.
type Task[T any] struct {
	once   sync.Once
	ch     chan taskResult[T]
	result taskResult[T]
}

// Run starts fn in a new goroutine and returns a Task to retrieve the result.
func Run[T any](fn func() (T, error)) *Task[T] {
	t := &Task[T]{ch: make(chan taskResult[T], 1)}
	go func() {
		v, err := fn()
		t.ch <- taskResult[T]{value: v, err: err}
	}()
	return t
}

func (t *Task[T]) await() taskResult[T] {
	t.once.Do(func() {
		t.result = <-t.ch
	})
	return t.result
}

// Await blocks until the task is complete and returns the value and any error.
// It is safe to call Await multiple times or from multiple goroutines.
func (t *Task[T]) Await() (T, error) {
	r := t.await()
	return r.value, r.err
}

// MustAwait blocks until the task is complete and returns the value.
// It panics if the task returned an error.
func (t *Task[T]) MustAwait() T {
	v, err := t.Await()
	if err != nil {
		panic("task: MustAwait called on failed Task: " + err.Error())
	}
	return v
}

// AwaitAll waits for all tasks and returns their values in order.
// The first error encountered is returned; remaining tasks continue running.
func AwaitAll[T any](tasks []*Task[T]) ([]T, error) {
	results := make([]T, len(tasks))
	for i, t := range tasks {
		v, err := t.Await()
		if err != nil {
			return nil, err
		}
		results[i] = v
	}
	return results, nil
}

// Map starts a new Task that applies fn to the result of t when it completes.
// If t fails, the mapped task propagates the error without calling fn.
// This is a package-level function because Go methods cannot introduce
// additional type parameters.
func Map[T, R any](t *Task[T], fn func(T) R) *Task[R] {
	return Run(func() (R, error) {
		v, err := t.Await()
		if err != nil {
			var zero R
			return zero, err
		}
		return fn(v), nil
	})
}
