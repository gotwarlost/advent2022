package dec6

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

//go:embed test-input-part2.txt
var testInput2 string

func TestPuzzleP1(t *testing.T) {
	assert.Equal(t, 11, run(testInput, 4))
}

func TestPuzzleP2(t *testing.T) {
	assert.Equal(t, 26, run(testInput2, 14))
}
