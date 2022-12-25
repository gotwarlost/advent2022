package dec25

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestSnafuToDecimal(t *testing.T) {
	tests := []struct {
		snafu    string
		expected int
	}{
		{"1=-0-2", 1747},
		{"12111", 906},
		{"2=0=", 198},
		{"21", 11},
		{"2=01", 201},
		{"111", 31},
		{"20012", 1257},
		{"112", 32},
		{"1=-1=", 353},
		{"1-12", 107},
		{"12", 7},
		{"1=", 3},
		{"122", 37},
		{"1=====--10201111-0", 381535316395},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i+1), func(t *testing.T) {
			assert.EqualValues(t, test.expected, snafuToDecimal(test.snafu))
			assert.EqualValues(t, test.snafu, decimalToSnafu(int64(test.expected)))
		})
	}
}

func TestP1(t *testing.T) {
	assert.EqualValues(t, "2=-1=0", runP1(testInput))
	assert.EqualValues(t, "2-00=12=21-0=01--000", runP1(input))
}

func TestP2(t *testing.T) {
	// assert.EqualValues(t, 20, runP2(testInput))
	// assert.EqualValues(t, -1, runP2(input))
}
