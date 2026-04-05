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

func TestScan(t *testing.T) {
	t.Parallel()

	sum := func(v, acc int) int { return acc + v }

	cases := []struct {
		name string
		in   []int
		acc  int
		want []int
	}{
		{
			name: "empty",
			in:   []int{},
			acc:  0,
			want: []int{0},
		},
		{
			name: "running_sum",
			in:   []int{1, 2, 3, 4},
			acc:  0,
			want: []int{0, 1, 3, 6, 10},
		},
		{
			name: "non_zero_initial",
			in:   []int{1, 2, 3},
			acc:  10,
			want: []int{10, 11, 13, 16},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.Scan(tc.in, tc.acc, sum)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleScan() {
	in := []int{1, 2, 3, 4, 5}
	sums := slice.Scan(in, 0, func(v, acc int) int { return acc + v })
	fmt.Printf("%v\n", sums)
	// Output:
	// [0 1 3 6 10 15]
}
