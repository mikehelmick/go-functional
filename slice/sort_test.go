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

func TestSortBy(t *testing.T) {
	t.Parallel()

	type person struct {
		Name string
		Age  int
	}

	byAge := func(p person) int { return p.Age }

	cases := []struct {
		name string
		in   []person
		want []person
	}{
		{
			name: "empty",
			in:   []person{},
			want: []person{},
		},
		{
			name: "already_sorted",
			in:   []person{{"Alice", 20}, {"Bob", 30}, {"Carol", 40}},
			want: []person{{"Alice", 20}, {"Bob", 30}, {"Carol", 40}},
		},
		{
			name: "reverse_sorted",
			in:   []person{{"Carol", 40}, {"Bob", 30}, {"Alice", 20}},
			want: []person{{"Alice", 20}, {"Bob", 30}, {"Carol", 40}},
		},
		{
			name: "unsorted",
			in:   []person{{"Bob", 30}, {"Alice", 20}, {"Carol", 40}},
			want: []person{{"Alice", 20}, {"Bob", 30}, {"Carol", 40}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			original := make([]person, len(tc.in))
			copy(original, tc.in)

			got := slice.SortBy(tc.in, byAge)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
			// original slice must not be modified
			if diff := cmp.Diff(original, tc.in); diff != "" {
				t.Fatalf("input was mutated (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleSortBy() {
	type person struct {
		Name string
		Age  int
	}
	people := []person{{"Bob", 30}, {"Alice", 20}, {"Carol", 40}}
	sorted := slice.SortBy(people, func(p person) int { return p.Age })
	for _, p := range sorted {
		fmt.Printf("%s: %d\n", p.Name, p.Age)
	}
	// Output:
	// Alice: 20
	// Bob: 30
	// Carol: 40
}
