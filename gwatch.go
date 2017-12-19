package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	util "github.com/robertely/gwatch/util"

	getopt "github.com/pborman/getopt/v2"
	ui "github.com/robertely/termui"
)

// config stores run time options
type config struct {
	Arguments []string
	// Watch flags were going to mock
	Exec     bool
	Interval float64
	NoTitle  bool
	Beep     bool
	ErrExit  bool
	Help     bool
	Version  bool
	// differences bool // N/A
	// precise bool // Naw not going to do this.
	// chgexit bool // N/A
	// color bool // N/A

	// gwatch specific flags...
	ShowStats   bool
	TruncFloats bool
	ExecTimeOut int
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

// Validate Validates config as populated by getopt
func (c *config) Validate() {
	// Version
	if c.Version == true {
		// TODO version managed by Makefile
		fmt.Println("gwatch from https://github.com/robertely/gwatch 0.0.3")
		os.Exit(0)
	}

	// Help text
	if c.Help == true {
		fmt.Println("graphing watch: expects numerical values, graphs the first one it sees.")
		fmt.Println("")
		getopt.Usage()
		os.Exit(0)
	}

	// no input handler
	if len(c.Arguments) == 0 {
		getopt.Usage()
		os.Exit(1)
	}
}

func shellOutForNum(args []string) (float64, error) {
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	if conf.Exec == true {
		cmd = exec.Command(args[0], args[1:]...)
	}
	out, err := cmd.Output()
	// exitCode := what....?
	//TODO do literally any thing with the exit code [exec.Command().ProcessState????]
	if err != nil {
		return 0, errors.New("Exit Non Zero")
	}
	parsed, _ := util.ParseOutSingle(string(out))
	// TODO: check for "strconv.ParseFloat"
	// TODO: check for "No numeral found"
	// If we have overran float64 or ~1.7*10^308
	if math.IsInf(parsed, 0) {
		return parsed, errors.New("Inf")
	}
	return parsed, nil
}

// warningdialog returns a ui.Bufferer displaying msg.
// The intended use is as a popover.
func warningdialog(msg string) ui.Bufferer {
	warn := ui.NewPar(msg)
	warn.Height = util.LineCount(msg) + 2    // 2 is room for boarder
	warn.Width = util.MaxLineLength(msg) + 3 // 3 adds a little padding.
	warn.Y = ui.TermHeight()/2 - warn.Height/2
	warn.X = ui.TermWidth()/2 - warn.Width/2
	warn.BorderLabel = "Warning"
	warn.BorderFg = ui.ColorYellow
	return warn
}

// timeSeries thinly wraps []float64 adding a hard limit(Capacity.)
type timeSeries struct {
	Series   []float64
	Capacity int
}

// append adds to ts.Series and truncates when at ts.Capacity
func (ts *timeSeries) append(next float64) {
	if len(ts.Series) < ts.Capacity {
		ts.Series = append(ts.Series, next)
	} else {
		ts.Series = append(ts.Series[1:], next)
	}
}

// getMax Gets maximum value in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func (ts *timeSeries) getMax(rng int) (max float64) {
	for _, i := range ts.Series {
		if i > max {
			max = i
		}
	}
	return
}

// getMin Gets minumum value in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func (ts *timeSeries) getMin(rng int) (min float64) {
	min = ts.Series[0]
	for _, i := range ts.Series {
		if i < min {
			min = i
		}
	}
	return
}

// getAvg Gets simple average for in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func (ts *timeSeries) getAvg(rng int) (avg float64) {
	return
}

// getStDev Gets standard deviation in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func (ts *timeSeries) getStDev(rng int) (stdev float64) {
	// TODO: This.
	return
}

// genXBasic generates xaxis labels.
// This is resposnable for that weird flip you see when the graph fills.
// I'm working around termui limitations here. Fixing this in termiu would be an undertaking.
// If I was inclined to go that far I would probably write my own graphing library.
func genXBasic(length int, inverse bool) []string {
	s := make([]string, length)
	if inverse {
		for i := len(s) - 1; i >= 0; i-- {
			s[i] = strconv.Itoa(i - (len(s) - 1))
		}
	} else {
		for i := 0; i <= len(s)-1; i++ {
			s[i] = strconv.Itoa(i)
		}
	}
	return s
}

func renderloop(g *ui.LineChart) {
	ts := timeSeries{Capacity: 125000} // 1MB/64bits "reasonable maximum"
	// render loop
	for {
		// putting this in the loop helps to deal with window changes.
		g.Width = ui.TermWidth()
		g.Height = ui.TermHeight()
		nextval, err := shellOutForNum(conf.Arguments)
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
			//append next value
			ts.append(nextval)
			// merge ts.series with g.Data and recalculate DataLabels
			if len(ts.Series) >= g.GetCapacity() {
				g.DataLabels = genXBasic(g.GetCapacity(), true)
				g.Data = ts.Series[(len(ts.Series) - g.GetCapacity()):]
			} else {
				g.DataLabels = genXBasic(g.GetCapacity(), false)
				g.Data = ts.Series
			}
			// Render
			ui.Clear()
			ui.Render(g)
		}
		// Sleep
		time.Sleep(time.Millisecond * time.Duration(conf.Interval*1000))
	}
}

func init() {
	// watch like function
	getopt.FlagLong(&conf.Beep, "beep", 'b', "beep if command has a non-zero exit")
	getopt.FlagLong(&conf.Interval, "interval", 'n', "seconds to wait between updates")
	getopt.FlagLong(&conf.NoTitle, "no-title", 't', "turn off header")
	getopt.FlagLong(&conf.ErrExit, "errexit", 'e', "exit if command has a non-zero exit")
	getopt.FlagLong(&conf.Exec, "exec", 'x', "directly pass command to os.exec instead of \"sh -c\"")
	//meta
	getopt.FlagLong(&conf.Help, "help", 'h', "display this help and exit")
	getopt.FlagLong(&conf.Version, "version", 'v', "output version information and exit")
}

func main() {
	// parse arguments
	getopt.Parse()
	conf.Arguments = getopt.Args()
	conf.Validate()

	// Build UI
	if err := ui.Init(); err != nil {
		ui.Close()
		panic(err)
	}

	// Clean up. Not calling this leaves your terminal in a bad state.
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

	// Create graph
	g := ui.NewLineChart()

	// Set title
	if conf.NoTitle != true {
		g.BorderLabel = "Every " + strconv.FormatFloat(conf.Interval, 'f', -1, 64) + "s: " + strings.Join(conf.Arguments, " ")
	}

	// Handle resize. I strongly suspect this isn't working perfectly in OSX
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		g.Width = ui.TermWidth()
		g.Height = ui.TermHeight()
		ui.Render(g)
	})

	// Start rendering
	go renderloop(g)

	// Blocks and reacts to keyboard
	ui.Loop()
}
