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

func TestMinBy(t *testing.T) {
	t.Parallel()

	type person struct {
		Name string
		Age  int
	}
	byAge := func(p person) int { return p.Age }

	cases := []struct {
		name      string
		in        []person
		wantVal   person
		wantFound bool
	}{
		{
			name:      "empty",
			in:        []person{},
			wantVal:   person{},
			wantFound: false,
		},
		{
			name:      "single",
			in:        []person{{"Alice", 30}},
			wantVal:   person{"Alice", 30},
			wantFound: true,
		},
		{
			name:      "multiple",
			in:        []person{{"Bob", 30}, {"Alice", 20}, {"Carol", 40}},
			wantVal:   person{"Alice", 20},
			wantFound: true,
		},
		{
			name:      "returns_first_on_tie",
			in:        []person{{"Alice", 20}, {"Bob", 20}, {"Carol", 40}},
			wantVal:   person{"Alice", 20},
			wantFound: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, found := slice.MinBy(tc.in, byAge)
			if diff := cmp.Diff(tc.wantFound, found); diff != "" {
				t.Fatalf("found mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantVal, got); diff != "" {
				t.Fatalf("value mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestMaxBy(t *testing.T) {
	t.Parallel()

	type person struct {
		Name string
		Age  int
	}
	byAge := func(p person) int { return p.Age }

	cases := []struct {
		name      string
		in        []person
		wantVal   person
		wantFound bool
	}{
		{
			name:      "empty",
			in:        []person{},
			wantVal:   person{},
			wantFound: false,
		},
		{
			name:      "single",
			in:        []person{{"Alice", 30}},
			wantVal:   person{"Alice", 30},
			wantFound: true,
		},
		{
			name:      "multiple",
			in:        []person{{"Bob", 30}, {"Carol", 40}, {"Alice", 20}},
			wantVal:   person{"Carol", 40},
			wantFound: true,
		},
		{
			name:      "returns_first_on_tie",
			in:        []person{{"Alice", 40}, {"Bob", 40}, {"Carol", 20}},
			wantVal:   person{"Alice", 40},
			wantFound: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, found := slice.MaxBy(tc.in, byAge)
			if diff := cmp.Diff(tc.wantFound, found); diff != "" {
				t.Fatalf("found mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantVal, got); diff != "" {
				t.Fatalf("value mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleMinBy() {
	type person struct {
		Name string
		Age  int
	}
	people := []person{{"Bob", 30}, {"Alice", 20}, {"Carol", 40}}
	youngest, _ := slice.MinBy(people, func(p person) int { return p.Age })
	fmt.Printf("%s: %d\n", youngest.Name, youngest.Age)
	// Output:
	// Alice: 20
}

func ExampleMaxBy() {
	type person struct {
		Name string
		Age  int
	}
	people := []person{{"Bob", 30}, {"Alice", 20}, {"Carol", 40}}
	oldest, _ := slice.MaxBy(people, func(p person) int { return p.Age })
	fmt.Printf("%s: %d\n", oldest.Name, oldest.Age)
	// Output:
	// Carol: 40
}
