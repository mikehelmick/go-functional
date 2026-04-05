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

// Command circuitbreaker demonstrates using genserver to implement a
// circuit breaker — a resilience pattern that stops cascading failures
// by fast-failing requests when a downstream dependency is unhealthy.
//
// State machine:
//
//	Closed ──(threshold failures)──▶ Open ──(timeout elapsed)──▶ HalfOpen
//	  ▲                                                               │
//	  └──────────────────(success)─────────────────────────────────── ┘
//	                     (failure returns to Open)
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/mikehelmick/go-functional/genserver"
)

// ── Mode ────────────────────────────────────────────────────────────────────

// Mode is the circuit breaker's operating mode.
type Mode int

const (
	Closed   Mode = iota // normal: all requests pass through
	Open                 // tripped: requests are rejected immediately
	HalfOpen             // probing: one test request is allowed through
)

func (m Mode) String() string {
	switch m {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// ── Config ───────────────────────────────────────────────────────────────────

// Config holds circuit breaker tuning parameters.
type Config struct {
	// Threshold is the number of consecutive failures that trip the breaker.
	Threshold int
	// Timeout is how long the breaker stays Open before probing with a
	// single test request (transitioning to HalfOpen).
	Timeout time.Duration
}

// ── GenServer state and request types ────────────────────────────────────────

type cbState struct {
	mode     Mode
	failures int
	openedAt time.Time
	cfg      Config
}

type cbOp int

const (
	opAllow   cbOp = iota // Call: is the next request permitted?
	opSuccess             // Cast: the last request succeeded
	opFailure             // Cast: the last request failed
)

// ── Server implementation ────────────────────────────────────────────────────

type cbServer struct{ cfg Config }

func (s cbServer) Init() cbState {
	return cbState{mode: Closed, cfg: s.cfg}
}

// HandleCall handles the synchronous opAllow request. It returns whether the
// caller may proceed and the (possibly updated) state.
func (cbServer) HandleCall(req cbOp, st cbState) (bool, cbState) {
	switch st.mode {
	case Closed:
		return true, st
	case Open:
		// Transition to HalfOpen once the timeout has elapsed.
		if time.Since(st.openedAt) >= st.cfg.Timeout {
			st.mode = HalfOpen
			return true, st
		}
		return false, st
	case HalfOpen:
		return true, st
	}
	return false, st
}

// HandleCast handles async outcome reports (opSuccess / opFailure).
func (cbServer) HandleCast(req cbOp, st cbState) cbState {
	switch req {
	case opSuccess:
		if st.mode == HalfOpen {
			// Probe succeeded — close the circuit and reset the counter.
			st.mode = Closed
			st.failures = 0
		}
	case opFailure:
		st.failures++
		if st.mode == HalfOpen || (st.mode == Closed && st.failures >= st.cfg.Threshold) {
			st.mode = Open
			st.openedAt = time.Now()
		}
	}
	return st
}

// ── Public API ───────────────────────────────────────────────────────────────

// ErrOpen is returned by Do when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker open")

// CircuitBreaker wraps a genserver to provide a goroutine-safe circuit breaker.
// All state transitions are serialised through the genserver loop.
type CircuitBreaker struct {
	srv *genserver.GenServer[cbState, cbOp, bool]
}

// New creates and starts a CircuitBreaker with the given Config.
func New(cfg Config) *CircuitBreaker {
	return &CircuitBreaker{
		srv: genserver.Start[cbState, cbOp, bool](cbServer{cfg: cfg}),
	}
}

// Do executes fn if the circuit allows it, records the outcome, and returns
// ErrOpen if the request was rejected because the circuit is open.
func (cb *CircuitBreaker) Do(fn func() error) error {
	if !cb.srv.Call(opAllow) {
		return ErrOpen
	}
	if err := fn(); err != nil {
		cb.srv.Cast(opFailure)
		return err
	}
	cb.srv.Cast(opSuccess)
	return nil
}

// Stop shuts down the underlying genserver.
func (cb *CircuitBreaker) Stop() {
	cb.srv.Stop()
}

// ── Demo ─────────────────────────────────────────────────────────────────────

func main() {
	cb := New(Config{Threshold: 3, Timeout: 200 * time.Millisecond})
	defer cb.Stop()

	errFlaky := errors.New("service unavailable")

	report := func(label string, err error) {
		switch {
		case err == nil:
			fmt.Printf("  %-28s ✓ ok\n", label)
		case errors.Is(err, ErrOpen):
			fmt.Printf("  %-28s ✗ rejected (circuit open)\n", label)
		default:
			fmt.Printf("  %-28s ✗ failed: %v\n", label, err)
		}
	}

	fmt.Println("Phase 1 — normal operation (circuit closed)")
	report("request", cb.Do(func() error { return nil }))

	fmt.Println("\nPhase 2 — three consecutive failures trip the breaker")
	for i := range 3 {
		report(fmt.Sprintf("request %d", i+1), cb.Do(func() error { return errFlaky }))
	}

	fmt.Println("\nPhase 3 — circuit open, requests are rejected immediately")
	for i := range 3 {
		report(fmt.Sprintf("request %d", i+1), cb.Do(func() error { return nil }))
	}

	fmt.Printf("\nPhase 4 — waiting %s for recovery timeout...\n", 200*time.Millisecond)
	time.Sleep(250 * time.Millisecond)

	fmt.Println("\nPhase 5 — half-open: probe succeeds, circuit closes")
	report("probe", cb.Do(func() error { return nil }))

	fmt.Println("\nPhase 6 — normal operation resumed")
	for i := range 3 {
		report(fmt.Sprintf("request %d", i+1), cb.Do(func() error { return nil }))
	}
}
