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

package slice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/slice"
)

func TestFlatMap(t *testing.T) {
	t.Parallel()

	split := func(s string) []string { return strings.Split(s, ",") }

	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "empty",
			in:   []string{},
			want: []string{},
		},
		{
			name: "single_element_single_result",
			in:   []string{"a"},
			want: []string{"a"},
		},
		{
			name: "single_element_multiple_results",
			in:   []string{"a,b,c"},
			want: []string{"a", "b", "c"},
		},
		{
			name: "multiple_elements",
			in:   []string{"a,b", "c,d"},
			want: []string{"a", "b", "c", "d"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.FlatMap(tc.in, split)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   [][]int
		want []int
	}{
		{
			name: "empty",
			in:   [][]int{},
			want: []int{},
		},
		{
			name: "single_inner_slice",
			in:   [][]int{{1, 2, 3}},
			want: []int{1, 2, 3},
		},
		{
			name: "multiple_inner_slices",
			in:   [][]int{{1, 2}, {3, 4}, {5}},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty_inner_slices",
			in:   [][]int{{}, {1, 2}, {}},
			want: []int{1, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.Flatten(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleFlatMap() {
	words := []string{"hello world", "foo bar"}
	split := func(s string) []string { return strings.Split(s, " ") }
	got := slice.FlatMap(words, split)
	fmt.Printf("%v\n", got)
	// Output:
	// [hello world foo bar]
}

func ExampleFlatten() {
	in := [][]int{{1, 2}, {3, 4}, {5}}
	got := slice.Flatten(in)
	fmt.Printf("%v\n", got)
	// Output:
	// [1 2 3 4 5]
}
