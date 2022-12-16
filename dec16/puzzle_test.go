package dec16

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.Equal(t, 1651, runP1(testInput))
	// assert.Equal(t, 0, runP1(input))
}

func TestP2(t *testing.T) {
	// assert.Equal(t, 0, runP2(testInput))
	// assert.Equal(t, 11840879211051, runP2(input))
}
