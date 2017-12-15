package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncNormal(t *testing.T) {
	s := truncate("Concurrency is not parallelism.", 12)
	assert.Equal(t, "parallelism.", s)
}

func TestTruncNegative(t *testing.T) {
	s := truncate("Concurrency is not parallelism.", -10)
	assert.Equal(t, "", s)
}

func TestTruncWide(t *testing.T) {
	s := truncate("Concurrency is not parallelism.", 100)
	assert.Equal(t, "Concurrency is not parallelism.", s)
}

func TestParseOutSingleNormalNegative(t *testing.T) {
	s, _ := parseOutSingle("Concurrency is -1,223.234   .4not parallelism.")
	assert.Equal(t, -1223.234, s)
}

func TestParseOutSingleNormal(t *testing.T) {
	s, _ := parseOutSingle("Concurrency is 1,223.234   .4not parallelism.")
	assert.Equal(t, 1223.234, s)
}

func TestParseOutSingleBadNumeral(t *testing.T) {
	s, err := parseOutSingle("Concurrency is 1,2.23.234   .4not parallelism.")
	assert.Equal(t, float64(0), s)
	assert.EqualError(t, err, "strconv.ParseFloat: parsing \"12.23.234\": invalid syntax")
}

func TestParseOutSingleNoNumeral(t *testing.T) {
	s, err := parseOutSingle("Concurrency is not parallelism")
	assert.Equal(t, float64(0), s)
	assert.EqualError(t, err, "No numeral found")
}

func TestParseOutSingleDanglingDelimiter(t *testing.T) {
	s, err := parseOutSingle("Concurrency is not parallelism.")
	assert.Equal(t, float64(0), s)
	// I would rather this be the error, need to tweak the regex
	// assert.EqualError(t, err, "No numeral found")
	assert.EqualError(t, err, "strconv.ParseFloat: parsing \".\": invalid syntax")
}
