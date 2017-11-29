package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	flag "github.com/spf13/pflag"
)

// https://godoc.org/github.com/pborman/getopt

type timeSeries struct {
	Series []float64
}

func shellOutNum(cmd string) float64 {
	out, _ := exec.Command("sh", "-c", cmd).Output()
	r := regexp.MustCompile("[\\d,\\.]+")
	cleaned := r.FindAllString(string(out), 1)
	parsed, _ := strconv.ParseFloat(cleaned[0], 64)
	return parsed
}

func (ts *timeSeries) getMax() (max float64) {
	for _, i := range ts.Series {
		if i > max {
			max = i
		}
	}
	return
}

func (ts *timeSeries) getMin() (min float64) {
	min = ts.Series[0]
	for _, i := range ts.Series {
		if i < min {
			min = i
		}
	}
	return
}

func main() {
	// parse arguments
	ConfInterval := flag.Float64P("interval", "n", 1, "seconds to wait between updates")
	flag.Parse()
	ConfArguments := flag.Args()
	fmt.Println(ConfArguments)

	if len(ConfArguments) == 0 {
		flag.Usage()
	}

	// Build UI
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer func() {
		ui.Close()
		fmt.Print("\033[2J") // Clear
	}()

	// build graph
	x := timeSeries{Series: []float64{}}
	g := ui.NewLineChart()
	g.Data = x.Series
	g.Height = ui.TermHeight()
	g.Width = ui.TermWidth()
	g.BorderLabel = "Every " + strconv.FormatFloat(*ConfInterval, 'f', -1, 64) + "s: " + strings.Join(ConfArguments, " ")

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/C-x", func(ui.Event) {
		ui.StopLoop()
	})

	go func() {
		// none of this math is right. How do you get the capacity of g.Data ???
		for {
			nextval := shellOutNum(strings.Join(ConfArguments, " "))

			x.Series = append(x.Series, nextval)

			if len(x.Series) > ui.TermWidth()*2 { // Brail is 2 wide
				g.Data = x.Series[len(x.Series)-ui.TermWidth()*2:]
			} else {
				g.Data = x.Series
			}
			g.Width = ui.TermWidth()
			g.Height = ui.TermHeight()
			ui.Render(g)
			time.Sleep(time.Millisecond * time.Duration(*ConfInterval*1000))
		}
	}()

	// ui.Handle("/sys/wnd/resize", func(e ui.Event) {
	// 	ui.Clear()
	// 	g.Width = ui.TermWidth()
	// 	g.Height = ui.TermHeight()
	// 	ui.Render(g)
	// })

	ui.Loop()

}
