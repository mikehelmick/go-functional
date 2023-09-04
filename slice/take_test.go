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
