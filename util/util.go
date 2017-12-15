package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
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
