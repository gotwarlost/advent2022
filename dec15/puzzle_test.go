package dec15

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.Equal(t, 26, runP1(testInput, 10))
	assert.Equal(t, 5394423, runP1(input, 2000000))
}

func TestP2(t *testing.T) {
	assert.Equal(t, 56000011, runP2(testInput))
	assert.Equal(t, 11840879211051, runP2(input))
}
