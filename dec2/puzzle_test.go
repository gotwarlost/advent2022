package dec2

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestPuzzleP1(t *testing.T) {
	score := runPart1(testInput)
	assert.Equal(t, 15, score)
}

func TestPuzzleP2(t *testing.T) {
	score := runPart2(testInput)
	assert.Equal(t, 12, score)
}
