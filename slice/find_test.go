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

func TestFind(t *testing.T) {
	t.Parallel()

	isEven := func(x int) bool { return x%2 == 0 }

	cases := []struct {
		name      string
		in        []int
		wantVal   int
		wantFound bool
	}{
		{
			name:      "empty",
			in:        []int{},
			wantVal:   0,
			wantFound: false,
		},
		{
			name:      "no_match",
			in:        []int{1, 3, 5},
			wantVal:   0,
			wantFound: false,
		},
		{
			name:      "first_element",
			in:        []int{2, 3, 5},
			wantVal:   2,
			wantFound: true,
		},
		{
			name:      "middle_element",
			in:        []int{1, 4, 5},
			wantVal:   4,
			wantFound: true,
		},
		{
			name:      "returns_first_match",
			in:        []int{1, 2, 4},
			wantVal:   2,
			wantFound: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotVal, gotFound := slice.Find(tc.in, isEven)
			if diff := cmp.Diff(tc.wantVal, gotVal); diff != "" {
				t.Fatalf("value mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantFound, gotFound); diff != "" {
				t.Fatalf("found mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestFindIndex(t *testing.T) {
	t.Parallel()

	isEven := func(x int) bool { return x%2 == 0 }

	cases := []struct {
		name      string
		in        []int
		wantIdx   int
		wantFound bool
	}{
		{
			name:      "empty",
			in:        []int{},
			wantIdx:   -1,
			wantFound: false,
		},
		{
			name:      "no_match",
			in:        []int{1, 3, 5},
			wantIdx:   -1,
			wantFound: false,
		},
		{
			name:      "first_element",
			in:        []int{2, 3, 5},
			wantIdx:   0,
			wantFound: true,
		},
		{
			name:      "middle_element",
			in:        []int{1, 4, 5},
			wantIdx:   1,
			wantFound: true,
		},
		{
			name:      "returns_first_match_index",
			in:        []int{1, 2, 4},
			wantIdx:   1,
			wantFound: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotIdx, gotFound := slice.FindIndex(tc.in, isEven)
			if diff := cmp.Diff(tc.wantIdx, gotIdx); diff != "" {
				t.Fatalf("index mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantFound, gotFound); diff != "" {
				t.Fatalf("found mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleFind() {
	isEven := func(x int) bool { return x%2 == 0 }
	in := []int{1, 3, 4, 6}
	val, found := slice.Find(in, isEven)
	fmt.Printf("Find even in %v: val=%v found=%v\n", in, val, found)

	noMatch := []int{1, 3, 5}
	val, found = slice.Find(noMatch, isEven)
	fmt.Printf("Find even in %v: val=%v found=%v\n", noMatch, val, found)
	// Output:
	// Find even in [1 3 4 6]: val=4 found=true
	// Find even in [1 3 5]: val=0 found=false
}

func ExampleFindIndex() {
	isEven := func(x int) bool { return x%2 == 0 }
	in := []int{1, 3, 4, 6}
	idx, found := slice.FindIndex(in, isEven)
	fmt.Printf("FindIndex even in %v: idx=%v found=%v\n", in, idx, found)

	noMatch := []int{1, 3, 5}
	idx, found = slice.FindIndex(noMatch, isEven)
	fmt.Printf("FindIndex even in %v: idx=%v found=%v\n", noMatch, idx, found)
	// Output:
	// FindIndex even in [1 3 4 6]: idx=2 found=true
	// FindIndex even in [1 3 5]: idx=-1 found=false
}
