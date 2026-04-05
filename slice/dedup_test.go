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

func TestDedup(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{
			name: "empty",
			in:   []int{},
			want: []int{},
		},
		{
			name: "no_duplicates",
			in:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "all_same",
			in:   []int{1, 1, 1},
			want: []int{1},
		},
		{
			name: "consecutive_duplicates",
			in:   []int{1, 1, 2, 3, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "non_consecutive_duplicates_kept",
			in:   []int{1, 2, 1, 2},
			want: []int{1, 2, 1, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.Dedup(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestDedupBy(t *testing.T) {
	t.Parallel()

	isEven := func(x int) bool { return x%2 == 0 }

	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{
			name: "empty",
			in:   []int{},
			want: []int{},
		},
		{
			name: "no_consecutive_same_key",
			in:   []int{1, 2, 3, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "consecutive_same_key",
			in:   []int{1, 3, 2, 4, 5},
			want: []int{1, 2, 5},
		},
		{
			name: "keeps_first_of_run",
			in:   []int{2, 4, 6, 1, 3},
			want: []int{2, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.DedupBy(tc.in, isEven)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleDedup() {
	in := []int{1, 1, 2, 3, 3, 3, 2}
	fmt.Printf("%v\n", slice.Dedup(in))
	// Output:
	// [1 2 3 2]
}

func ExampleDedupBy() {
	in := []int{1, 3, 2, 4, 5}
	isEven := func(x int) bool { return x%2 == 0 }
	fmt.Printf("%v\n", slice.DedupBy(in, isEven))
	// Output:
	// [1 2 5]
}
