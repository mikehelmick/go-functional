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

package slice

import (
	"cmp"
	"slices"
)

// SortBy returns a new slice sorted in ascending order by the key returned
// by keyFn. The original slice is not modified.
func SortBy[T any, K cmp.Ordered](in []T, keyFn func(T) K) []T {
	rtn := make([]T, len(in))
	copy(rtn, in)
	slices.SortFunc(rtn, func(a, b T) int {
		return cmp.Compare(keyFn(a), keyFn(b))
	})
	return rtn
}
