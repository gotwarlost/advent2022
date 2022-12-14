package dec14

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.Equal(t, 24, runP1(testInput))
	assert.Equal(t, 6072, runP1(input))
}

func TestP2(t *testing.T) {
	assert.Equal(t, 93, runP2(testInput))
	assert.Equal(t, 26729, runP2(input))
}
