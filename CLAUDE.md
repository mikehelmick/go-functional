# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run tests (short mode, shuffled)
make test

# Run tests with race detector and coverage
make test-acc

# View coverage report (after test-acc)
make test-coverage

# Lint
make lint

# Run a single test
go test -run TestName ./slice/
go test -run TestName ./maps/
```

## Architecture

This is a Go generics library providing functional programming primitives. It requires Go 1.24+.

Two packages:

**`slice/`** — operations on slices:
- `Map[T, R]` / `MapToPtr[T]` — transform elements
- `Filter[T]` — select elements by predicate (`MatchFn[T]`)
- `FoldL[T, A]` / `FoldR[T, A]` — reduce from left/right
- `Take[T]` — take first N elements
- `Group[T]` — group elements by key
- `Frequencies[T]` — count occurrences
- `matchers.go` defines shared `MatchFn[T]` type

**`maps/`** — operations on maps:
- `ToMap[K, T]` — convert slice to map via key function (errors on duplicates)
- `ToSlice[K, T]` — convert map values to slice
- `Keys[K, T]` — extract map keys as slice

Each function is generic and uses type constraints (`any`, `comparable`). Tests use `github.com/google/go-cmp/cmp` for comparisons.
