package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeSeriesAppend(t *testing.T) {
	ts := timeSeries{Capacity: 10}
	ts.append(float64(0))
	assert.Equal(t, float64(0), ts.Series[0])
	for i := 1; i <= 40; i++ {
		ts.append(float64(i))
	}
}
