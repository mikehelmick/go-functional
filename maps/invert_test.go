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

func TestInvert(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   map[string]int
		want map[int]string
	}{
		{
			name: "empty",
			in:   map[string]int{},
			want: map[int]string{},
		},
		{
			name: "unique_values",
			in:   map[string]int{"a": 1, "b": 2, "c": 3},
			want: map[int]string{1: "a", 2: "b", 3: "c"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := maps.Invert(tc.in)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func ExampleInvert() {
	codes := map[string]int{"USD": 1, "EUR": 2, "GBP": 3}
	byCode := maps.Invert(codes)
	fmt.Printf("%s\n", byCode[1])
	fmt.Printf("%s\n", byCode[2])
	// Output:
	// USD
	// EUR
}
