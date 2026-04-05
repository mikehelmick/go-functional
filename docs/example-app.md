---
layout: default
title: Example Applications
nav_order: 4
render_with_liquid: false
---

# Example Applications
{: .no_toc }

<details open markdown="block">
  <summary>Contents</summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

---

## Word Frequency Service

[View source on GitHub](https://github.com/mikehelmick/go-functional/blob/main/examples/wordfreq/main.go){: .btn .btn-outline }

---

`examples/wordfreq` is an OTP-style word frequency service. It is a self-contained runnable program that demonstrates all major packages working together.

```bash
go run github.com/mikehelmick/go-functional/examples/wordfreq@latest
```

## Architecture

```
supervisor
  └── reporter goroutine   (prints progress every 5 ms)

agent[map[string]int]      (accumulates merged word counts)

task × N                   (one per document, run concurrently)
  └── pipeline.Pipe2       (strings.Fields → slice.Map(ToLower))
  └── slice.Frequencies    (counts words in the document)

slice.SortBy / Take / Filter   (post-process final results)
maps.Keys                       (extract word list)
```

## Walkthrough

### 1 — Start the supervisor

A `supervisor.OneForOne` supervisor keeps a progress-reporter goroutine alive for the duration of the program. If it crashes for any reason, the supervisor restarts it automatically.

```go
sup := supervisor.Start(supervisor.OneForOne, []supervisor.ChildSpec{
    {
        Name: "reporter",
        Start: func(ctx context.Context) error {
            for {
                n := agent.GetWith(freqs, func(m map[string]int) int { return len(m) })
                fmt.Printf("  [reporter] unique words indexed: %d\n", n)
                select {
                case <-ctx.Done():
                    return nil
                case <-time.After(5 * time.Millisecond):
                }
            }
        },
    },
})
defer sup.Stop()
```

### 2 — Process documents concurrently with tasks

Each document is processed in its own goroutine via `task.Run`. The text pipeline is built with `pipeline.Pipe2` and `slice.Frequencies`.

```go
var tokenise = pipeline.Pipe2(
    strings.Fields,
    func(words []string) []string {
        return slice.Map(words, strings.ToLower)
    },
)

tasks := slice.Map(corpus, func(doc string) *task.Task[map[string]int] {
    return task.Run(func() (map[string]int, error) {
        return slice.Frequencies(tokenise(doc)), nil
    })
})

results, _ := task.AwaitAll(tasks)
```

### 3 — Merge into the agent

Each document's frequency map is merged into the shared `agent` serially. The agent's goroutine guarantees no data races without any explicit locking.

```go
for _, r := range results {
    freqs.Update(func(cur map[string]int) map[string]int {
        for word, count := range r {
            cur[word] += count
        }
        return cur
    })
}
```

### 4 — Query and display results

Once all work is done, `slice.SortBy`, `slice.Take`, and `slice.Filter` transform the final map into a ranked list.

```go
type wordCount struct {
    word  string
    count int
}

entries := slice.Map(maps.Keys(final), func(w string) wordCount {
    return wordCount{w, final[w]}
})
sorted  := slice.SortBy(entries, func(e wordCount) int { return -e.count })

for _, e := range slice.Take(sorted, 5) {
    fmt.Printf("  %-14s %d\n", e.word, e.count)
}

repeated := slice.Filter(sorted, func(e wordCount) bool { return e.count > 1 })
fmt.Printf("%d words appear more than once\n", len(repeated))
```

## Sample output

```
  [reporter] unique words indexed: 0

Top 5 words:
  the            4
  programming    4
  to             3
  lazy           3
  and            3

14 of 46 unique words appear more than once
```

---

## Circuit Breaker

[View source on GitHub](https://github.com/mikehelmick/go-functional/blob/main/examples/circuitbreaker/main.go){: .btn .btn-outline }

`examples/circuitbreaker` implements a circuit breaker using `genserver` as the state machine. A circuit breaker stops cascading failures by fast-failing requests when a downstream dependency is unhealthy.

```bash
go run github.com/mikehelmick/go-functional/examples/circuitbreaker@latest
```

### State machine

```
Closed ──(threshold failures)──▶ Open ──(timeout elapsed)──▶ HalfOpen
  ▲                                                               │
  └──────────────────(success)────────────────────────────────── ┘
                     (failure returns to Open)
```

| Mode | Behaviour |
|---|---|
| **Closed** | Normal operation — all requests pass through |
| **Open** | Tripped — requests are rejected immediately (fast-fail) |
| **HalfOpen** | Probing — one test request is allowed; success closes, failure reopens |

### Why genserver fits perfectly

The circuit breaker's state must be mutated atomically: reading the current mode and transitioning it must happen as a single unit. `genserver` provides this guarantee for free — `HandleCall` and `HandleCast` run serially in a single goroutine, so no mutex is needed.

- **`Call(opAllow)`** — synchronous gate check; reads the mode and may transition `Open → HalfOpen`
- **`Cast(opSuccess)` / `Cast(opFailure)`** — async outcome reports; update the failure counter and drive state transitions

### Implementation

```go
type cbServer struct{ cfg Config }

func (cbServer) HandleCall(req cbOp, st cbState) (bool, cbState) {
    switch st.mode {
    case Closed:
        return true, st
    case Open:
        if time.Since(st.openedAt) >= st.cfg.Timeout {
            st.mode = HalfOpen   // transition on first probe attempt
            return true, st
        }
        return false, st
    case HalfOpen:
        return true, st
    }
    return false, st
}

func (cbServer) HandleCast(req cbOp, st cbState) cbState {
    switch req {
    case opSuccess:
        if st.mode == HalfOpen {
            st.mode, st.failures = Closed, 0
        }
    case opFailure:
        st.failures++
        if st.mode == HalfOpen || (st.mode == Closed && st.failures >= st.cfg.Threshold) {
            st.mode, st.openedAt = Open, time.Now()
        }
    }
    return st
}
```

The public `Do` method wraps the genserver calls into a clean API:

```go
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
```

### Sample output

```
Phase 1 — normal operation (circuit closed)
  request                      ✓ ok

Phase 2 — three consecutive failures trip the breaker
  request 1                    ✗ failed: service unavailable
  request 2                    ✗ failed: service unavailable
  request 3                    ✗ failed: service unavailable

Phase 3 — circuit open, requests are rejected immediately
  request 1                    ✗ rejected (circuit open)
  request 2                    ✗ rejected (circuit open)
  request 3                    ✗ rejected (circuit open)

Phase 4 — waiting 200ms for recovery timeout...

Phase 5 — half-open: probe succeeds, circuit closes
  probe                        ✓ ok

Phase 6 — normal operation resumed
  request 1                    ✓ ok
  request 2                    ✓ ok
  request 3                    ✓ ok
```
