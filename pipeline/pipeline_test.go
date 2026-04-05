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

package pipeline_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/pipeline"
)

func TestPipe(t *testing.T) {
	t.Parallel()

	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }

	cases := []struct {
		name string
		fns  []func(int) int
		in   int
		want int
	}{
		{
			name: "no_functions_identity",
			fns:  []func(int) int{},
			in:   42,
			want: 42,
		},
		{
			name: "single_function",
			fns:  []func(int) int{double},
			in:   5,
			want: 10,
		},
		{
			name: "left_to_right_order",
			fns:  []func(int) int{double, addOne},
			in:   5,
			want: 11, // (5*2)+1
		},
		{
			name: "three_functions",
			fns:  []func(int) int{addOne, double, addOne},
			in:   3,
			want: 9, // ((3+1)*2)+1
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := pipeline.Pipe(tc.fns...)(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestPipe2(t *testing.T) {
	t.Parallel()
	f := pipeline.Pipe2(
		func(s string) int { return len(s) },
		func(n int) bool { return n > 3 },
	)
	if diff := cmp.Diff(true, f("hello")); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(false, f("hi")); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestPipe3(t *testing.T) {
	t.Parallel()
	f := pipeline.Pipe3(
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) int { return len(s) },
	)
	if diff := cmp.Diff(5, f("  hello  ")); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestPipe4(t *testing.T) {
	t.Parallel()
	f := pipeline.Pipe4(
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) []string { return strings.Split(s, " ") },
		func(ss []string) int { return len(ss) },
	)
	if diff := cmp.Diff(3, f("  foo bar baz  ")); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func ExamplePipe() {
	transform := pipeline.Pipe(
		func(x int) int { return x * 2 },
		func(x int) int { return x + 1 },
	)
	fmt.Println(transform(5))
	// Output:
	// 11
}

func ExamplePipe2() {
	f := pipeline.Pipe2(
		strings.TrimSpace,
		strings.ToUpper,
	)
	fmt.Println(f("  hello  "))
	// Output:
	// HELLO
}
