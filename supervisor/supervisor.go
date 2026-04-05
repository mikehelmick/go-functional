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

package supervisor

import (
	"context"
	"sync"
)

// Strategy controls how the supervisor responds to a child failure.
type Strategy int

const (
	// OneForOne restarts only the child that failed, leaving others running.
	OneForOne Strategy = iota
	// OneForAll restarts all children when any one of them fails.
	OneForAll
)

// ChildSpec defines a supervised worker. Start runs in its own goroutine.
// A nil return signals a clean exit (no restart). A non-nil return signals
// a crash and triggers a restart according to the supervisor's Strategy.
// Start must honour context cancellation to allow clean shutdown.
type ChildSpec struct {
	Name  string
	Start func(ctx context.Context) error
}

type exitMsg struct {
	idx int
	err error
}

// Supervisor manages a set of worker goroutines with automatic restart.
// All methods are safe to call from multiple goroutines.
// Methods must not be called after Stop.
type Supervisor struct {
	strategy Strategy
	specs    []ChildSpec
	stopCh   chan struct{}
	stopped  chan struct{}
	stopOnce sync.Once
}

// Start launches the supervisor and all children immediately.
func Start(strategy Strategy, specs []ChildSpec) *Supervisor {
	s := &Supervisor{
		strategy: strategy,
		specs:    specs,
		stopCh:   make(chan struct{}),
		stopped:  make(chan struct{}),
	}
	go s.loop()
	return s
}

func (s *Supervisor) loop() {
	defer close(s.stopped)

	exitCh := make(chan exitMsg, len(s.specs))
	ctx, cancel := s.startAll(exitCh)

	alive := len(s.specs)
	for alive > 0 {
		select {
		case <-s.stopCh:
			cancel()
			for alive > 0 {
				<-exitCh
				alive--
			}
			return
		case exit := <-exitCh:
			alive--
			if exit.err == nil {
				// Clean exit — do not restart.
				continue
			}
			// Crashed — check if we are already shutting down.
			select {
			case <-s.stopCh:
				cancel()
				for alive > 0 {
					<-exitCh
					alive--
				}
				return
			default:
			}
			switch s.strategy {
			case OneForOne:
				alive++
				go s.runChild(ctx, exit.idx, exitCh)
			case OneForAll:
				// Cancel context to stop all running children, drain their
				// exits, then restart every child with a fresh context.
				cancel()
				for alive > 0 {
					<-exitCh
					alive--
				}
				ctx, cancel = s.startAll(exitCh)
				alive = len(s.specs)
			}
		}
	}
	cancel()
}

// startAll creates a fresh context, launches all children, and returns the
// context and its cancel function. The caller is responsible for calling cancel.
func (s *Supervisor) startAll(exitCh chan exitMsg) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	for i := range s.specs {
		go s.runChild(ctx, i, exitCh)
	}
	return ctx, cancel
}

func (s *Supervisor) runChild(ctx context.Context, idx int, exitCh chan<- exitMsg) {
	err := s.specs[idx].Start(ctx)
	exitCh <- exitMsg{idx: idx, err: err}
}

// Stop cancels all children and blocks until the supervisor loop exits.
// Calling Stop more than once is safe; subsequent calls are no-ops.
func (s *Supervisor) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
		<-s.stopped
	})
}
