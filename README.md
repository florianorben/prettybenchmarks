# Prettybenchmarks

Prettybenchmarks formats your go benchmarks into nice looking sorted tables

## Usage
Pipe go bench results into "pb"
Works with and without -benchmem flag

    go test -bench [-benchmem] | pb

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
