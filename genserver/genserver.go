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

package genserver

import "sync"

// Server defines the callbacks for a GenServer process.
// S is the state type, Req is the request/message type, and Resp is the
// response type returned by synchronous calls.
type Server[S, Req, Resp any] interface {
	// Init returns the initial state of the server.
	Init() S
	// HandleCall processes a synchronous request. It returns a response
	// sent back to the caller and the updated state.
	HandleCall(req Req, state S) (Resp, S)
	// HandleCast processes an asynchronous message. It returns the
	// updated state. No response is sent to the caller.
	HandleCast(msg Req, state S) S
}

// msg is an internal envelope used by the server loop.
type msg[Req, Resp any] struct {
	req     Req
	replyCh chan Resp // nil for Cast; non-nil for Call
	stop    bool
}

// GenServer is a running server process. All operations are goroutine-safe.
// Methods must not be called after Stop.
type GenServer[S, Req, Resp any] struct {
	msgCh    chan msg[Req, Resp]
	doneCh   chan struct{}
	stopOnce sync.Once
}

// Start initialises the server by calling srv.Init() and launches the
// message loop in a new goroutine.
func Start[S, Req, Resp any](srv Server[S, Req, Resp]) *GenServer[S, Req, Resp] {
	g := &GenServer[S, Req, Resp]{
		msgCh:  make(chan msg[Req, Resp]),
		doneCh: make(chan struct{}),
	}
	go g.loop(srv, srv.Init())
	return g
}

func (g *GenServer[S, Req, Resp]) loop(srv Server[S, Req, Resp], state S) {
	defer close(g.doneCh)
	for m := range g.msgCh {
		if m.stop {
			return
		}
		if m.replyCh != nil {
			resp, newState := srv.HandleCall(m.req, state)
			state = newState
			m.replyCh <- resp
		} else {
			state = srv.HandleCast(m.req, state)
		}
	}
}

// Call sends a synchronous request and blocks until the server returns a
// response.
func (g *GenServer[S, Req, Resp]) Call(req Req) Resp {
	replyCh := make(chan Resp, 1)
	g.msgCh <- msg[Req, Resp]{req: req, replyCh: replyCh}
	return <-replyCh
}

// Cast sends an asynchronous message and returns immediately without
// waiting for the server to process it.
func (g *GenServer[S, Req, Resp]) Cast(req Req) {
	g.msgCh <- msg[Req, Resp]{req: req}
}

// Stop shuts down the server and blocks until the loop goroutine exits.
// Calling Stop more than once is safe; subsequent calls are no-ops.
func (g *GenServer[S, Req, Resp]) Stop() {
	g.stopOnce.Do(func() {
		g.msgCh <- msg[Req, Resp]{stop: true}
		<-g.doneCh
	})
}
