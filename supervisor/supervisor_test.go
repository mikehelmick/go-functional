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

package supervisor_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mikehelmick/go-functional/supervisor"
)

var errCrash = errors.New("crash")

func TestOneForOneRestartsCrashedChild(t *testing.T) {
	t.Parallel()

	var count atomic.Int32
	restartedCh := make(chan struct{})

	specs := []supervisor.ChildSpec{
		{
			Name: "crasher",
			Start: func(ctx context.Context) error {
				if count.Add(1) == 1 {
					return errCrash
				}
				close(restartedCh)
				<-ctx.Done()
				return nil
			},
		},
	}

	sup := supervisor.Start(supervisor.OneForOne, specs)
	select {
	case <-restartedCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for restart")
	}
	sup.Stop()

	if count.Load() < 2 {
		t.Errorf("expected at least 2 starts, got %d", count.Load())
	}
}

func TestOneForOneOnlyRestartsCrashedChild(t *testing.T) {
	t.Parallel()

	var counts [2]atomic.Int32
	crasherRestartedCh := make(chan struct{})

	specs := []supervisor.ChildSpec{
		{
			Name: "crasher",
			Start: func(ctx context.Context) error {
				if counts[0].Add(1) == 1 {
					return errCrash
				}
				close(crasherRestartedCh)
				<-ctx.Done()
				return nil
			},
		},
		{
			Name: "stable",
			Start: func(ctx context.Context) error {
				counts[1].Add(1)
				<-ctx.Done()
				return nil
			},
		},
	}

	sup := supervisor.Start(supervisor.OneForOne, specs)
	select {
	case <-crasherRestartedCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for crasher restart")
	}
	sup.Stop()

	if counts[0].Load() < 2 {
		t.Errorf("crasher: expected at least 2 starts, got %d", counts[0].Load())
	}
	if counts[1].Load() != 1 {
		t.Errorf("stable: expected exactly 1 start, got %d", counts[1].Load())
	}
}

func TestOneForAllRestartsAllOnCrash(t *testing.T) {
	t.Parallel()

	var counts [2]atomic.Int32
	// closed when both children have been started a second time
	bothRestartedCh := make(chan struct{})
	var closeOnce atomic.Bool

	makeSpec := func(idx int) supervisor.ChildSpec {
		return supervisor.ChildSpec{
			Name: fmt.Sprintf("child-%d", idx),
			Start: func(ctx context.Context) error {
				n := counts[idx].Add(1)
				if idx == 0 && n == 1 {
					return errCrash // first child crashes on its first run
				}
				// Signal when both are on their second run.
				if counts[0].Load() >= 2 && counts[1].Load() >= 2 {
					if closeOnce.CompareAndSwap(false, true) {
						close(bothRestartedCh)
					}
				}
				<-ctx.Done()
				return nil
			},
		}
	}

	specs := []supervisor.ChildSpec{makeSpec(0), makeSpec(1)}
	sup := supervisor.Start(supervisor.OneForAll, specs)

	select {
	case <-bothRestartedCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for all children to restart")
	}
	sup.Stop()

	if counts[0].Load() < 2 {
		t.Errorf("child-0: expected at least 2 starts, got %d", counts[0].Load())
	}
	if counts[1].Load() < 2 {
		t.Errorf("child-1: expected at least 2 starts, got %d", counts[1].Load())
	}
}

func TestCleanExitNoRestart(t *testing.T) {
	t.Parallel()

	var count atomic.Int32
	specs := []supervisor.ChildSpec{
		{
			Name: "one-shot",
			Start: func(_ context.Context) error {
				count.Add(1)
				return nil // clean exit — should not restart
			},
		},
	}

	sup := supervisor.Start(supervisor.OneForOne, specs)
	// Allow time for any unexpected restart.
	time.Sleep(50 * time.Millisecond)
	sup.Stop()

	if count.Load() != 1 {
		t.Errorf("expected exactly 1 start, got %d", count.Load())
	}
}

func TestStop(t *testing.T) {
	t.Parallel()

	started := make(chan struct{}, 1)
	stopped := make(chan struct{})

	specs := []supervisor.ChildSpec{
		{
			Name: "blocker",
			Start: func(ctx context.Context) error {
				started <- struct{}{}
				<-ctx.Done()
				close(stopped)
				return nil
			},
		},
	}

	sup := supervisor.Start(supervisor.OneForOne, specs)
	<-started
	sup.Stop()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("child did not stop after supervisor Stop()")
	}
}

func TestStopIdempotent(t *testing.T) {
	t.Parallel()

	specs := []supervisor.ChildSpec{
		{
			Name: "noop",
			Start: func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			},
		},
	}

	sup := supervisor.Start(supervisor.OneForOne, specs)
	sup.Stop()
	sup.Stop() // must not panic or deadlock
}

func ExampleStart() {
	ready := make(chan struct{})
	specs := []supervisor.ChildSpec{
		{
			Name: "worker",
			Start: func(ctx context.Context) error {
				close(ready)
				<-ctx.Done()
				return nil
			},
		},
	}
	sup := supervisor.Start(supervisor.OneForOne, specs)
	<-ready
	sup.Stop()
	fmt.Println("stopped")
	// Output:
	// stopped
}
