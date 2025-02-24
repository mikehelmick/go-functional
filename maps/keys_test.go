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

package maps_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mikehelmick/go-functional/maps"
)

func TestKeys(t *testing.T) {
	t.Parallel()

	input := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
	}

	keys := maps.Keys(input)

	want := []string{"one", "two", "three", "four", "five", "six"}

	opts := cmp.Options{cmpopts.SortSlices(func(a, b string) bool { return a < b })}
	if diff := cmp.Diff(want, keys, opts); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}
