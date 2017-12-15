package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func truncate(st string, ln int) string {
	if ln <= 0 {
		return ""
	}
	if len(st) > ln {
		return st[len(st)-ln:]
	}
	return st
}

func parseOutSingle(s string) (float64, error) {
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

// lineCount counts the nubmer of lines in a string (split on \n)
func lineCount(s string) int {
	return len(strings.Split(strings.TrimSuffix(s, "\n"), "\n"))
}

// maxLineLength take s string and returns the length of the largest single line.
func maxLineLength(s string) int {
	i := 0
	for _, line := range strings.Split(strings.TrimSuffix(s, "\n"), "\n") {
		if len(line) > i {
			i = len(line)
		}
	}
	return i
}
