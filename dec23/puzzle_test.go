package dec23

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.EqualValues(t, 110, runP1(testInput))
	assert.EqualValues(t, 3882, runP1(input))
}

func TestP2(t *testing.T) {
	assert.EqualValues(t, 20, runP2(testInput))
	assert.EqualValues(t, -1, runP2(input))
}
