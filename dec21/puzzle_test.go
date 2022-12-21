package dec21

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.EqualValues(t, 152, runP1(testInput))
	assert.EqualValues(t, 276156919469632, runP1(input))
}

func TestP2(t *testing.T) {
	assert.EqualValues(t, 301, runP2(testInput))
	assert.EqualValues(t, 3441198826073, runP2(input))
}
