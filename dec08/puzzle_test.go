package dec08

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test-input.txt
var testInput string

func TestAll(t *testing.T) {
	assert.Equal(t, 21, run(testInput))
	assert.Equal(t, 1733, run(input))

	assert.Equal(t, 8, run2(testInput))
	assert.Equal(t, 284648, run2(input))
}
