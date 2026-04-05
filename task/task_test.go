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

package task_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/task"
)

var errTest = errors.New("task error")

func TestRunAndAwait(t *testing.T) {
	t.Parallel()
	got, err := task.Run(func() (int, error) { return 42, nil }).Await()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(42, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestAwaitError(t *testing.T) {
	t.Parallel()
	_, err := task.Run(func() (int, error) { return 0, errTest }).Await()
	if !errors.Is(err, errTest) {
		t.Fatalf("expected errTest, got %v", err)
	}
}

func TestMustAwait(t *testing.T) {
	t.Parallel()
	got := task.Run(func() (int, error) { return 7, nil }).MustAwait()
	if diff := cmp.Diff(7, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMustAwaitPanicsOnError(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on MustAwait of failed task")
		}
	}()
	task.Run(func() (int, error) { return 0, errTest }).MustAwait()
}

func TestAwaitIdempotent(t *testing.T) {
	t.Parallel()
	tsk := task.Run(func() (int, error) { return 99, nil })
	for range 5 {
		got, err := tsk.Await()
		if err != nil || got != 99 {
			t.Fatalf("expected 99/nil, got %v/%v", got, err)
		}
	}
}

func TestAwaitConcurrent(t *testing.T) {
	t.Parallel()
	tsk := task.Run(func() (int, error) { return 42, nil })
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := tsk.Await()
			if err != nil || got != 42 {
				t.Errorf("expected 42/nil, got %v/%v", got, err)
			}
		}()
	}
	wg.Wait()
}

func TestAwaitAll(t *testing.T) {
	t.Parallel()
	tasks := []*task.Task[int]{
		task.Run(func() (int, error) { return 1, nil }),
		task.Run(func() (int, error) { return 2, nil }),
		task.Run(func() (int, error) { return 3, nil }),
	}
	got, err := task.AwaitAll(tasks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff([]int{1, 2, 3}, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestAwaitAllError(t *testing.T) {
	t.Parallel()
	tasks := []*task.Task[int]{
		task.Run(func() (int, error) { return 1, nil }),
		task.Run(func() (int, error) { return 0, errTest }),
		task.Run(func() (int, error) { return 3, nil }),
	}
	_, err := task.AwaitAll(tasks)
	if !errors.Is(err, errTest) {
		t.Fatalf("expected errTest, got %v", err)
	}
}

func TestMap(t *testing.T) {
	t.Parallel()
	tsk := task.Run(func() (int, error) { return 21, nil })
	mapped := task.Map(tsk, func(v int) string { return fmt.Sprintf("%d", v*2) })
	got, err := mapped.Await()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff("42", got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMapPropagatesError(t *testing.T) {
	t.Parallel()
	tsk := task.Run(func() (int, error) { return 0, errTest })
	mapped := task.Map(tsk, func(v int) string { return fmt.Sprintf("%d", v) })
	_, err := mapped.Await()
	if !errors.Is(err, errTest) {
		t.Fatalf("expected errTest to propagate, got %v", err)
	}
}

func ExampleRun() {
	t := task.Run(func() (int, error) { return 42, nil })
	v, _ := t.Await()
	fmt.Println(v)
	// Output:
	// 42
}

func ExampleAwaitAll() {
	tasks := []*task.Task[int]{
		task.Run(func() (int, error) { return 1, nil }),
		task.Run(func() (int, error) { return 2, nil }),
		task.Run(func() (int, error) { return 3, nil }),
	}
	results, _ := task.AwaitAll(tasks)
	fmt.Println(results)
	// Output:
	// [1 2 3]
}
