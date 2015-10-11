// Copyright 2015 Florian Orben. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//Prettybenchmarks formats your go benchmarks into nice looking sorted tables
//
//Prettybenchmarks
//
//Works with and without -benchmem flag
//
//If you provide a time interval (either ns, µs (or us), ms, s), each benchmark's runtime will be
//converted to that interval. If left blank, a suitable value will automatically be chosen
//
//    go test -bench=YOUR_PKG [-benchmem] | pb [timeinterval]
//Example
//    go test -bench=. -benchmem | pb ms
//
package main

import "github.com/florianorben/prettybenchmarks/prettybenchmarks"

func main() {
	prettybenchmarks.Main()
}
