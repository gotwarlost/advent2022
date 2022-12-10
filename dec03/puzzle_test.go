package dec03

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestPuzzleP1(t *testing.T) {
	out := runPart1(testInput)
	assert.Equal(t, 157, out)
	out = runPart1(input)
	assert.Equal(t, 8088, out)
}

func TestPuzzleP2(t *testing.T) {
	out := runPart2(testInput)
	assert.Equal(t, 70, out)
	out = runPart2(input)
	assert.Equal(t, 2522, out)
}
