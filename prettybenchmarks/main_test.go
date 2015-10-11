package prettybenchmarks

import (
	"reflect"
	"testing"
)

var tests = []struct {
	input        [][]byte
	expected     *results
	expectedInfo *benchmarkInfo
	extraLines   []string
}{
	{
		[][]byte{
			[]byte("PASS\n"),
			[]byte("Benchmark_NewSmallReq-8      	  100000	     21618 ns/op	    2739 B/op	      45 allocs/op\n"),
			[]byte("BenchmarkNewLargeReq-8      	   10000	    122245 ns/op	   29823 B/op	      54 allocs/op\n"),
			[]byte("Benchmark_NewSmallReqProto-8 	  100000	     15594 ns/op	    2691 B/op	      44 allocs/op\n"),
			[]byte("BenchmarkNewLargeReqProto-8 	   10000	    170835 ns/op	   26706 B/op	      53 allocs/op\n"),
			[]byte("Benchmark_UnmarshalSmallReq-8	  100000	     22231 ns/op	    3570 B/op	     100 allocs/op\n"),
			[]byte("Benchmark_UnmarshalLargeReq_10-8    5000	    342400 ns/op	   60385 B/op	    1680 allocs/op\n"),
			[]byte("Benchmark_UnmarshalLargeReq_100-8   5000	    342400 ns/op	   60385 B/op	    1680 allocs/op\n"),
			[]byte("Benchmark_UnmarshalLargeReq_1000-8  5000	    342400 ns/op	   60385 B/op	    1680 allocs/op\n"),
			[]byte("ok  	github.com/foobar/baz	11.164s\n"),
		},
		&results{
			"NewSmallReq": []*result{
				{"NewSmallReq", -1, 100000, float64(21618), 2739, 45},
			},
			"UnmarshalSmallReq": []*result{
				{"UnmarshalSmallReq", -1, 100000, float64(22231), 3570, 100},
			},
			"NewLargeReq": []*result{
				{"NewLargeReq", -1, 10000, float64(122245), 29823, 54},
			},
			"NewSmallReqProto": []*result{
				{"NewSmallReqProto", -1, 100000, float64(15594), 2691, 44},
			},
			"NewLargeReqProto": []*result{
				{"NewLargeReqProto", -1, 10000, float64(170835), 26706, 53},
			},
			"UnmarshalLargeReq": []*result{
				{"UnmarshalLargeReq", 10, 5000, float64(342400), 60385, 1680},
				{"UnmarshalLargeReq", 100, 5000, float64(342400), 60385, 1680},
				{"UnmarshalLargeReq", 1000, 5000, float64(342400), 60385, 1680},
			},
		},
		&benchmarkInfo{true, true, "µs"},
		[]string{
			"PASS\n",
			"ok  	github.com/foobar/baz	11.164s\n",
		},
	},
	{
		[][]byte{
			[]byte("FOO\n"),
			[]byte("\n"),
			[]byte("BAR\n\n"),
			[]byte("Benchmark_NewSmallReq-8      	  100000	     21 ns/op	    2739 B/op	      45 allocs/op\n"),
			[]byte("BenchmarkNewLargeReq-8      	   10000	    122 ns/op	   29823 B/op	      54 allocs/op\n"),
			[]byte("NO_BenchmarkNewLargeReq-8      	   10000	    210 ns/op	   29823 B/op	      54 allocs/op\n"),
			[]byte("ok  	github.com/foobar/baz	11.164s\n"),
			[]byte("fail  	github.com/foobar/baz	11.164s\n"),
			[]byte("??  	github.com/foobar/baz	11.164s\n"),
		},
		&results{
			"NewSmallReq": []*result{
				{"NewSmallReq", -1, 100000, float64(21), 2739, 45},
			},
			"NewLargeReq": []*result{
				{"NewLargeReq", -1, 10000, float64(122), 29823, 54},
			},
		},
		&benchmarkInfo{false, true, "µs"},
		[]string{
			"FOO\n",
			"\n",
			"BAR\n\n",
			"NO_BenchmarkNewLargeReq-8      	   10000	    210 ns/op	   29823 B/op	      54 allocs/op\n",
			"ok  	github.com/foobar/baz	11.164s\n",
			"fail  	github.com/foobar/baz	11.164s\n",
			"??  	github.com/foobar/baz	11.164s\n",
		},
	},

	{
		[][]byte{
			[]byte("PASS\n"),
			[]byte("Benchmark_NewSmallReq-8      	  100000	     21618 ns/op\n"),
			[]byte("BenchmarkNewLargeReq-8      	   10000	    122245 ns/op\n"),
			[]byte("Benchmark_NewSmallReqProto-8 	  100000	     15594 ns/op\n"),
			[]byte("BenchmarkNewLargeReqProto-8 	   10000	    170835 ns/op\n"),
			[]byte("Benchmark_UnmarshalSmallReq-8	  100000	     22231 ns/op\n"),
			[]byte("Benchmark_UnmarshalLargeReq_10-8    5000	    342400 ns/op\n"),
			[]byte("Benchmark_UnmarshalLargeReq_100-8   5000	    342400 ns/op\n"),
			[]byte("Benchmark_UnmarshalLargeReq_1000-8  5000	    342400 ns/op\n"),
			[]byte("ok  	github.com/foobar/baz	22222.164s\n"),
		},
		&results{
			"NewSmallReq": []*result{
				{"NewSmallReq", -1, 100000, float64(21618), -1, -1},
			},
			"UnmarshalSmallReq": []*result{
				{"UnmarshalSmallReq", -1, 100000, float64(22231), -1, -1},
			},
			"NewLargeReq": []*result{
				{"NewLargeReq", -1, 10000, float64(122245), -1, -1},
			},
			"NewSmallReqProto": []*result{
				{"NewSmallReqProto", -1, 100000, float64(15594), -1, -1},
			},
			"NewLargeReqProto": []*result{
				{"NewLargeReqProto", -1, 10000, float64(170835), -1, -1},
			},
			"UnmarshalLargeReq": []*result{
				{"UnmarshalLargeReq", 10, 5000, float64(342400), -1, -1},
				{"UnmarshalLargeReq", 100, 5000, float64(342400), -1, -1},
				{"UnmarshalLargeReq", 1000, 5000, float64(342400), -1, -1},
			},
		},
		&benchmarkInfo{true, false, "µs"},
		[]string{
			"PASS\n",
			"ok  	github.com/foobar/baz	22222.164s\n",
		},
	},
	{
		[][]byte{
			[]byte("FOO\n"),
			[]byte("\n"),
			[]byte("BAR\n\n"),
			[]byte("Benchmark_NewSmallReq-8      	  100000	     216180000 ns/op\n"),
			[]byte("NO_BenchmarkNewLargeReq-8      	   10000	    1200021245 ns/op\n"),
			[]byte("BenchmarkNewLargeReq-8      	   10000	    1222450320 ns/op\n"),
			[]byte("ok  	github.com/foobar/baz	11.164s\n"),
			[]byte("fail  	github.com/foobar/baz	11.164s\n"),
			[]byte("?foo?  	github.com/foobar/baz	11.164s\n"),
		},
		&results{
			"NewSmallReq": []*result{
				{"NewSmallReq", -1, 100000, float64(216180000), -1, -1},
			},
			"NewLargeReq": []*result{
				{"NewLargeReq", -1, 10000, float64(1222450320), -1, -1},
			},
		},
		&benchmarkInfo{false, false, "µs"},
		[]string{
			"FOO\n",
			"\n",
			"BAR\n\n",
			"NO_BenchmarkNewLargeReq-8      	   10000	    1200021245 ns/op\n",
			"ok  	github.com/foobar/baz	11.164s\n",
			"fail  	github.com/foobar/baz	11.164s\n",
			"?foo?  	github.com/foobar/baz	11.164s\n",
		},
	},
}

func Test_newResults(t *testing.T) {
	for _, tt := range tests {
		actual := newResults(tt.input)

		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("Constructing new results %s: expected %#v, actual %#v\n", tt.input, tt.expected, actual)
		}
	}
}

func Test_benchInfo(t *testing.T) {
	for _, tt := range tests {
		tmp := newBenchmark(tt.input)
		actual := tmp.info

		if !reflect.DeepEqual(actual, tt.expectedInfo) {
			t.Errorf("Constructing new info %s: expected %#v, actual %#v\n", tt.input, tt.expectedInfo, actual)
		}
	}
}

func Test_unparsableLines(t *testing.T) {
	for _, tt := range tests {
		actual := make([]string, 0)
		for _, line := range tt.input {
			_, err := newResult(line)
			if err != nil {
				actual = append(actual, err.Error())
				continue
			}
		}

		if !reflect.DeepEqual(actual, tt.extraLines) {
			t.Errorf("Awaiting following unparsable lines for input %s: expected %#v, actual %#v\n", tt.input, tt.extraLines, actual)
		}
	}
}

func Test_bold(t *testing.T) {
	s := "foo"
	expected := "\033[1mfoo\033[0m"
	actual := bold(s)
	if expected != actual {
		t.Errorf("Formatting bold string %: expected %v, actual %v", s, expected, actual)
	}
}
func Test_green(t *testing.T) {
	s := "foo"
	expected := "\033[32mfoo\033[0m"
	actual := green(s)
	if expected != actual {
		t.Errorf("Formatting green string %: expected %v, actual %v", s, expected, actual)
	}
}
func Test_red(t *testing.T) {
	s := "foo"
	expected := "\033[31mfoo\033[0m"
	actual := red(s)
	if expected != actual {
		t.Errorf("Formatting red string %: expected %v, actual %v", s, expected, actual)
	}
}
func Test_gray(t *testing.T) {
	s := "foo"
	expected := "\033[90mfoo\033[0m"
	actual := gray(s)
	if expected != actual {
		t.Errorf("Formatting gray string %: expected %v, actual %v", s, expected, actual)
	}
}
