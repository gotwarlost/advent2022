package dec11

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.EqualValues(t, 10605, runP1(testInput))
	// assert.Equal(t, 316888, runP1(input))
}

func TestP2(t *testing.T) {
	assert.EqualValues(t, 2713310158, runP2(testInput))
	// assert.Equal(t, 316888, runP2(input))
}
