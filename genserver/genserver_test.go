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

package genserver_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikehelmick/go-functional/genserver"
)

// --- counter server used across tests ---

type counterCmd int

const (
	cmdGet counterCmd = iota
	cmdIncrement
	cmdAdd
	cmdReset
)

type counterReq struct {
	cmd counterCmd
	val int
}

// counterServer is a simple integer counter implementing genserver.Server.
type counterServer struct{}

func (counterServer) Init() int { return 0 }

func (counterServer) HandleCall(req counterReq, state int) (int, int) {
	switch req.cmd {
	case cmdGet:
		return state, state
	case cmdIncrement:
		return state + 1, state + 1
	case cmdAdd, cmdReset:
		return state, state
	}
	return state, state
}

func (counterServer) HandleCast(req counterReq, state int) int {
	switch req.cmd {
	case cmdAdd:
		return state + req.val
	case cmdReset:
		return 0
	case cmdGet, cmdIncrement:
		return state
	}
	return state
}

// --- tests ---

func TestCallGet(t *testing.T) {
	t.Parallel()
	gs := genserver.Start[int, counterReq, int](counterServer{})
	defer gs.Stop()

	got := gs.Call(counterReq{cmd: cmdGet})
	if diff := cmp.Diff(0, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestCallIncrement(t *testing.T) {
	t.Parallel()
	gs := genserver.Start[int, counterReq, int](counterServer{})
	defer gs.Stop()

	gs.Call(counterReq{cmd: cmdIncrement})
	gs.Call(counterReq{cmd: cmdIncrement})
	got := gs.Call(counterReq{cmd: cmdGet})
	if diff := cmp.Diff(2, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestCastAdd(t *testing.T) {
	t.Parallel()
	gs := genserver.Start[int, counterReq, int](counterServer{})
	defer gs.Stop()

	gs.Cast(counterReq{cmd: cmdAdd, val: 5})
	gs.Cast(counterReq{cmd: cmdAdd, val: 3})
	// Call serializes after the casts.
	got := gs.Call(counterReq{cmd: cmdGet})
	if diff := cmp.Diff(8, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestCastReset(t *testing.T) {
	t.Parallel()
	gs := genserver.Start[int, counterReq, int](counterServer{})
	defer gs.Stop()

	gs.Call(counterReq{cmd: cmdIncrement})
	gs.Call(counterReq{cmd: cmdIncrement})
	gs.Cast(counterReq{cmd: cmdReset})
	got := gs.Call(counterReq{cmd: cmdGet})
	if diff := cmp.Diff(0, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestStopIsIdempotent(t *testing.T) {
	t.Parallel()
	gs := genserver.Start[int, counterReq, int](counterServer{})
	gs.Stop()
	gs.Stop()
}

func TestConcurrentCalls(t *testing.T) {
	t.Parallel()
	gs := genserver.Start[int, counterReq, int](counterServer{})
	defer gs.Stop()

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gs.Call(counterReq{cmd: cmdIncrement})
		}()
	}
	wg.Wait()

	got := gs.Call(counterReq{cmd: cmdGet})
	if diff := cmp.Diff(100, got); diff != "" {
		t.Fatalf("mismatch (-want, +got):\n%s", diff)
	}
}

// --- example ---

func ExampleGenServer_Call() {
	gs := genserver.Start[int, counterReq, int](counterServer{})
	defer gs.Stop()

	gs.Cast(counterReq{cmd: cmdAdd, val: 10})
	gs.Call(counterReq{cmd: cmdIncrement})
	fmt.Println(gs.Call(counterReq{cmd: cmdGet}))
	// Output:
	// 11
}
