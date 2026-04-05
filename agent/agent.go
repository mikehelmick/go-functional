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

package agent

import "sync"

// call is an internal message sent to the agent's loop goroutine.
type call[S any] struct {
	// fn transforms the current state and returns (newState, replyValue).
	fn func(S) (S, any)
	// replyCh receives the reply value. nil for async calls (Cast).
	replyCh chan any
	// stop signals the loop to exit.
	stop bool
}

// Agent holds state of type S in a dedicated goroutine. All operations are
// goroutine-safe and serialized. Methods must not be called after Stop.
type Agent[S any] struct {
	callCh   chan call[S]
	doneCh   chan struct{}
	stopOnce sync.Once
}

// New creates and starts an Agent with the given initial state.
func New[S any](initial S) *Agent[S] {
	a := &Agent[S]{
		callCh: make(chan call[S]),
		doneCh: make(chan struct{}),
	}
	go a.loop(initial)
	return a
}

func (a *Agent[S]) loop(state S) {
	defer close(a.doneCh)
	for c := range a.callCh {
		if c.stop {
			return
		}
		newState, reply := c.fn(state)
		state = newState
		if c.replyCh != nil {
			c.replyCh <- reply
		}
	}
}

// Get returns the current state.
func (a *Agent[S]) Get() S {
	replyCh := make(chan any, 1)
	a.callCh <- call[S]{
		fn:      func(s S) (S, any) { return s, s },
		replyCh: replyCh,
	}
	return (<-replyCh).(S) //nolint:forcetypeassert
}

// Update applies fn to the current state, stores the result, and blocks
// until the update is complete.
func (a *Agent[S]) Update(fn func(S) S) {
	replyCh := make(chan any, 1)
	a.callCh <- call[S]{
		fn:      func(s S) (S, any) { return fn(s), nil },
		replyCh: replyCh,
	}
	<-replyCh
}

// Cast applies fn to the current state asynchronously. It returns
// immediately without waiting for the update to complete.
func (a *Agent[S]) Cast(fn func(S) S) {
	a.callCh <- call[S]{
		fn: func(s S) (S, any) { return fn(s), nil },
	}
}

// Stop shuts down the agent. It blocks until the loop goroutine exits.
// Calling Stop more than once is safe; subsequent calls are no-ops.
func (a *Agent[S]) Stop() {
	a.stopOnce.Do(func() {
		a.callCh <- call[S]{stop: true}
		<-a.doneCh
	})
}

// GetWith applies fn to the agent's state and returns the result.
// Unlike Get, the return type R may differ from the state type S.
// This is a package-level function because Go methods cannot introduce
// additional type parameters.
func GetWith[S, R any](a *Agent[S], fn func(S) R) R {
	replyCh := make(chan any, 1)
	a.callCh <- call[S]{
		fn:      func(s S) (S, any) { return s, fn(s) },
		replyCh: replyCh,
	}
	return (<-replyCh).(R) //nolint:forcetypeassert
}
