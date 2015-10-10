package main

import (
	"bufio"
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
	r  = regexp.MustCompile(`\s+`)
	r2 = regexp.MustCompile(`-\d+$`)
	r3 = regexp.MustCompile(`(?i:)(Benchmark_?)`)
)

type benchLine struct {
	Name         string
	FnIterations int
	Iterations   int
	Speed        float64
	Bps          int
	Aps          int
	Status       string
}

var lines [][]byte

func main() {
	// old := os.Stdout // keep backup of the real stdout
	// r, w, err := os.Pipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// os.Stdout = w
	//
	// // print() // fails if called here "fatal error: all goroutines are asleep - deadlock"
	//
	// outC := make(chan bytes.Buffer)
	// // copy the output in a separate goroutine so printing can't block indefinitely
	// go func() {
	// 	var buf bytes.Buffer
	// 	io.Copy(&buf, r)
	// 	outC <- buf
	// }()
	//
	// // back to normal state
	// w.Close()
	// os.Stdout = old // restoring the real stdout
	// out := <-outC
	//
	// // reading our temp stdout
	//
	// //fmt.Print(out)
	//
	// var reader *bufio.Reader
	// if len(out.Bytes()) > 0 {
	// 	reader = bufio.NewReader(&out)
	// } else {
	// 	reader = bufio.NewReader(os.Stdin)
	// }

	fmt.Print("\n")
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

	benchMap := getNameToBenchmarkLinesMap(lines)
	hasFnIterations, speedTblHeader := getInfoForLines(benchMap)

	table := termtables.CreateTable()
	if hasFnIterations {
		table.AddHeaders(bold("Name"), bold("Iterations"), bold("No."), bold(speedTblHeader+"/op"), bold("B/op"), bold("allocations/op"), "")
	} else {
		table.AddHeaders(bold("Name"), bold("No."), bold(speedTblHeader+"/op"), bold("B/op"), bold("allocations/op"), "")
	}

	i := len(benchMap)
	sorted := make([]string, 0, len(benchMap))
	for name := range benchMap {
		sorted = append(sorted, name)
	}
	sort.Sort(sort.StringSlice(sorted))

	for _, benchName := range sorted {
		benchLines := benchMap[benchName]

		for j, b := range benchLines {
			var name string
			if j == 0 {
				name = bold(b.Name)
			}

			if hasFnIterations {
				fnIterations := strconv.Itoa(b.FnIterations)

				if fnIterations == "-1" {
					fnIterations = ""
				}

				table.AddRow(name, fnIterations, b.Iterations, b.Speed, b.Bps, b.Aps, b.Status)
			} else {
				table.AddRow(name, b.Iterations, b.Speed, b.Bps, b.Aps, b.Status)
			}
		}

		i--
		if i > 0 {
			table.AddSeparator()
		}
	}

	table.SetAlign(termtables.AlignLeft, 1)
	table.SetAlign(termtables.AlignRight, 2)
	table.SetAlign(termtables.AlignRight, 3)
	table.SetAlign(termtables.AlignRight, 4)
	table.SetAlign(termtables.AlignRight, 5)
	table.SetAlign(termtables.AlignRight, 6)

	fmt.Print("\r")
	footer := strings.TrimSpace(string(lines[len(lines)-1]))
	footer = strings.Replace(footer, "\t", "     ", -1)
	padLeft := (table.Style.Width - len(footer)) / 4
	for i := 0; i < padLeft; i++ {
		fmt.Print(bold(" "))
	}
	fmt.Print(bold(footer))
	fmt.Print("\n")
	for i := 0; i < padLeft-1; i++ {
		fmt.Print(bold(" "))
	}
	fmt.Print(bold("+"))
	for range footer {
		fmt.Print(bold("-"))
	}
	fmt.Print(bold("+"))
	fmt.Print("\n")

	fmt.Println(table.Render())
}

func getNameToBenchmarkLinesMap(lines [][]byte) map[string][]*benchLine {
	benchMap := make(map[string][]*benchLine)

	for _, l := range lines[1 : len(lines)-1] {
		bl := newbenchLineFromString(string(l))
		if _, ok := benchMap[bl.Name]; !ok {
			benchMap[bl.Name] = make([]*benchLine, 0)
		}
		benchMap[bl.Name] = append(benchMap[bl.Name], bl)
	}

	return benchMap
}

func newbenchLineFromString(s string) *benchLine {
	//qwe := "Benchmark_Sort_1000-8                	   20000	     84093 ns/op	      32 B/op	       1 allocs/op"

	y := r.Split(s, -1)
	x := r2.ReplaceAllString(y[0], "")
	c := r3.ReplaceAllString(x, "")
	lastIndex := strings.LastIndex(c, "_")

	var name string
	var fnIter int

	if lastIndex > -1 {
		name = c[:lastIndex]
		fnIter, _ = strconv.Atoi(c[lastIndex+1:])
	} else {
		name = c
		fnIter = -1
	}

	iter, err := strconv.Atoi(y[1])
	if err != nil {
		iter = -1
	}
	speed, err := strconv.ParseFloat(y[2], 64)
	if err != nil {
		speed = -1
	}
	bps, err := strconv.Atoi(y[4])
	if err != nil {
		bps = -1
	}
	aps, err := strconv.Atoi(y[6])
	if err != nil {
		aps = -1
	}

	return &benchLine{
		Name:         name,
		FnIterations: fnIter,
		Iterations:   iter,
		Speed:        speed,
		Bps:          bps,
		Aps:          aps,
	}
}

func getInfoForLines(bm map[string][]*benchLine) (bool, string) {
	var slowest float64
	var hasIterations = false
	var speedTblHeader = "ns"

	for _, bl := range bm {

		bestSpeedPerLine := float64(-1)
		bestLine := 0
		worstLine := 0

		for i, l := range bl {

			if l.FnIterations > -1 {
				hasIterations = true
			}

			if slowest < l.Speed {
				slowest = l.Speed
				worstLine = i
			}

			if bestSpeedPerLine == -1 || bestSpeedPerLine > l.Speed {
				bestSpeedPerLine = l.Speed
				bestLine = i
			}
		}

		bl[worstLine].Status = red("✘")
		bl[bestLine].Status = green("✔")
	}

	switch {
	case slowest <= 1e3:
		speedTblHeader = "ns"
	case slowest > 1e3 && slowest <= 1e6:
		speedTblHeader = "µs"
		updateSpeedVals(bm, float64(1e3))
	case slowest > 1e6 && slowest <= 1e9:
		speedTblHeader = "ms"
		updateSpeedVals(bm, float64(1e6))
	case slowest > 1e9:
		speedTblHeader = "s"
		updateSpeedVals(bm, float64(1e9))
	}

	return hasIterations, speedTblHeader
}

func updateSpeedVals(bm map[string][]*benchLine, f float64) {
	for _, bl := range bm {
		for _, l := range bl {
			l.Speed = l.Speed / f
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

func blue(s string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", s)
}

func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}
