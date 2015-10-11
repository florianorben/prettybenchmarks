# Prettybenchmarks

Prettybenchmarks formats your go benchmarks into nice looking sorted tables

## Installation
    go get github.com/florianorben/prettybenchmarks/cmd/pb

## Usage
Pipe go bench results into "pb"

Works with and without -benchmem flag

If you provide a time interval (either *ns*, *µs* (or *us*), *ms*, *s*), each benchmark's runtime will be converted to that interval. If left blank, a suitable value will automatically be chosen


    go test -bench=YOUR_PKG [-benchmem] | pb [timeinterval]

Example:

    go test -bench=. -benchmem | pb ms

## Features
- Removes clutter in benchmark's names (e.g. Benchmark_, -8 etc.)
- Automatically groups benchmarks if you use Benchmark_FN_XXX notation, where XXX is the number of iterations you run the benchmark (see screenshots)
- Optionally convert *ns* runtime values into a more-readable value (>1000 µs, > 1000000 ms, > 1000000000 s)
- Prints a table ;)

## Screenshots
**Before**

![Before](https://raw.githubusercontent.com/wiki/florianorben/prettybenchmarks/before.png "Before")

**After**

![After](https://raw.githubusercontent.com/wiki/florianorben/prettybenchmarks/after.png "After")

## Misc
Shoutout to [panicparse](https://github.com/maruel/panicparse) for giving me some inspiration

---

[![Build Status](https://travis-ci.org/florianorben/prettybenchmarks.svg?branch=master)](https://travis-ci.org/florianorben/prettybenchmarks)
