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

func ExampleFoldL() {
	sumFn := func(x int, acc int) int {
		return acc + x
	}

	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := slice.FoldL(in, 0, sumFn)

	fmt.Printf("sum(%v) == %v\n", in, sum)

	// Output:
	// sum([1 2 3 4 5 6 7 8 9 10]) == 55
}

func ExampleFoldR() {
	subFn := func(x int, acc int) int {
		return acc - x
	}

	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := slice.FoldR(in, 20, subFn)

	fmt.Printf("foldr_sub(%v) == %v\n", in, sum)

	// Output:
	// foldr_sub([1 2 3 4 5 6 7 8 9 10]) == -35
}
