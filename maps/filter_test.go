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

func TestFilterMap(t *testing.T) {
	t.Parallel()

	positiveValue := func(_ string, v int) bool { return v > 0 }

	cases := []struct {
		name string
		in   map[string]int
		want map[string]int
	}{
		{
			name: "empty",
			in:   map[string]int{},
			want: map[string]int{},
		},
		{
			name: "all_pass",
			in:   map[string]int{"a": 1, "b": 2},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "none_pass",
			in:   map[string]int{"a": -1, "b": -2},
			want: map[string]int{},
		},
		{
			name: "partial",
			in:   map[string]int{"a": 1, "b": -2, "c": 3},
			want: map[string]int{"a": 1, "c": 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := maps.FilterMap(tc.in, positiveValue)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleFilterMap() {
	scores := map[string]int{"Alice": 90, "Bob": 45, "Carol": 75}
	passing := maps.FilterMap(scores, func(_ string, v int) bool { return v >= 60 })
	fmt.Printf("Alice passing: %v\n", passing["Alice"])
	fmt.Printf("Bob passing: %v\n", passing["Bob"] != 0)
	fmt.Printf("Carol passing: %v\n", passing["Carol"])
	// Output:
	// Alice passing: 90
	// Bob passing: false
	// Carol passing: 75
}
