package dec22

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.EqualValues(t, 6032, runP1(testInput))
	assert.EqualValues(t, 95358, runP1(input))
}

func TestP2(t *testing.T) {
	assert.EqualValues(t, 5031, runP2(testInput, testRouteMap()))
	assert.EqualValues(t, 144361, runP2(input, mainRouteMap()))
}
