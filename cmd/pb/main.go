// Copyright 2015 Florian Orben. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// prettybenchmarks: format go benchmarks into tables
//
// Usage: Pipe your benchmark results into "pb"
//    go test -bench=. [-benchmem] | pb [ns/µs/ms/s]

package main

import "github.com/florianorben/prettybenchmarks/prettybenchmarks"

func main() {
	prettybenchmarks.Main()
}
