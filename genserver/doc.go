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

// Package genserver provides a generic server process inspired by Elixir's
// OTP GenServer. State is owned by a single goroutine; Call and Cast
// messages are serialized through a single channel to preserve ordering.
//
// Call is synchronous: the caller blocks until HandleCall returns a response.
// Cast is asynchronous: the caller returns immediately after sending.
//
// Without tail-call optimization in Go, the server loop is an explicit
// channel select rather than a recursive function, but the semantics are
// equivalent.
package genserver
