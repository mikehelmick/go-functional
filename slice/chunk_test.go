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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/slice"
)

func TestChunkEvery(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   []int
		n    int
		want [][]int
	}{
		{
			name: "empty",
			in:   []int{},
			n:    2,
			want: [][]int{},
		},
		{
			name: "zero_size",
			in:   []int{1, 2, 3},
			n:    0,
			want: [][]int{},
		},
		{
			name: "negative_size",
			in:   []int{1, 2, 3},
			n:    -1,
			want: [][]int{},
		},
		{
			name: "even_chunks",
			in:   []int{1, 2, 3, 4},
			n:    2,
			want: [][]int{{1, 2}, {3, 4}},
		},
		{
			name: "uneven_chunks",
			in:   []int{1, 2, 3, 4, 5},
			n:    2,
			want: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name: "chunk_larger_than_input",
			in:   []int{1, 2},
			n:    5,
			want: [][]int{{1, 2}},
		},
		{
			name: "chunk_size_one",
			in:   []int{1, 2, 3},
			n:    1,
			want: [][]int{{1}, {2}, {3}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.ChunkEvery(tc.in, tc.n)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestChunkBy(t *testing.T) {
	t.Parallel()

	isEven := func(x int) bool { return x%2 == 0 }

	cases := []struct {
		name string
		in   []int
		want [][]int
	}{
		{
			name: "empty",
			in:   []int{},
			want: [][]int{},
		},
		{
			name: "all_same_key",
			in:   []int{2, 4, 6},
			want: [][]int{{2, 4, 6}},
		},
		{
			name: "alternating",
			in:   []int{1, 2, 3, 4},
			want: [][]int{{1}, {2}, {3}, {4}},
		},
		{
			name: "consecutive_groups",
			in:   []int{1, 3, 2, 4, 5},
			want: [][]int{{1, 3}, {2, 4}, {5}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.ChunkBy(tc.in, isEven)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleChunkEvery() {
	in := []int{1, 2, 3, 4, 5}
	fmt.Printf("%v\n", slice.ChunkEvery(in, 2))
	// Output:
	// [[1 2] [3 4] [5]]
}

func ExampleChunkBy() {
	in := []int{1, 3, 2, 4, 5}
	isEven := func(x int) bool { return x%2 == 0 }
	fmt.Printf("%v\n", slice.ChunkBy(in, isEven))
	// Output:
	// [[1 3] [2 4] [5]]
}
