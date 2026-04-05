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

// Command wordfreq is an OTP-style word frequency service demonstrating how
// go-functional's concurrency primitives compose:
//
//   - supervisor manages a progress-reporter goroutine with automatic restart
//   - agent accumulates word frequencies from concurrent workers
//   - task processes each document in a separate goroutine
//   - pipeline composes the text-normalisation steps
//   - slice/maps transform and query the final results
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mikehelmick/go-functional/agent"
	fmaps "github.com/mikehelmick/go-functional/maps"
	"github.com/mikehelmick/go-functional/pipeline"
	"github.com/mikehelmick/go-functional/slice"
	"github.com/mikehelmick/go-functional/supervisor"
	"github.com/mikehelmick/go-functional/task"
)

// corpus is the set of "documents" that will be processed concurrently.
var corpus = []string{
	"the quick brown fox jumps over the lazy dog",
	"go is an open source programming language that makes it easy to build software",
	"functional programming makes code easier to reason about and test",
	"the go programming language supports generics and has a strong concurrency model",
	"brown foxes are quick and lazy dogs sleep all day long",
	"the lazy programmer writes functional code to reason about programming problems",
}

// tokenise splits a document into lower-case words using a composed pipeline.
var tokenise = pipeline.Pipe2(
	strings.Fields,
	func(words []string) []string {
		return slice.Map(words, strings.ToLower)
	},
)

func main() {
	// agent holds the running word-frequency totals.
	freqs := agent.New(map[string]int{})

	// supervisor keeps a progress-reporter goroutine alive. If the reporter
	// crashes (returns a non-nil error), supervisor restarts it automatically.
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

	// tasks process each document concurrently using slice.Frequencies.
	tasks := slice.Map(corpus, func(doc string) *task.Task[map[string]int] {
		return task.Run(func() (map[string]int, error) {
			return slice.Frequencies(tokenise(doc)), nil
		})
	})

	results, _ := task.AwaitAll(tasks)

	// Merge all frequency maps into the agent, one at a time.
	for _, r := range results {
		freqs.Update(func(cur map[string]int) map[string]int {
			for word, count := range r {
				cur[word] += count
			}
			return cur
		})
	}

	sup.Stop()
	final := freqs.Get()
	freqs.Stop()

	// Build a sorted list of (word, count) pairs, descending by frequency.
	type wordCount struct {
		word  string
		count int
	}
	entries := slice.Map(fmaps.Keys(final), func(w string) wordCount {
		return wordCount{w, final[w]}
	})
	sorted := slice.SortBy(entries, func(e wordCount) int { return -e.count })

	fmt.Println("\nTop 5 words:")
	for _, e := range slice.Take(sorted, 5) {
		fmt.Printf("  %-14s %d\n", e.word, e.count)
	}

	repeated := slice.Filter(sorted, func(e wordCount) bool { return e.count > 1 })
	fmt.Printf("\n%d of %d unique words appear more than once\n", len(repeated), len(final))
}
