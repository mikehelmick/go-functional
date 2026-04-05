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

// MapValues returns a new map with the same keys as in, with each value
// transformed by fn.
func MapValues[K comparable, V any, R any](in map[K]V, fn func(V) R) map[K]R {
	rtn := make(map[K]R, len(in))
	for k, v := range in {
		rtn[k] = fn(v)
	}
	return rtn
}
