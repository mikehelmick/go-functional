---
layout: default
title: Functional Primitives
nav_order: 2
render_with_liquid: false
---

# Functional Primitives
{: .no_toc }

<details open markdown="block">
  <summary>Contents</summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

---

## slice

Operations on typed slices. Every function is generic and returns a new slice without mutating the input.

### Transforming

```go
nums := []int{1, 2, 3, 4, 5}

doubled := slice.Map(nums, func(n int) int { return n * 2 })
// [2, 4, 6, 8, 10]

words := slice.Map([]string{"hello", "world"}, strings.ToUpper)
// ["HELLO", "WORLD"]
```

`FlatMap` maps and flattens in one pass; `Flatten` flattens an already-nested slice.

```go
sentences := []string{"hello world", "foo bar"}
words := slice.FlatMap(sentences, strings.Fields)
// ["hello", "world", "foo", "bar"]
```

### Filtering

```go
evens, odds := slice.Partition(nums, func(n int) bool { return n%2 == 0 })
// evens=[2,4]  odds=[1,3,5]

evens = slice.Filter(nums, func(n int) bool { return n%2 == 0 })
odds  = slice.Reject(nums, func(n int) bool { return n%2 == 0 })
```

### Reducing

```go
sum  := slice.FoldL(nums, 0, func(acc, n int) int { return acc + n }) // 15
prod := slice.FoldR(nums, 1, func(n, acc int) int { return n * acc }) // 120

// Scan returns every intermediate accumulator.
running := slice.Scan(nums, 0, func(acc, n int) int { return acc + n })
// [1, 3, 6, 10, 15]
```

### Searching

```go
val, ok  := slice.Find(nums, func(n int) bool { return n > 3 })      // 4, true
idx, ok  := slice.FindIndex(nums, func(n int) bool { return n > 3 }) // 3, true
allBig   := slice.All(nums, func(n int) bool { return n > 0 })        // true
anyBig   := slice.Any(nums, func(n int) bool { return n > 4 })        // true
```

### Sorting and extremes

```go
// SortBy sorts ascending by the key function.
sorted := slice.SortBy(nums, func(n int) int { return -n }) // [5,4,3,2,1]

type Person struct {
    Name string
    Age  int
}
people := []Person{
    {Name: "Alice", Age: 30},
    {Name: "Bob",   Age: 25},
    {Name: "Carol", Age: 35},
}
youngest, _ := slice.MinBy(people, func(p Person) int { return p.Age }) // Bob
oldest,   _ := slice.MaxBy(people, func(p Person) int { return p.Age }) // Carol
```

### Taking and chunking

```go
slice.Take(nums, 3)              // [1, 2, 3]
slice.TakeWhile(nums, func(n int) bool { return n < 4 }) // [1, 2, 3]
slice.TakeEvery(nums, 2)        // [1, 3, 5]

slice.ChunkEvery(nums, 2)       // [[1,2],[3,4],[5]]
slice.ChunkBy(nums, func(n int) string {
    if n%2 == 0 { return "even" }
    return "odd"
}) // [[1],[2],[3],[4],[5]]
```

### Deduplication

```go
slice.Dedup([]int{1, 1, 2, 3, 3, 3, 2}) // [1, 2, 3, 2]  (consecutive only)

slice.DedupBy(people, func(p Person) int { return p.Age })
```

### Counting and grouping

```go
freqs := slice.Frequencies([]string{"a", "b", "a", "c", "a", "b"})
// map[a:3 b:2 c:1]

grouped := slice.GroupBy(people, func(p Person) string {
    if p.Age < 30 { return "young" }
    return "senior"
})
```

### Zipping

```go
pairs := slice.Zip([]int{1, 2, 3}, []string{"a", "b", "c"})
// [{1 a} {2 b} {3 c}]

ints, strs := slice.Unzip(pairs)
```

---

## maps

Operations on `map[K]V` where `K` is `comparable`.

```go
scores := map[string]int{"alice": 10, "bob": 7, "carol": 10}

// Transform values
doubled := maps.MapValues(scores, func(v int) int { return v * 2 })
// map[alice:20 bob:14 carol:20]

// Filter by value
passing := maps.FilterMap(scores, func(v int) bool { return v >= 8 })
// map[alice:10 carol:10]

// Merge — right-hand values win on conflict
a := map[string]int{"x": 1, "y": 2}
b := map[string]int{"y": 99, "z": 3}
merged := maps.Merge(a, b)
// map[x:1 y:99 z:3]

// Merge with custom conflict resolution
mergedWith := maps.MergeWith(a, b, func(va, vb int) int { return va + vb })
// map[x:1 y:101 z:3]

// Invert (values become keys)
inverted := maps.Invert(map[string]string{"en": "English", "fr": "French"})
// map[English:en French:fr]

// Keys / ToSlice / ToMap
keys := maps.Keys(scores)

byName, err := maps.ToMap(people, func(p Person) (string, error) {
    return p.Name, nil
})
```

---

## optional

`Maybe[T]` represents a value that may or may not be present.

```go
some := optional.Some(42)
none := optional.None[int]()

// Safe extraction
val := some.GetOrElse(0)  // 42
val  = none.GetOrElse(-1) // -1

// Transform — None propagates automatically
doubled := optional.Map(some, func(v int) int { return v * 2 }) // Some(84)
doubled  = optional.Map(none, func(v int) int { return v * 2 }) // None

// Unwrap to (T, bool)
v, ok := some.Unwrap() // 42, true
v, ok  = none.Unwrap() // 0, false
```

---

## result

`Result[T]` wraps `(T, error)` for chainable error handling.

```go
// Construct from any (T, error) pair
r := result.Of(strconv.Atoi("42"))   // Ok(42)
bad := result.Of(strconv.Atoi("??")) // Err(...)

// Transform — errors propagate automatically
doubled := result.Map(r, func(v int) int { return v * 2 }) // Ok(84)
doubled  = result.Map(bad, func(v int) int { return v * 2 }) // Err(...)

// Unwrap
val, err := r.Unwrap() // 42, nil
val, err  = bad.Unwrap() // 0, <error>

// Fallback
val = bad.GetOrElse(0) // 0
```

---

## pipeline

Function composition, left-to-right. Use `Pipe` for same-type chains; `Pipe2`–`Pipe4` for type-changing chains.

```go
// Same-type composition
clean := pipeline.Pipe(
    strings.TrimSpace,
    strings.ToLower,
)
clean("  Hello World  ") // "hello world"

// Type-changing chain: string → []string → int
countWords := pipeline.Pipe2(
    strings.Fields,
    func(ws []string) int { return len(ws) },
)
countWords("one two three") // 3

// Three-step chain: string → []string → []string → map[string]int
analyse := pipeline.Pipe3(
    strings.Fields,
    func(ws []string) []string {
        return slice.Map(ws, strings.ToLower)
    },
    slice.Frequencies[string],
)
analyse("Go go GO") // map[go:3]
```
