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

package maps_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/maps"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		a    map[string]int
		b    map[string]int
		want map[string]int
	}{
		{
			name: "both_empty",
			a:    map[string]int{},
			b:    map[string]int{},
			want: map[string]int{},
		},
		{
			name: "no_overlap",
			a:    map[string]int{"a": 1},
			b:    map[string]int{"b": 2},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "b_wins_on_conflict",
			a:    map[string]int{"a": 1, "b": 2},
			b:    map[string]int{"b": 99, "c": 3},
			want: map[string]int{"a": 1, "b": 99, "c": 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := maps.Merge(tc.a, tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestMergeWith(t *testing.T) {
	t.Parallel()

	sum := func(a, b int) int { return a + b }

	cases := []struct {
		name string
		a    map[string]int
		b    map[string]int
		want map[string]int
	}{
		{
			name: "both_empty",
			a:    map[string]int{},
			b:    map[string]int{},
			want: map[string]int{},
		},
		{
			name: "no_overlap",
			a:    map[string]int{"a": 1},
			b:    map[string]int{"b": 2},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "conflict_resolved_by_fn",
			a:    map[string]int{"a": 1, "b": 2},
			b:    map[string]int{"b": 3, "c": 4},
			want: map[string]int{"a": 1, "b": 5, "c": 4},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := maps.MergeWith(tc.a, tc.b, sum)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleMerge() {
	defaults := map[string]int{"timeout": 30, "retries": 3}
	overrides := map[string]int{"retries": 5, "verbose": 1}
	config := maps.Merge(defaults, overrides)
	fmt.Printf("timeout: %d\n", config["timeout"])
	fmt.Printf("retries: %d\n", config["retries"])
	fmt.Printf("verbose: %d\n", config["verbose"])
	// Output:
	// timeout: 30
	// retries: 5
	// verbose: 1
}

func ExampleMergeWith() {
	a := map[string]int{"x": 1, "y": 2}
	b := map[string]int{"y": 3, "z": 4}
	merged := maps.MergeWith(a, b, func(va, vb int) int { return va + vb })
	fmt.Printf("x: %d\n", merged["x"])
	fmt.Printf("y: %d\n", merged["y"])
	fmt.Printf("z: %d\n", merged["z"])
	// Output:
	// x: 1
	// y: 5
	// z: 4
}
