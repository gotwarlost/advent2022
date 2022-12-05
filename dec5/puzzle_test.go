package dec5

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
	assert.Equal(t, "CMZ", out)
	out = runPart1(input)
	assert.Equal(t, "LBLVVTVLP", out)
}

func TestPuzzleP2(t *testing.T) {
	out := runPart2(testInput)
	assert.Equal(t, "MCD", out)
	out = runPart2(input)
	assert.Equal(t, "TPFFBDRJD", out)
}
