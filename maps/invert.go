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

package maps

// Invert returns a new map with keys and values swapped. If multiple keys
// in the input map to the same value, the resulting key will map to one of
// the original keys non-deterministically.
func Invert[K comparable, V comparable](in map[K]V) map[V]K {
	rtn := make(map[V]K, len(in))
	for k, v := range in {
		rtn[v] = k
	}
	return rtn
}
