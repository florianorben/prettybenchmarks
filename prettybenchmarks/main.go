package prettybenchmarks

import (
	"bufio"
	"flag"
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

const (
	fmtInt     = "#,###."
	fmtFloat   = "#,###.###"
	fmtFloatNS = "#,###."
)

type (
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
	regExByWhitespace = regexp.MustCompile(`\s+`)
	regExByRuns       = regexp.MustCompile(`-\d+$`)
	regExByIterations = regexp.MustCompile(`(?i:)(^Benchmark_?)`)
	regExIsBenchmark  = regExByIterations
	linePassed        = "PASS"
	lineSkipped       = "SKIP"
	lineFail          = "FAIL"
)

var (
	lines           [][]byte
	unparsableLines []string
	table           *termtables.Table
	bench           *benchmark
	timing          string
)

func init() {
	setTiming()
}

func Main() {

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

	if len(lines) == 0 {
		os.Exit(0)
	}

	bench = newBenchmark(lines)

	table = termtables.CreateTable()
	table.Style.Alignment = termtables.AlignRight
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

	for _, l := range l {
		bl, err := newResult(l)
		if err != nil {
			unparsableLines = append(unparsableLines, err.Error())
			continue
		}
		if _, ok := benchMap[bl.Name]; !ok {
			benchMap[bl.Name] = make([]*result, 0)
		}
		benchMap[bl.Name] = append(benchMap[bl.Name], bl)

	}

	for _, r := range benchMap {
		sort.Sort(sortByFnIterations(r))
	}

	return &benchMap
}

func newResult(b []byte) (*result, error) {
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
	parts := regExByWhitespace.Split(s, -1)
	if len(parts) < 4 || !regExIsBenchmark.MatchString(parts[0]) {
		return nil, fmt.Errorf("%s", s)
	}

	nameRuns := regExByRuns.ReplaceAllString(parts[0], "")
	nameIterations := regExByIterations.ReplaceAllString(nameRuns, "")
	lastIndex := strings.LastIndex(nameIterations, "_")

	if lastIndex > -1 {
		name = nameIterations[:lastIndex]
		fnIter, _ = strconv.Atoi(nameIterations[lastIndex+1:])
	} else {
		name = nameIterations
		fnIter = -1
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
	}, nil
}

func newBenchmarkInfo(r *results) *benchmarkInfo {
	var (
		slowest         float64
		hasFnIterations bool
		benchmemUsed    bool
	)

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

	if timing == "" {
		switch {
		case slowest <= 1e3:
			timing = "ns"
		case slowest > 1e3 && slowest <= 1e6:
			timing = "µs"
		case slowest > 1e6 && slowest <= 1e9:
			timing = "ms"
		case slowest > 1e9:
			timing = "s"
		}
	}

	switch timing {
	case "ns":
		//ns is default, dont't do anything
	case "µs":
		updateSpeedVals(r, float64(1e3))
	case "ms":
		updateSpeedVals(r, float64(1e6))
	case "s":
		updateSpeedVals(r, float64(1e9))
	}

	return &benchmarkInfo{hasFnIterations, benchmemUsed, timing}
}

func updateSpeedVals(r *results, f float64) {
	for _, bl := range *r {
		for _, l := range bl {
			l.Speed = l.Speed / f
		}
	}
}

func footer() string {
	var footer []byte

	footer = append(footer, []byte{10}...)
	footer = append(footer, []byte((bold("Summary:"))+"\n")...)
	footer = append(footer, []byte((bold("+------+"))+"\n")...)

	for _, line := range unparsableLines {
		tmp := strings.TrimSpace(line)
		switch {
		case tmp == linePassed:
			footer = append(footer, []byte(green(bold(tmp))+"\n")...)
		case tmp == lineSkipped:
			footer = append(footer, []byte(gray(bold(tmp))+"\n")...)
		case tmp == lineFail:
			footer = append(footer, []byte(red(bold(tmp))+"\n")...)
		default:
			footer = append(footer, []byte(tmp+"\n")...)
		}
	}

	return string(footer)
}

func addTableHeader(t *termtables.Table) {
	var lenLongestName int
	for name := range *bench.results {
		if tmpLen := len(name); tmpLen > lenLongestName {
			lenLongestName = tmpLen
		}
	}

	// add padding to first col since alignment in header columns does not work
	// padding of longest name + len("name") + 1 padding right
	nameCol := make([]byte, 0, lenLongestName+4+1)
	nameCol = append(nameCol, []byte("Name")...)
	for i := 0; i < lenLongestName; i++ {
		nameCol = append(nameCol, byte(32))
	}

	if bench.info.benchmemUsed {
		if bench.info.hasFnIterations {
			t.AddHeaders(bold(string(nameCol)), bold("Iterations"), bold("Runs"), bold(bench.info.suggestedTiming+"/op"), bold("B/op"), bold("allocations/op"))
		} else {
			t.AddHeaders(bold(string(nameCol)), bold("Runs"), bold(bench.info.suggestedTiming+"/op"), bold("B/op"), bold("allocations/op"))
		}
	} else {
		if bench.info.hasFnIterations {
			t.AddHeaders(bold(string(nameCol)), bold("Iterations"), bold("Runs"), bold(bench.info.suggestedTiming+"/op"))
		} else {
			t.AddHeaders(bold(string(nameCol)), bold("Runs"), bold(bench.info.suggestedTiming+"/op"))
		}
	}
}

func addTableBody(t *termtables.Table) {
	floatFmt := fmtFloat
	if bench.info.suggestedTiming == "ns" {
		floatFmt = fmtFloatNS
	}

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
					} else {
						i, err := strconv.Atoi(fnIterations)
						if err != nil {
							fnIterations = ""
						} else {
							fnIterations = RenderInteger(fmtInt, i)
						}
					}

					t.AddRow(name, fnIterations, RenderInteger(fmtInt, b.Runs), RenderFloat(floatFmt, b.Speed), RenderInteger(fmtInt, b.Bps), RenderInteger(fmtInt, b.Aps))
				} else {
					t.AddRow(name, RenderInteger(fmtInt, b.Runs), RenderFloat(floatFmt, b.Speed), RenderInteger(fmtInt, b.Bps), RenderInteger(fmtInt, b.Aps))
				}
			} else {
				if bench.info.hasFnIterations {
					fnIterations := strconv.Itoa(b.FnIterations)

					if fnIterations == "-1" {
						fnIterations = ""
					} else {
						i, err := strconv.Atoi(fnIterations)
						if err != nil {
							fnIterations = ""
						} else {
							fnIterations = RenderInteger(fmtInt, i)
						}
					}

					t.AddRow(name, fnIterations, RenderInteger(fmtInt, b.Runs), RenderFloat(floatFmt, b.Speed))
				} else {
					t.AddRow(name, RenderInteger(fmtInt, b.Runs), RenderFloat(floatFmt, b.Speed))
				}
			}
		}

		i--
		if i > 0 {
			t.AddSeparator()
		}
	}

	t.SetAlign(termtables.AlignLeft, 1)
}

func setTiming() {
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		if lowerArg := strings.ToLower(args[0]); lowerArg == "ns" || lowerArg == "us" || lowerArg == "µs" || lowerArg == "ms" || lowerArg == "s" {
			if lowerArg == "us" {
				lowerArg = "µs"
			}
			timing = lowerArg
		}
	}
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

func green(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func gray(s string) string {
	return fmt.Sprintf("\033[90m%s\033[0m", s)
}
