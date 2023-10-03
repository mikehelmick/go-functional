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
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/mikehelmick/go-functional/slice"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %v", p.Name, p.Age)
}

func ExampleMap() {
	redactAge := func(p Person) Person {
		return Person{
			Name: p.Name,
		}
	}

	people := []Person{
		{Name: "Jim", Age: 28},
		{Name: "Dwight", Age: 32},
	}

	redacted := slice.Map(people, redactAge)

	fmt.Printf("%+v\n", people)
	fmt.Printf("%+v\n", redacted)

	// Output:
	// [Jim: 28 Dwight: 32]
	// [Jim: 0 Dwight: 0]
}

func TestMap(t *testing.T) {
	t.Parallel()

	fn := func(a int) string {
		return fmt.Sprintf("%04d", a)
	}

	got := slice.Map[int, string]([]int{1, 2, 3, 4, 5}, fn)

	want := []string{
		"0001", "0002", "0003", "0004", "0005",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestMapPtr(t *testing.T) {

	in := []Person{
		{
			Name: "Alice",
			Age:  30,
		},
		{
			Name: "Bob",
			Age:  31,
		},
	}

	got := slice.MapToPtr(in)

	want := []*Person{&in[0], &in[1]}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}
