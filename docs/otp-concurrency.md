---
layout: default
title: OTP Concurrency
nav_order: 3
render_with_liquid: false
---

# OTP-Inspired Concurrency
{: .no_toc }

These packages bring Elixir's OTP patterns to Go: typed goroutines with well-defined lifecycles,
serialised state ownership, and supervised restart.

<details open markdown="block">
  <summary>Contents</summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

---

## agent

`Agent[S]` owns state of type `S` in a dedicated goroutine. All reads and writes are serialised without requiring a mutex.

```go
counter := agent.New(0)

// Synchronous update — blocks until applied
counter.Update(func(n int) int { return n + 1 })

// Asynchronous update — returns immediately
counter.Cast(func(n int) int { return n + 1 })

// Read
fmt.Println(counter.Get()) // 2

counter.Stop()
```

When the result type differs from the state type, use the package-level `GetWith`:

```go
wordFreqs := agent.New(map[string]int{})

// Extract a projection without exposing the full map
uniqueCount := agent.GetWith(wordFreqs, func(m map[string]int) int {
    return len(m)
})
```

{: .note }
`Cast` is fire-and-forget. If ordering matters, use `Update`.

---

## genserver

`GenServer[S, Req, Resp]` runs a single-goroutine state machine. Implement the `Server` interface and start it with `genserver.Start`.

- **`Call`** — synchronous; blocks the caller until `HandleCall` returns a response.
- **`Cast`** — asynchronous; the caller returns immediately after enqueueing the message.

```go
type counterServer struct{}

func (counterServer) Init() int                               { return 0 }
func (counterServer) HandleCall(req string, n int) (int, int) { return n, n } // get
func (counterServer) HandleCast(req string, n int) int        { return n + 1 } // inc

srv := genserver.Start[int, string, int](counterServer{})
defer srv.Stop()

srv.Cast("inc")
srv.Cast("inc")
fmt.Println(srv.Call("get")) // 2
```

Because `GenServer` is a generic type, request and response types are fixed at construction. For servers that handle multiple message kinds, use a tagged struct or interface as the request type.

```go
type Req struct {
    Op    string
    Value int
}

type stackServer struct{}

func (stackServer) Init() []int { return nil }
func (stackServer) HandleCall(r Req, s []int) (int, []int) {
    if len(s) == 0 { return 0, s }
    return s[0], s
}
func (stackServer) HandleCast(r Req, s []int) []int {
    if r.Op == "push" { return append([]int{r.Value}, s...) }
    if r.Op == "pop" && len(s) > 0 { return s[1:] }
    return s
}
```

---

## task

`Task[T]` runs a function in a goroutine and provides a future-like handle. `Await` blocks until the result is ready. It is safe to call `Await` from multiple goroutines — the result is computed exactly once.

```go
// Fire three operations concurrently
tasks := []*task.Task[int]{
    task.Run(func() (int, error) { return fetchA() }),
    task.Run(func() (int, error) { return fetchB() }),
    task.Run(func() (int, error) { return fetchC() }),
}

// Collect all results in order; returns the first error encountered
results, err := task.AwaitAll(tasks)
```

`MustAwait` panics on error — useful in initialisation code where failure is not recoverable.

```go
config := task.Run(loadConfig).MustAwait()
```

`Map` transforms the result type without waiting:

```go
numTask  := task.Run(func() (int, error) { return 21, nil })
strTask  := task.Map(numTask, func(v int) string { return fmt.Sprint(v * 2) })
s, _     := strTask.Await() // "42"
```

---

## supervisor

`Supervisor` manages a set of goroutines and restarts them on failure. Workers are defined by a `ChildSpec`.

### ChildSpec contract

| Return value | Meaning |
|---|---|
| `nil` | Clean exit — do not restart |
| non-nil error | Crash — apply restart strategy |

Workers must honour context cancellation to allow clean shutdown:

```go
Start: func(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return nil // clean exit when supervisor stops
        default:
            if err := doWork(); err != nil {
                return err // triggers restart
            }
        }
    }
},
```

### Strategies

**`OneForOne`** — only the crashed child is restarted; others continue running.

```go
sup := supervisor.Start(supervisor.OneForOne, []supervisor.ChildSpec{
    {Name: "fetcher",   Start: runFetcher},
    {Name: "processor", Start: runProcessor},
})
defer sup.Stop()
```

**`OneForAll`** — when any child crashes, all children are stopped and the entire set is restarted together. Use when children share state that must stay consistent.

```go
sup := supervisor.Start(supervisor.OneForAll, []supervisor.ChildSpec{
    {Name: "producer", Start: runProducer},
    {Name: "consumer", Start: runConsumer},
})
```

### Shutdown

`Stop` cancels all children's contexts and blocks until every goroutine has exited. Calling `Stop` more than once is safe.

```go
sup.Stop() // idempotent
```
