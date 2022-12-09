package dec9

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.Equal(t, 13, run(testInput, 2))
	assert.Equal(t, 5695, run(input, 2))
}

func TestP2(t *testing.T) {
	assert.Equal(t, 1, run(testInput, 10))
	assert.Equal(t, 2434, run(input, 10))
}
