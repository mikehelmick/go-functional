[![Go Reference](https://pkg.go.dev/badge/github.com/mikehelmick/go-functional.svg)](https://pkg.go.dev/github.com/mikehelmick/go-functional)
![CI workflow](https://github.com/mikehelmick/go-functional/actions/workflows/ci.yaml/badge.svg)

# go-functional

Functional programming primitives and OTP-inspired concurrency patterns for Go generics.

Requires **Go 1.24+**. See the [documentation site](https://mikehelmick.github.io/go-functional/) for full API docs and examples.

## Installation

```bash
go get github.com/mikehelmick/go-functional
```

## Packages

### Functional primitives

| Package | Description |
|---|---|
| [`slice`](https://pkg.go.dev/github.com/mikehelmick/go-functional/slice) | Filter, Map, Fold, Find, Zip, Chunk, Sort, Dedup, Scan, and more |
| [`maps`](https://pkg.go.dev/github.com/mikehelmick/go-functional/maps) | MapValues, FilterMap, Merge, MergeWith, Invert, Keys, ToSlice |
| [`optional`](https://pkg.go.dev/github.com/mikehelmick/go-functional/optional) | `Maybe[T]` — Some or None with Map and GetOrElse |
| [`result`](https://pkg.go.dev/github.com/mikehelmick/go-functional/result) | `Result[T]` — wraps `(T, error)` with Map and chaining |
| [`pipeline`](https://pkg.go.dev/github.com/mikehelmick/go-functional/pipeline) | Function composition: `Pipe`, `Pipe2`, `Pipe3`, `Pipe4` |

### OTP-inspired concurrency

Inspired by Elixir's OTP abstractions. Each package provides a typed, goroutine-safe building block.

| Package | Description |
|---|---|
| [`agent`](https://pkg.go.dev/github.com/mikehelmick/go-functional/agent) | Goroutine-owned mutable state with `Get` / `Update` / `Cast` |
| [`genserver`](https://pkg.go.dev/github.com/mikehelmick/go-functional/genserver) | Generic request/response server loop — `Call` (sync) and `Cast` (async) |
| [`task`](https://pkg.go.dev/github.com/mikehelmick/go-functional/task) | Typed async work unit — `Run`, `Await`, `AwaitAll`, `Map` |
| [`supervisor`](https://pkg.go.dev/github.com/mikehelmick/go-functional/supervisor) | Managed goroutine restart with `OneForOne` and `OneForAll` strategies |

## Quick examples

### Slice operations

```go
nums := []int{1, 2, 3, 4, 5, 6}

evens  := slice.Filter(nums, func(n int) bool { return n%2 == 0 })  // [2 4 6]
doubled := slice.Map(nums, func(n int) int { return n * 2 })         // [2 4 6 8 10 12]
sum    := slice.FoldL(nums, 0, func(acc, n int) int { return acc + n }) // 21
top3   := slice.Take(slice.SortBy(nums, func(n int) int { return -n }), 3) // [6 5 4]
```

### Task — concurrent work

```go
tasks := []*task.Task[int]{
    task.Run(func() (int, error) { return expensiveA() }),
    task.Run(func() (int, error) { return expensiveB() }),
    task.Run(func() (int, error) { return expensiveC() }),
}
results, err := task.AwaitAll(tasks) // waits for all, returns in order
```

### Agent — shared mutable state

```go
counter := agent.New(0)
counter.Update(func(n int) int { return n + 1 })
fmt.Println(counter.Get()) // 1
counter.Stop()
```

### Supervisor — automatic restart

```go
sup := supervisor.Start(supervisor.OneForOne, []supervisor.ChildSpec{
    {
        Name: "worker",
        Start: func(ctx context.Context) error {
            // Runs until ctx is cancelled.
            // Returning a non-nil error triggers an automatic restart.
            return runWorker(ctx)
        },
    },
})
defer sup.Stop()
```

### GenServer — serialised state with call/cast

```go
type CounterServer struct{}

func (CounterServer) Init() int                                  { return 0 }
func (CounterServer) HandleCall(req string, n int) (int, int)    { return n, n } // return current
func (CounterServer) HandleCast(req string, n int) int           { return n + 1 } // increment

srv := genserver.Start[int, string, int](CounterServer{})
srv.Cast("inc")
srv.Cast("inc")
fmt.Println(srv.Call("get")) // 2
srv.Stop()
```

## Example application

See [`examples/wordfreq`](examples/wordfreq/main.go) for an OTP-style word frequency service that combines `supervisor`, `agent`, `task`, `pipeline`, `slice`, and `maps`.

## Attribution

Parts of this library were written with the assistance of [Claude Code](https://claude.ai/code).
