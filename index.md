---
layout: default
title: Home
nav_order: 1
description: "Functional programming primitives and OTP-inspired concurrency patterns for Go generics."
permalink: /
---

# go-functional
{: .fs-9 }

Functional programming primitives and OTP-inspired concurrency patterns for Go generics.
{: .fs-6 .fw-300 }

[![Go Reference](https://pkg.go.dev/badge/github.com/mikehelmick/go-functional.svg)](https://pkg.go.dev/github.com/mikehelmick/go-functional)
![CI](https://github.com/mikehelmick/go-functional/actions/workflows/ci.yaml/badge.svg)

[Get started](#installation){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/mikehelmick/go-functional){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## Installation

```bash
go get github.com/mikehelmick/go-functional
```

Requires **Go 1.24+**.

---

## Packages at a glance

### Functional primitives

| Package | Highlights |
|---|---|
| [`slice`]({% link docs/functional-primitives.md %}#slice) | Filter, Map, Fold, Find, Zip, Chunk, Sort, Dedup, Scan, Frequencies |
| [`maps`]({% link docs/functional-primitives.md %}#maps) | MapValues, FilterMap, Merge, MergeWith, Invert |
| [`optional`]({% link docs/functional-primitives.md %}#optional) | `Maybe[T]` — Some or None |
| [`result`]({% link docs/functional-primitives.md %}#result) | `Result[T]` — chainable error handling |
| [`pipeline`]({% link docs/functional-primitives.md %}#pipeline) | Function composition: `Pipe`, `Pipe2`, `Pipe3`, `Pipe4` |

### OTP-inspired concurrency

| Package | Highlights |
|---|---|
| [`agent`]({% link docs/otp-concurrency.md %}#agent) | Goroutine-owned mutable state |
| [`genserver`]({% link docs/otp-concurrency.md %}#genserver) | Serialised request/response loop |
| [`task`]({% link docs/otp-concurrency.md %}#task) | Typed async work with `Await` / `AwaitAll` |
| [`supervisor`]({% link docs/otp-concurrency.md %}#supervisor) | Automatic goroutine restart (`OneForOne` / `OneForAll`) |

---

## Quick look

```go
// Concurrent document processing with OTP building blocks
freqs := agent.New(map[string]int{})

tasks := slice.Map(docs, func(doc string) *task.Task[map[string]int] {
    return task.Run(func() (map[string]int, error) {
        return slice.Frequencies(strings.Fields(doc)), nil
    })
})

results, _ := task.AwaitAll(tasks)
for _, r := range results {
    freqs.Update(func(cur map[string]int) map[string]int {
        for word, n := range r { cur[word] += n }
        return cur
    })
}
```

See the [example application]({% link docs/example-app.md %}) for the full OTP-style word frequency service.
