// Copyright 2015 Florian Orben. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//Package prettybenchmarks formats your go benchmarks into nice looking sorted tables
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
//
//Turns
//    Benchmark_NewSmallReq-8      	  100000	     21618 ns/op	    2739 B/op	      45 allocs/op
//    BenchmarkNewLargeReq-8      	   10000	    122245 ns/op	   29823 B/op	      54 allocs/op
//    Benchmark_NewSmallReqProto-8 	  100000	     15594 ns/op	    2691 B/op	      44 allocs/op
//
//into
//    +-----------------------+---------+---------+--------+----------------+
//    | Name                  |    Runs |   µs/op |   B/op | allocations/op |
//    +-----------------------+---------+---------+--------+----------------+
//    | NewLargeReq           |  10,000 | 109.805 | 29,823 |             54 |
//    +-----------------------+---------+---------+--------+----------------+
//    | NewSmallReq           | 100,000 |  14.122 |  2,739 |             45 |
//    +-----------------------+---------+---------+--------+----------------+
//    | NewSmallReqProto      | 100,000 |  13.959 |  2,691 |             44 |
//    +-----------------------+---------+---------+--------+----------------+
package main

import "github.com/florianorben/prettybenchmarks/prettybenchmarks"

func main() {
	prettybenchmarks.Main()
}
