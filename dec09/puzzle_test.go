package dec09

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test-input.txt
var testInput string

//go:embed test-input2.txt
var testInput2 string

func TestP1(t *testing.T) {
	assert.Equal(t, 13, countVisited(testInput, 2))
	assert.Equal(t, 5695, countVisited(input, 2))
}

func TestP2(t *testing.T) {
	assert.Equal(t, 1, countVisited(testInput, 10))
	assert.Equal(t, 36, countVisited(testInput2, 10))
	assert.Equal(t, 2434, countVisited(input, 10))
}
