package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	getopt "github.com/pborman/getopt/v2"
)

type config struct {
	Arguments []string
	// Watch flags were going to mock
	// differences bool // not sure what I could do with this.
	// precise bool // Naw not going to do this.
	// chgexit bool // N/A here
	// color bool // N/A here	// precise bool // Naw not going to do this.
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

func shellOutNum(cmd string) (float64, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return 0, errors.New("Command could not be executed")
		// do something more useful here...
	}
	r := regexp.MustCompile("[\\d,\\.]+")
	cleaned := r.FindAllString(string(out), 1)
	if len(cleaned) == 0 {
		return 0, errors.New("no numerical output detected")
	}
	parsed, _ := strconv.ParseFloat(cleaned[0], 64)
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

	// build graph
	x := timeSeries{Series: []float64{}}
	g := ui.NewLineChart()
	g.Data = x.Series
	g.Height = ui.TermHeight()
	g.Width = ui.TermWidth()
	g.BorderLabel = "Every " + strconv.FormatFloat(conf.Interval, 'f', -1, 64) + "s: " + strings.Join(conf.Arguments, " ")

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

	go func() {
		// none of this math is right. How do you get the capacity of g.Data ???
		for {
			nextval, err := shellOutNum(strings.Join(conf.Arguments, " "))
			if err != nil {
				warn := warningdialog(err.Error())
				ui.Render(g, warn)
				// not a good idea make this an else
				time.Sleep(time.Millisecond * time.Duration(conf.Interval*1000))
				continue
			}

			x.Series = append(x.Series, nextval)

			if len(x.Series) > ui.TermWidth()*2 { // Brail is 2 wide
				g.Data = x.Series[len(x.Series)-ui.TermWidth()*2:]
			} else {
				g.Data = x.Series
			}
			g.Width = ui.TermWidth()
			g.Height = ui.TermHeight()
			// g.DataLabels = []string{}
			ui.Render(g)
			time.Sleep(time.Millisecond * time.Duration(conf.Interval*1000))
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
