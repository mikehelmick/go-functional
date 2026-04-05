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

func TestPartition(t *testing.T) {
	t.Parallel()

	isEven := func(x int) bool { return x%2 == 0 }

	cases := []struct {
		name        string
		in          []int
		wantMatches []int
		wantRejects []int
	}{
		{
			name:        "empty",
			in:          []int{},
			wantMatches: []int{},
			wantRejects: []int{},
		},
		{
			name:        "all_match",
			in:          []int{2, 4, 6},
			wantMatches: []int{2, 4, 6},
			wantRejects: []int{},
		},
		{
			name:        "none_match",
			in:          []int{1, 3, 5},
			wantMatches: []int{},
			wantRejects: []int{1, 3, 5},
		},
		{
			name:        "partial_match",
			in:          []int{1, 2, 3, 4, 5},
			wantMatches: []int{2, 4},
			wantRejects: []int{1, 3, 5},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotMatches, gotRejects := slice.Partition(tc.in, isEven)
			if diff := cmp.Diff(tc.wantMatches, gotMatches); diff != "" {
				t.Fatalf("matches mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantRejects, gotRejects); diff != "" {
				t.Fatalf("rejects mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExamplePartition() {
	isEven := func(x int) bool { return x%2 == 0 }
	in := []int{1, 2, 3, 4, 5, 6}
	evens, odds := slice.Partition(in, isEven)
	fmt.Printf("Evens: %v\n", evens)
	fmt.Printf("Odds:  %v\n", odds)
	// Output:
	// Evens: [2 4 6]
	// Odds:  [1 3 5]
}
