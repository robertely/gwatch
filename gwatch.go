package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	getopt "github.com/pborman/getopt/v2"
	ui "github.com/robertely/termui"
)

type config struct {
	Arguments []string
	// Watch flags were going to mock
	// differences bool // not sure what I could do with this.
	// precise bool // Naw not going to do this.
	// chgexit bool // N/A here
	// color bool // N/A here
	// precise bool // Naw not going to do this.
	Interval float64
	NoTitle  bool
	Beep     bool
	ErrExit  bool
	Exec     bool
	Help     bool
	Version  bool
	// gwatch specific flags...
}

var conf = config{
	Arguments: make([]string, 0, 0),
	Interval:  2,
	NoTitle:   false,
	Beep:      false,
	ErrExit:   false,
	Exec:      false,
	Help:      false,
	Version:   false}

func shellOutForNum(cmd string) (float64, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	// do literally any thing with the exit code
	if err != nil {
		return 0, errors.New("Exit Non Zero")
	}
	// TODO deal with weird people who use "," as a decimal.
	// TODO deal with "," as a thousands mark.
	// TODO don't use regex for this. Write a function.
	r := regexp.MustCompile("[\\d,\\.]+")
	cleaned := r.FindAllString(string(out), 1)
	if len(cleaned) == 0 {
		return 0, errors.New("NaN")
	}
	parsed, err := strconv.ParseFloat(cleaned[0], 64)
	if err != nil {
		return 0, errors.New("Parse Failure")
	}
	// If we have overran float64 or ~1.7*10^308
	if math.IsInf(parsed, 0) {
		return parsed, errors.New("Inf")
	}
	return parsed, nil
}

func warningdialog(msg string) ui.Bufferer {
	warn := ui.NewPar(msg)
	warn.Height = 4
	warn.Width = 34
	warn.Y = ui.TermHeight()/2 - warn.Height/2
	warn.X = ui.TermWidth()/2 - warn.Width/2
	warn.BorderLabel = "Warning"
	warn.BorderFg = ui.ColorYellow
	return warn
}

type timeSeries struct {
	Series []float64
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

func genXBasic(length int) []string {
	s := make([]string, length)
	for i := len(s) - 1; i >= 0; i-- {
		s[i] = strconv.Itoa(i - (len(s) - 1))
	}
	return s
}

func renderloop() {
	// build graph
	g := ui.NewLineChart()
	g.Height = ui.TermHeight()
	g.Width = ui.TermWidth()
	if conf.NoTitle != true {
		g.BorderLabel = "Every " + strconv.FormatFloat(conf.Interval, 'f', -1, 64) + "s: " + strings.Join(conf.Arguments, " ")
	}
	ts := timeSeries{}
	// render loop
	for {
		nextval, err := shellOutForNum(strings.Join(conf.Arguments, " "))
		if err != nil {
			if conf.Beep == true {
				fmt.Fprintf(os.Stdout, "\a")
			}
			if conf.ErrExit == true {
				ui.StopLoop()
				// report this error somehow
			}
			if err.Error() == "Exit Non Zero" {
				warn := warningdialog("Command exited non-zero\nCheck Escaping")
				ui.Render(g, warn)
			} else if err.Error() == "NaN" {
				warn := warningdialog("No numbers found in output\nCheck Escaping")
				ui.Render(g, warn)
			} else {
				warn := warningdialog("Unknown error occured")
				ui.Render(g, warn)
			}
		} else {
			ts.Series = append(ts.Series, nextval)
			if len(ts.Series) > g.GetCapacity() {
				g.Data = ts.Series[(len(ts.Series) - g.GetCapacity()):]
			} else {
				g.Data = ts.Series
			}
			// putting this in the loop deals with window changes.
			g.Width = ui.TermWidth()
			g.Height = ui.TermHeight()
			// render
			ui.Render(g)
		}
		time.Sleep(time.Millisecond * time.Duration(conf.Interval*1000))
	}
}

func init() {
	// watch like function
	getopt.FlagLong(&conf.Beep, "beep", 'b', "beep if command has a non-zero exit")
	getopt.FlagLong(&conf.Interval, "interval", 'n', "seconds to wait between updates")
	getopt.FlagLong(&conf.NoTitle, "no-title", 't', "turn off header")
	getopt.FlagLong(&conf.ErrExit, "errexit", 'e', "exit if command has a non-zero exit")
	getopt.FlagLong(&conf.Exec, "exec", 'x', "pass command to exec instead of \"sh -c\"")
	//meta
	getopt.FlagLong(&conf.Help, "help", 'h', "display this help and exit")
	getopt.FlagLong(&conf.Version, "version", 'v', "output version information and exit")
}

func main() {
	// parse arguments
	getopt.Parse()
	conf.Arguments = getopt.Args()
	// Version
	if conf.Version == true {
		fmt.Println("gwatch from https://github.com/robertely/gwatch 0.0.1")
		os.Exit(0)
	}

	// Help text
	if conf.Help == true {
		fmt.Println("graphing watch: expects numerical values, graphs the first one it sees.")
		fmt.Println("")
		getopt.Usage()
		os.Exit(0)
	}

	// no input handler
	if len(conf.Arguments) == 0 {
		getopt.Usage()
		os.Exit(1)
	}

	// Build UI
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer func() {
		ui.Close()
		fmt.Print("\033[2J") // Clear
	}()

	// Handle various keyboard exits.
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/C-x", func(ui.Event) {
		ui.StopLoop()
	})

	// Start rendering
	go renderloop()
	// Blocks and reacts to keyboard
	ui.Loop()

}
