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

func ExampleFilter() {
	type Account struct {
		Name   string
		Writer bool
	}

	input := []*Account{
		{
			Name:   "Bob",
			Writer: false,
		},
		{
			Name:   "Alice",
			Writer: true,
		},
		{
			Name:   "Steve",
			Writer: false,
		},
	}

	// Function to check if an element is even
	isWriter := func(x *Account) bool {
		return x.Writer
	}

	writers := slice.Filter(input, isWriter)
	fmt.Printf("Writers:\n")
	for _, writer := range writers {
		fmt.Printf(" - name: %v\n", writer.Name)
	}

	isReadOnly := func(x *Account) bool {
		return !x.Writer
	}

	readers := slice.Filter(input, isReadOnly)
	fmt.Printf("Readers:\n")
	for _, reader := range readers {
		fmt.Printf(" - name: %v\n", reader.Name)
	}

	// Output:
	// Writers:
	//  - name: Alice
	// Readers:
	//  - name: Bob
	//  - name: Steve
}
