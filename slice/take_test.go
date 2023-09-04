// Copyright 2023 the go-functional authors
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

package slice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/slice"
)

func TestTake(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   []int
		amt  int
		want []int
	}{
		{
			name: "zero_amount",
			in:   []int{1, 2, 3},
			amt:  0,
			want: []int{},
		},
		{
			name: "one_amt",
			in:   []int{1, 2, 3},
			amt:  1,
			want: []int{1},
		},
		{
			name: "mid_amt",
			in:   []int{1, 2, 3},
			amt:  2,
			want: []int{1, 2},
		},
		{
			name: "full_amt",
			in:   []int{1, 2, 3},
			amt:  3,
			want: []int{1, 2, 3},
		},
		{
			name: "over_amt",
			in:   []int{1, 2, 3},
			amt:  4,
			want: []int{1, 2, 3},
		},
		{
			name: "negative_one",
			in:   []int{1, 2, 3},
			amt:  -1,
			want: []int{3},
		},
		{
			name: "negative_mid",
			in:   []int{1, 2, 3},
			amt:  -2,
			want: []int{2, 3},
		},
		{
			name: "negative_full",
			in:   []int{1, 2, 3},
			amt:  -3,
			want: []int{1, 2, 3},
		},
		{
			name: "negative_over",
			in:   []int{1, 2, 3},
			amt:  -4,
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := slice.Take(tc.in, tc.amt)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleTakeEvery() {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	takeTwo := slice.TakeEvery(in, 2)
	fmt.Printf("Take 2: %+v\n", takeTwo)

	takeZero := slice.TakeEvery(in, 0)
	fmt.Printf("Take 0: %+v\n", takeZero)

	takeAll := slice.TakeEvery(in, 27)
	fmt.Printf("Take all: %+v\n", takeAll)

	// Output:
	// Take 2: [1 3 5 7 9]
	// Take 0: []
	// Take all: [1]
}

func ExampleTakeWhile() {
	in := []string{"dip", "drive", "dodge", "swerve"}

	dPrefix := slice.TakeWhile(in, func(s string) bool {
		return strings.HasPrefix(s, "d")
	})
	fmt.Printf("dPrefix: %+v\n", dPrefix)

	diPrefix := slice.TakeWhile(in, func(s string) bool {
		return strings.HasPrefix(s, "di")
	})
	fmt.Printf("diPrefix: %+v\n", diPrefix)

	sPrefix := slice.TakeWhile(in, func(s string) bool {
		return strings.HasPrefix(s, "s")
	})
	fmt.Printf("sPrefix: %+v\n", sPrefix)

	all := slice.TakeWhile(in, func(s string) bool { return true })
	fmt.Printf("all: %+v\n", all)

	// Output:
	// dPrefix: [dip drive dodge]
	// diPrefix: [dip]
	// sPrefix: []
	// all: [dip drive dodge swerve]
}
