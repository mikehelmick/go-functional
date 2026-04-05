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

package agent_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/agent"
)

func TestGetReturnsInitialState(t *testing.T) {
	t.Parallel()
	a := agent.New(42)
	defer a.Stop()
	if diff := cmp.Diff(42, a.Get()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	a := agent.New(0)
	defer a.Stop()
	a.Update(func(s int) int { return s + 10 })
	if diff := cmp.Diff(10, a.Get()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestCast(t *testing.T) {
	t.Parallel()
	a := agent.New(0)
	defer a.Stop()
	a.Cast(func(s int) int { return s + 1 })
	// Synchronize by calling Get, which is serialized after Cast.
	if diff := cmp.Diff(1, a.Get()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestGetWith(t *testing.T) {
	t.Parallel()
	a := agent.New([]string{"hello", "world"})
	defer a.Stop()
	length := agent.GetWith(a, func(s []string) int { return len(s) })
	if diff := cmp.Diff(2, length); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestStopIsIdempotent(t *testing.T) {
	t.Parallel()
	a := agent.New(0)
	a.Stop()
	a.Stop() // must not panic or deadlock
}

func TestConcurrentUpdates(t *testing.T) {
	t.Parallel()
	a := agent.New(0)
	defer a.Stop()

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.Update(func(s int) int { return s + 1 })
		}()
	}
	wg.Wait()

	if diff := cmp.Diff(100, a.Get()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func ExampleAgent_Get() {
	a := agent.New(42)
	defer a.Stop()
	fmt.Println(a.Get())
	// Output:
	// 42
}

func ExampleAgent_Update() {
	a := agent.New(0)
	defer a.Stop()
	a.Update(func(s int) int { return s + 10 })
	fmt.Println(a.Get())
	// Output:
	// 10
}

func ExampleGetWith() {
	a := agent.New(map[string]int{"a": 1, "b": 2, "c": 3})
	defer a.Stop()
	size := agent.GetWith(a, func(m map[string]int) int { return len(m) })
	fmt.Println(size)
	// Output:
	// 3
}
