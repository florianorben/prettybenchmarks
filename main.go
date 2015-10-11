// prettybenchmark: format go benchmarks into tables
//
// Usage: Pipe your benchmark results into "pb"
//    go test -bench=. [-benchmem] | pb
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apcera/termtables"
)

var (
	byWhitespace = regexp.MustCompile(`\s+`)
	byRuns       = regexp.MustCompile(`-\d+$`)
	byIterations = regexp.MustCompile(`(?i:)(Benchmark_?)`)
)

type (
	Foo       struct{}
	benchmark struct {
		info    *benchmarkInfo
		results *results
	}
	benchmarkInfo struct {
		hasFnIterations bool
		benchmemUsed    bool
		suggestedTiming string
	}
	results map[string][]*result
	result  struct {
		Name         string
		FnIterations int
		Runs         int
		Speed        float64
		Bps          int
		Aps          int
	}
)

type sortByFnIterations []*result

func (b sortByFnIterations) Len() int           { return len(b) }
func (b sortByFnIterations) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b sortByFnIterations) Less(i, j int) bool { return b[i].FnIterations < b[j].FnIterations }

var (
	lines [][]byte
	table *termtables.Table
	bench *benchmark
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	quit := make(chan bool)
	go loading(quit)

	for {
		text, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		lines = append(lines, text)
	}
	close(quit)

	bench = newBenchmark(lines)

	table = termtables.CreateTable()
	addTableHeader(table)
	addTableBody(table)

	fmt.Print("\r \n")
	fmt.Println(table.Render())
	fmt.Println(footer())
}

func newBenchmark(l [][]byte) *benchmark {
	results := newResults(l)
	return &benchmark{
		info:    newBenchmarkInfo(results),
		results: results,
	}
}

func newResults(l [][]byte) *results {
	benchMap := make(results)

	for _, l := range l[1 : len(l)-1] {
		bl := newResult(l)
		if bl != nil {
			if _, ok := benchMap[bl.Name]; !ok {
				benchMap[bl.Name] = make([]*result, 0)
			}
			benchMap[bl.Name] = append(benchMap[bl.Name], bl)
		}
	}

	for _, r := range benchMap {
		sort.Sort(sortByFnIterations(r))
	}

	return &benchMap
}

func newResult(b []byte) *result {
	var (
		name   string
		fnIter int
		bps    int
		aps    int
		err    error
		iter   int
		speed  float64
	)

	s := string(b)
	parts := byWhitespace.Split(s, -1)
	nameRuns := byRuns.ReplaceAllString(parts[0], "")
	nameIterations := byIterations.ReplaceAllString(nameRuns, "")
	lastIndex := strings.LastIndex(nameIterations, "_")

	if lastIndex > -1 {
		name = nameIterations[:lastIndex]
		fnIter, _ = strconv.Atoi(nameIterations[lastIndex+1:])
	} else {
		name = nameIterations
		fnIter = -1
	}

	//just print the line if it doesn't have the correct format
	if len(parts) < 4 {
		fmt.Println(s)
		return nil
	}

	iter, err = strconv.Atoi(parts[1])
	if err != nil {
		iter = -1
	}
	speed, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		speed = -1
	}

	if len(parts) > 5 {
		bps, err = strconv.Atoi(parts[4])
		if err != nil {
			bps = -1
		}
		aps, err = strconv.Atoi(parts[6])
		if err != nil {
			aps = -1
		}
	} else {

		//without benchmem
		bps = -1
		aps = -1
	}

	return &result{
		Name:         name,
		FnIterations: fnIter,
		Runs:         iter,
		Speed:        speed,
		Bps:          bps,
		Aps:          aps,
	}
}

func newBenchmarkInfo(r *results) *benchmarkInfo {
	var (
		slowest         float64
		hasFnIterations bool
		benchmemUsed    bool
		suggestedTiming string
	)
	suggestedTiming = "ns"

	for _, bl := range *r {
		for _, l := range bl {

			if l.FnIterations > -1 {
				hasFnIterations = true
			}

			if l.Aps > -1 && l.Bps > -1 {
				benchmemUsed = true
			}

			if slowest < l.Speed {
				slowest = l.Speed
			}
		}
	}

	switch {
	case slowest <= 1e3:
		suggestedTiming = "ns"
	case slowest > 1e3 && slowest <= 1e6:
		suggestedTiming = "µs"
		updateSpeedVals(r, float64(1e3))
	case slowest > 1e6 && slowest <= 1e9:
		suggestedTiming = "ms"
		updateSpeedVals(r, float64(1e6))
	case slowest > 1e9:
		suggestedTiming = "s"
		updateSpeedVals(r, float64(1e9))
	}

	return &benchmarkInfo{hasFnIterations, benchmemUsed, suggestedTiming}
}

func updateSpeedVals(r *results, f float64) {
	for _, bl := range *r {
		for _, l := range bl {
			l.Speed = l.Speed / f
		}
	}
}

func footer() string {

	lastLine := bytes.Replace(
		bytes.TrimSpace(lines[len(lines)-1]),
		[]byte{9},
		[]byte{32, 32, 32, 32, 32},
		-1,
	)
	footer := make([]byte, 0, len(lastLine)*2+1)

	footer = append(
		footer,
		lastLine...,
	)

	footer = append(footer, byte(10))
	footer = append(footer, byte(43))
	for i := 0; i < len(lastLine)-2; i++ {
		footer = append(footer, byte(45))
	}
	footer = append(footer, byte(43))

	return bold(string(footer))
}

func addTableHeader(t *termtables.Table) {
	if bench.info.benchmemUsed {
		if bench.info.hasFnIterations {
			t.AddHeaders(bold("Name"), bold("Iterations"), bold("Runs"), bold(bench.info.suggestedTiming+"/op"), bold("B/op"), bold("allocations/op"), "")
		} else {
			t.AddHeaders(bold("Name"), bold("Runs"), bold(bench.info.suggestedTiming+"/op"), bold("B/op"), bold("allocations/op"), "")
		}
	} else {
		if bench.info.hasFnIterations {
			t.AddHeaders(bold("Name"), bold("Iterations"), bold("Runs"), bold(bench.info.suggestedTiming+"/op"), "")
		} else {
			t.AddHeaders(bold("Name"), bold("Runs"), bold(bench.info.suggestedTiming+"/op"), "")
		}
	}
}

func addTableBody(t *termtables.Table) {
	i := len(*bench.results)
	sorted := make([]string, 0, i)
	for name := range *bench.results {
		sorted = append(sorted, name)
	}
	sort.Sort(sort.StringSlice(sorted))

	for _, benchName := range sorted {
		results := (*bench.results)[benchName]

		for j, b := range results {
			var name string
			if j == 0 {
				name = bold(b.Name)
			}

			if bench.info.benchmemUsed {
				if bench.info.hasFnIterations {
					fnIterations := strconv.Itoa(b.FnIterations)

					if fnIterations == "-1" {
						fnIterations = ""
					}

					t.AddRow(name, fnIterations, b.Runs, b.Speed, b.Bps, b.Aps, "⬅")
				} else {
					t.AddRow(name, b.Runs, b.Speed, b.Bps, b.Aps, "⬅")
				}
			} else {
				if bench.info.hasFnIterations {
					fnIterations := strconv.Itoa(b.FnIterations)

					if fnIterations == "-1" {
						fnIterations = ""
					}

					t.AddRow(name, fnIterations, b.Runs, b.Speed, "⬅")
				} else {
					t.AddRow(name, b.Runs, b.Speed, "⬅")
				}
			}
		}

		i--
		if i > 0 {
			t.AddSeparator()
		}
	}

	t.SetAlign(termtables.AlignLeft, 1)
	t.SetAlign(termtables.AlignRight, 2)
	t.SetAlign(termtables.AlignRight, 3)
	t.SetAlign(termtables.AlignRight, 4)
	t.SetAlign(termtables.AlignRight, 5)
	t.SetAlign(termtables.AlignRight, 6)
}

func loading(q chan bool) {
	states := []string{"|", "/", "-", "\\", "|", "/", "–", "\\"}
	current := 0

	for {
		select {
		case <-time.Tick(150 * time.Millisecond):
			fmt.Printf("\r%s", states[current])
			if current == len(states)-1 {
				current = 0
			} else {
				current++
			}
		case <-q:
			break
		}
	}
}

func bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}
