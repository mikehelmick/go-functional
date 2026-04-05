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

func TestZip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		a    []int
		b    []string
		want []slice.Pair[int, string]
	}{
		{
			name: "empty",
			a:    []int{},
			b:    []string{},
			want: []slice.Pair[int, string]{},
		},
		{
			name: "equal_length",
			a:    []int{1, 2, 3},
			b:    []string{"a", "b", "c"},
			want: []slice.Pair[int, string]{{First: 1, Second: "a"}, {First: 2, Second: "b"}, {First: 3, Second: "c"}},
		},
		{
			name: "first_shorter",
			a:    []int{1, 2},
			b:    []string{"a", "b", "c"},
			want: []slice.Pair[int, string]{{First: 1, Second: "a"}, {First: 2, Second: "b"}},
		},
		{
			name: "second_shorter",
			a:    []int{1, 2, 3},
			b:    []string{"a"},
			want: []slice.Pair[int, string]{{First: 1, Second: "a"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.Zip(tc.a, tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestUnzip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		in     []slice.Pair[int, string]
		wantAs []int
		wantBs []string
	}{
		{
			name:   "empty",
			in:     []slice.Pair[int, string]{},
			wantAs: []int{},
			wantBs: []string{},
		},
		{
			name:   "multiple",
			in:     []slice.Pair[int, string]{{First: 1, Second: "a"}, {First: 2, Second: "b"}, {First: 3, Second: "c"}},
			wantAs: []int{1, 2, 3},
			wantBs: []string{"a", "b", "c"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotAs, gotBs := slice.Unzip(tc.in)
			if diff := cmp.Diff(tc.wantAs, gotAs); diff != "" {
				t.Fatalf("first slice mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantBs, gotBs); diff != "" {
				t.Fatalf("second slice mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleZip() {
	names := []string{"Alice", "Bob", "Carol"}
	scores := []int{95, 87, 92}
	pairs := slice.Zip(names, scores)
	for _, p := range pairs {
		fmt.Printf("%s: %d\n", p.First, p.Second)
	}
	// Output:
	// Alice: 95
	// Bob: 87
	// Carol: 92
}

func ExampleUnzip() {
	pairs := []slice.Pair[string, int]{
		{First: "Alice", Second: 95},
		{First: "Bob", Second: 87},
	}
	names, scores := slice.Unzip(pairs)
	fmt.Printf("Names:  %v\n", names)
	fmt.Printf("Scores: %v\n", scores)
	// Output:
	// Names:  [Alice Bob]
	// Scores: [95 87]
}
