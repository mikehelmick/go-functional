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

import "maps"

// Merge combines two maps into a new map. When both maps contain the same
// key, the value from b takes precedence.
func Merge[K comparable, V any](a, b map[K]V) map[K]V {
	rtn := make(map[K]V, len(a)+len(b))
	maps.Copy(rtn, a)
	maps.Copy(rtn, b)
	return rtn
}

// MergeWith combines two maps into a new map. When both maps contain the
// same key, fn is called with the values from a and b to produce the
// merged value.
func MergeWith[K comparable, V any](a, b map[K]V, fn func(V, V) V) map[K]V {
	rtn := make(map[K]V, len(a)+len(b))
	maps.Copy(rtn, a)
	for k, v := range b {
		if existing, ok := rtn[k]; ok {
			rtn[k] = fn(existing, v)
		} else {
			rtn[k] = v
		}
	}
	return rtn
}
