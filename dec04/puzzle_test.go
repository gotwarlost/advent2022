package dec04

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestPuzzleP1(t *testing.T) {
	count := runPart1(testInput)
	assert.Equal(t, 2, count)
	count = runPart1(input)
	assert.Equal(t, 602, count)
}

func TestPuzzleP2(t *testing.T) {
	count := runPart2(testInput)
	assert.Equal(t, 4, count)
	count = runPart2(input)
	assert.Equal(t, 891, count)
}
