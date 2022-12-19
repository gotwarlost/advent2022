package dec18

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.Equal(t, 64, run(testInput, false))
	assert.Equal(t, 3448, run(input, false))
}

func TestP2(t *testing.T) {
	assert.Equal(t, 58, run(testInput, true))
	assert.Equal(t, 2052, run(input, true))
}
