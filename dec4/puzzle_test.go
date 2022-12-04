package dec4

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
}

func TestPuzzleP2(t *testing.T) {
	count := runPart2(testInput)
	assert.Equal(t, 4, count)
}
