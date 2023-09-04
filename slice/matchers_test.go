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

	"github.com/mikehelmick/go-functional/slice"
)

func ExampleAll() {
	// Function to check if an element is even
	isEven := func(x int) bool {
		return x%2 == 0
	}

	allEven := []int{2, 4, 6, 8, 10}
	fmt.Printf("All(%v) == %v\n", allEven, slice.All(allEven, isEven))

	withOds := []int{2, 4, 5, 6}
	fmt.Printf("All(%v) == %v\n", withOds, slice.All(withOds, isEven))

	// Output:
	// All([2 4 6 8 10]) == true
	// All([2 4 5 6]) == false
}

func ExampleAny() {
	// Function to check if an element is seven
	isSeven := func(x int) bool {
		return x == 7
	}

	noSevens := []int{5, 6, 8}
	fmt.Printf("All(%v) == %v\n", noSevens, slice.Any(noSevens, isSeven))

	withSevens := []int{5, 6, 7, 8}
	fmt.Printf("All(%v) == %v\n", withSevens, slice.Any(withSevens, isSeven))

	// Output:
	// All([5 6 8]) == false
	// All([5 6 7 8]) == true
}
