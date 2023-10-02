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

type user struct {
	ID   int
	Name string
	LDAP string
}

func TestFrequencies(t *testing.T) {
	t.Parallel()

	users := []*user{
		{
			ID:   1,
			Name: "Tom",
			LDAP: "tom",
		},
		{
			ID:   2,
			Name: "Old Tom",
			LDAP: "old.tom",
		},
		{
			ID:   3,
			Name: "Really Old Tom",
			LDAP: "rot",
		},
	}

	asMap, err := maps.ToMap(users, func(u *user) int {
		return u.ID
	})
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	wantMap := map[int]*user{
		1: users[0],
		2: users[1],
		3: users[2],
	}

	if diff := cmp.Diff(wantMap, asMap); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}

	// convert it back to a slice
	gotSlice := maps.ToSlice(asMap)
	lessFunc := func(a *user, b *user) bool {
		return a.ID < b.ID
	}
	if diff := cmp.Diff(users, gotSlice, cmpopts.SortSlices(lessFunc)); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}
