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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/slice"
)

func TestFrequencies(t *testing.T) {

	keyFn := func(x string) string {
		return strings.ToLower(x)
	}

	input := strings.Split("I have to practice my times tables over and over and over again so I can learn them", " ")

	got := slice.FrequenciesBy(input, keyFn)

	want := map[string]int{
		"again":    1,
		"and":      2,
		"can":      1,
		"have":     1,
		"i":        2,
		"learn":    1,
		"my":       1,
		"over":     3,
		"practice": 1,
		"so":       1,
		"tables":   1,
		"them":     1,
		"times":    1,
		"to":       1,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	// Same thing w/ frequencies
	got = slice.Frequencies(input)
	want = map[string]int{
		"again":    1,
		"and":      2,
		"can":      1,
		"have":     1,
		"I":        2,
		"learn":    1,
		"my":       1,
		"over":     3,
		"practice": 1,
		"so":       1,
		"tables":   1,
		"them":     1,
		"times":    1,
		"to":       1,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}
