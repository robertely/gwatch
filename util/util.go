package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat"
)

// Truncate a string to ln length
func Truncate(st string, ln int) string {
	if ln <= 0 {
		return ""
	}
	if len(st) > ln {
		return st[len(st)-ln:]
	}
	return st
}

// ParseOutSingle returns float from parsed string.
// Isolates the fist thing that looks like a number.
// Hard coded maximim of 1024 characters
func ParseOutSingle(s string) (float64, error) {
	r := regexp.MustCompile("(-?[\\d,\\.]+)")
	isolated := r.FindString(s)
	if len(isolated) == 0 {
		return 0, errors.New("No numeral found")
	}
	parsed, err := strconv.ParseFloat(strings.Replace(isolated, ",", "", -1), 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// GetMax Gets maximum value in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func GetMax(series []float64) (max float64) {
	for _, i := range series {
		if i > max {
			max = i
		}
	}
	return
}

// GetMin Gets minumum value in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func GetMin(series []float64) (min float64) {
	min = series[0]
	for _, i := range series {
		if i < min {
			min = i
		}
	}
	return
}

// GetAvg Gets simple average for in range of ts
// rng (range) allows you to work only with the data you are graphing and not the full capacity.
func GetAvg(series []float64) float64 {
	total := float64(0)
	for _, i := range series {
		total = total + i
	}
	return total / float64(len(series))
}

// GetStD Gets standard deviation in range of series
func GetStD(series []float64) (stdev float64) {
	return stat.StdDev(series, nil)
}

// LineCount counts the nubmer of lines in a string (split on \n)
func LineCount(s string) int {
	return len(strings.Split(strings.TrimSuffix(s, "\n"), "\n"))
}

// MaxLineLength take s string and returns the length of the largest single line.
func MaxLineLength(s string) int {
	i := 0
	for _, line := range strings.Split(strings.TrimSuffix(s, "\n"), "\n") {
		if len(line) > i {
			i = len(line)
		}
	}
	return i
}
