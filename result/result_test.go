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

package result_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/result"
)

var errTest = errors.New("test error")

func TestOK(t *testing.T) {
	t.Parallel()
	r := result.OK(42)
	if !r.IsOK() {
		t.Fatal("expected IsOK() == true")
	}
	if r.IsErr() {
		t.Fatal("expected IsErr() == false")
	}
	v, err := r.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(42, v); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestErr(t *testing.T) {
	t.Parallel()
	r := result.Err[int](errTest)
	if r.IsOK() {
		t.Fatal("expected IsOK() == false")
	}
	if !r.IsErr() {
		t.Fatal("expected IsErr() == true")
	}
	v, err := r.Get()
	if !errors.Is(err, errTest) {
		t.Fatalf("expected errTest, got %v", err)
	}
	if diff := cmp.Diff(0, v); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestErrPanicsOnNil(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when passing nil to Err")
		}
	}()
	result.Err[int](nil)
}

func TestOf(t *testing.T) {
	t.Parallel()
	r := result.Of(strconv.Atoi("42"))
	if diff := cmp.Diff(42, r.MustGet()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	bad := result.Of(strconv.Atoi("not-a-number"))
	if !bad.IsErr() {
		t.Fatal("expected error result")
	}
}

func TestMustGetPanicsOnErr(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on MustGet of error Result")
		}
	}()
	result.Err[int](errTest).MustGet()
}

func TestOrElse(t *testing.T) {
	t.Parallel()
	if diff := cmp.Diff(42, result.OK(42).OrElse(0)); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(99, result.Err[int](errTest).OrElse(99)); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMap(t *testing.T) {
	t.Parallel()
	double := func(v int) int { return v * 2 }

	got := result.Map(result.OK(21), double)
	if diff := cmp.Diff(42, got.MustGet()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	propagated := result.Map(result.Err[int](errTest), double)
	if !errors.Is(propagated.Error(), errTest) {
		t.Fatalf("expected errTest to propagate, got %v", propagated.Error())
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()
	safeParse := func(s string) result.Result[int] {
		return result.Of(strconv.Atoi(s))
	}

	got := result.FlatMap(result.OK("42"), safeParse)
	if diff := cmp.Diff(42, got.MustGet()); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	fromErr := result.FlatMap(result.Err[string](errTest), safeParse)
	if !errors.Is(fromErr.Error(), errTest) {
		t.Fatalf("expected errTest to propagate, got %v", fromErr.Error())
	}

	badParse := result.FlatMap(result.OK("not-a-number"), safeParse)
	if !badParse.IsErr() {
		t.Fatal("expected error result from failed parse")
	}
}

func ExampleOK() {
	r := result.OK(42)
	if v, err := r.Get(); err == nil {
		fmt.Println(v)
	}
	// Output:
	// 42
}

func ExampleOf() {
	r := result.Of(strconv.Atoi("123"))
	fmt.Println(r.MustGet())
	// Output:
	// 123
}

func ExampleMap() {
	r := result.Map(result.OK(21), func(v int) int { return v * 2 })
	fmt.Println(r.MustGet())
	// Output:
	// 42
}
