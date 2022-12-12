package dec12

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	x := runP1(testInput)
	assert.Equal(t, 31, x)
	x = runP1(input)
	assert.Equal(t, 517, x)
}

func TestP2(t *testing.T) {
	x := runP2(testInput)
	assert.Equal(t, 29, x)
	x = runP2(input)
	assert.Equal(t, 512, x)
}
