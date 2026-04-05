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

func TestMapValues(t *testing.T) {
	t.Parallel()

	double := func(v int) int { return v * 2 }

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
			name: "doubles_values",
			in:   map[string]int{"a": 1, "b": 2, "c": 3},
			want: map[string]int{"a": 2, "b": 4, "c": 6},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := maps.MapValues(tc.in, double)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleMapValues() {
	prices := map[string]float64{"apple": 1.00, "banana": 0.50}
	discounted := maps.MapValues(prices, func(v float64) float64 { return v * 0.9 })
	fmt.Printf("apple: %.2f\n", discounted["apple"])
	fmt.Printf("banana: %.2f\n", discounted["banana"])
	// Output:
	// apple: 0.90
	// banana: 0.45
}
