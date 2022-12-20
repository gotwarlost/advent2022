package dec20

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.EqualValues(t, 3, run(testInput, 1, 1))
	assert.EqualValues(t, 3346, run(input, 1, 1))
}

func TestP2(t *testing.T) {
	assert.EqualValues(t, 1623178306, run(testInput, 811589153, 10))
	assert.EqualValues(t, 4265712588168, run(input, 811589153, 10))
}
