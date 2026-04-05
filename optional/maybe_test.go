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

package optional_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/optional"
)

func TestSomeIsPresent(t *testing.T) {
	t.Parallel()
	m := optional.Some(42)
	if !m.IsPresent() {
		t.Fatal("expected IsPresent() == true")
	}
	if m.IsEmpty() {
		t.Fatal("expected IsEmpty() == false")
	}
	v, ok := m.Get()
	if !ok {
		t.Fatal("expected ok == true")
	}
	if diff := cmp.Diff(42, v); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestNoneIsEmpty(t *testing.T) {
	t.Parallel()
	m := optional.None[int]()
	if m.IsPresent() {
		t.Fatal("expected IsPresent() == false")
	}
	if !m.IsEmpty() {
		t.Fatal("expected IsEmpty() == true")
	}
	v, ok := m.Get()
	if ok {
		t.Fatal("expected ok == false")
	}
	if diff := cmp.Diff(0, v); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMustGet(t *testing.T) {
	t.Parallel()
	if diff := cmp.Diff(42, optional.Some(42).MustGet()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMustGetPanicsOnNone(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on MustGet of empty Maybe")
		}
	}()
	optional.None[int]().MustGet()
}

func TestOrElse(t *testing.T) {
	t.Parallel()
	if diff := cmp.Diff(42, optional.Some(42).OrElse(0)); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(99, optional.None[int]().OrElse(99)); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMap(t *testing.T) {
	t.Parallel()
	double := func(v int) int { return v * 2 }

	got := optional.Map(optional.Some(21), double)
	if diff := cmp.Diff(42, got.MustGet()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	empty := optional.Map(optional.None[int](), double)
	if !empty.IsEmpty() {
		t.Fatal("expected empty Maybe from mapping over None")
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()
	safeSqrt := func(v float64) optional.Maybe[float64] {
		if v < 0 {
			return optional.None[float64]()
		}
		return optional.Some(v * v)
	}

	got := optional.FlatMap(optional.Some(4.0), safeSqrt)
	if diff := cmp.Diff(16.0, got.MustGet()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	rejected := optional.FlatMap(optional.Some(-1.0), safeSqrt)
	if !rejected.IsEmpty() {
		t.Fatal("expected empty Maybe when fn returns None")
	}

	empty := optional.FlatMap(optional.None[float64](), safeSqrt)
	if !empty.IsEmpty() {
		t.Fatal("expected empty Maybe from FlatMap over None")
	}
}

func ExampleSome() {
	m := optional.Some("hello")
	if v, ok := m.Get(); ok {
		fmt.Println(v)
	}
	// Output:
	// hello
}

func ExampleNone() {
	m := optional.None[string]()
	fmt.Println(m.OrElse("default"))
	// Output:
	// default
}

func ExampleMap() {
	length := func(s string) int { return len(s) }
	fmt.Println(optional.Map(optional.Some("hello"), length).MustGet())
	fmt.Println(optional.Map(optional.None[string](), length).IsEmpty())
	// Output:
	// 5
	// true
}
