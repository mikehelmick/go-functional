---
layout: default
title: go-functional
---

# go-functional

Functional programming primitives and OTP-inspired concurrency patterns for Go generics.

[![Go Reference](https://pkg.go.dev/badge/github.com/mikehelmick/go-functional.svg)](https://pkg.go.dev/github.com/mikehelmick/go-functional)
![CI](https://github.com/mikehelmick/go-functional/actions/workflows/ci.yaml/badge.svg)

**Requires Go 1.24+**

```bash
go get github.com/mikehelmick/go-functional
```

---

## Functional primitives

### slice

Operations on typed slices. All functions are generic and allocate a new slice rather than mutating the input.

| Function | Description |
|---|---|
| `Filter` | Return elements matching a predicate |
| `Reject` | Return elements not matching a predicate |
| `Partition` | Split into matching and non-matching in one pass |
| `Map` | Transform each element |
| `FlatMap` / `Flatten` | Map then flatten nested slices |
| `FoldL` / `FoldR` | Reduce from left or right |
| `Scan` | FoldL variant returning all intermediate accumulators |
| `Find` / `FindIndex` | First element matching a predicate |
| `All` / `Any` | Universal and existential quantifiers |
| `Take` / `TakeWhile` / `TakeEvery` | Prefix and stride operations |
| `SortBy` | Sort by a key function (ascending) |
| `MinBy` / `MaxBy` | Extremes by a key function |
| `Zip` / `Unzip` | Pair and unpair two slices |
| `ChunkEvery` / `ChunkBy` | Split into fixed-size or key-grouped chunks |
| `Dedup` / `DedupBy` | Remove consecutive duplicates |
| `Frequencies` / `FrequenciesBy` | Count occurrences |
| `GroupBy` | Group elements by a key function |

```go
nums := []int{1, 2, 3, 4, 5, 6}

evens   := slice.Filter(nums, func(n int) bool { return n%2 == 0 }) // [2 4 6]
doubled := slice.Map(nums, func(n int) int { return n * 2 })         // [2 4 6 8 10 12]
sum     := slice.FoldL(nums, 0, func(acc, n int) int { return acc + n }) // 21
top3    := slice.Take(slice.SortBy(nums, func(n int) int { return -n }), 3) // [6 5 4]
freqs   := slice.Frequencies([]string{"a", "b", "a", "c", "a", "b"})
// map[a:3 b:2 c:1]
```

---

### maps

Operations on `map[K]V` where `K` is `comparable`.

| Function | Description |
|---|---|
| `Keys` | Extract keys as a slice |
| `ToSlice` | Convert map values to a slice via a function |
| `ToMap` | Convert a slice to a map via a key function |
| `MapValues` | Transform all values, keeping keys |
| `FilterMap` | Remove entries whose value fails a predicate |
| `Merge` | Combine two maps; right-hand values win on conflict |
| `MergeWith` | Combine two maps with a custom conflict resolver |
| `Invert` | Swap keys and values |

```go
scores := map[string]int{"alice": 10, "bob": 7, "carol": 10}

doubled := maps.MapValues(scores, func(v int) int { return v * 2 })
passing := maps.FilterMap(scores, func(v int) bool { return v >= 8 })
inverted := maps.Invert(map[string]string{"en": "English", "fr": "French"})
```

---

### optional

`Maybe[T]` represents a value that may or may not be present — `Some(v)` or `None[T]()`.

```go
m := optional.Some(42)
doubled := optional.Map(m, func(v int) int { return v * 2 }) // Some(84)
val := m.GetOrElse(0)                                         // 42

n := optional.None[int]()
val = n.GetOrElse(-1) // -1
```

---

### result

`Result[T]` wraps `(T, error)` for chainable error handling.

```go
r := result.Of(strconv.Atoi("42"))        // Ok(42)
doubled := result.Map(r, func(v int) int { return v * 2 }) // Ok(84)

bad := result.Of(strconv.Atoi("oops"))   // Err(...)
val, err := bad.Unwrap()                  // val=0, err=...
```

---

### pipeline

Function composition left-to-right. `Pipe` composes same-type functions; `Pipe2`–`Pipe4` handle type-changing chains.

```go
// Same-type composition
process := pipeline.Pipe(
    strings.TrimSpace,
    strings.ToLower,
)
fmt.Println(process("  Hello World  ")) // "hello world"

// Type-changing chain
countWords := pipeline.Pipe2(
    strings.Fields,                          // string → []string
    func(ws []string) int { return len(ws) }, // []string → int
)
fmt.Println(countWords("one two three")) // 3
```

---

## OTP-inspired concurrency

These packages bring Elixir's OTP patterns to Go: typed goroutines with well-defined lifecycles, serialised state, and supervised restart.

### agent

`Agent[S]` owns state of type `S` in a dedicated goroutine. All access is serialised and goroutine-safe.

```go
counter := agent.New(0)

counter.Update(func(n int) int { return n + 1 })
counter.Cast(func(n int) int { return n + 1 }) // async, no wait

fmt.Println(counter.Get()) // 2

// GetWith extracts a projection without exposing the full state
length := agent.GetWith(myMapAgent, func(m map[string]int) int { return len(m) })

counter.Stop()
```

---

### genserver

`GenServer[S, Req, Resp]` runs a single-goroutine state machine. `Call` is synchronous (blocks for a response); `Cast` is async.

```go
type CounterServer struct{}

func (CounterServer) Init() int                               { return 0 }
func (CounterServer) HandleCall(req string, n int) (int, int) { return n, n }
func (CounterServer) HandleCast(req string, n int) int        { return n + 1 }

srv := genserver.Start[int, string, int](CounterServer{})
srv.Cast("inc")
srv.Cast("inc")
fmt.Println(srv.Call("get")) // 2
srv.Stop()
```

---

### task

`Task[T]` runs a function in a goroutine and lets any number of callers `Await` the result. Safe to await from multiple goroutines; result is computed exactly once.

```go
// Fire and collect
tasks := []*task.Task[int]{
    task.Run(func() (int, error) { return fetchA() }),
    task.Run(func() (int, error) { return fetchB() }),
}
results, err := task.AwaitAll(tasks) // [resultA, resultB]

// Transform the result type
strTask := task.Map(tasks[0], func(v int) string { return fmt.Sprint(v) })
s, _ := strTask.Await()
```

---

### supervisor

`Supervisor` manages a set of goroutines and restarts them on failure according to a `Strategy`.

- **`OneForOne`** — restart only the failed child; others keep running.
- **`OneForAll`** — when any child fails, stop all and restart the full set.

A `ChildSpec.Start` function signals:
- `nil` return → clean exit, do not restart.
- non-nil return → crash, apply restart strategy.
- `ctx.Done()` → supervisor is stopping, return `nil`.

```go
sup := supervisor.Start(supervisor.OneForOne, []supervisor.ChildSpec{
    {
        Name: "fetcher",
        Start: func(ctx context.Context) error {
            return runFetcher(ctx) // restarted if this returns non-nil
        },
    },
    {
        Name: "processor",
        Start: func(ctx context.Context) error {
            return runProcessor(ctx)
        },
    },
})
defer sup.Stop()
```

---

## Example application

[`examples/wordfreq`](https://github.com/mikehelmick/go-functional/blob/main/examples/wordfreq/main.go) shows all packages working together in an OTP-style word-frequency service:

- `supervisor` manages a progress-reporter goroutine
- `agent` accumulates word frequencies from concurrent workers
- `task` processes each document in a separate goroutine
- `pipeline` composes the text-normalisation steps
- `slice` and `maps` transform and query the final results

---

## Attribution

Parts of this library were written with the assistance of [Claude Code](https://claude.ai/code).
