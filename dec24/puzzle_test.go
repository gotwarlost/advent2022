package dec24

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.EqualValues(t, 18, runP1(testInput, 25))
	assert.EqualValues(t, 305, runP1(input, mainBudget))
}

func TestP2(t *testing.T) {
	assert.EqualValues(t, 54, runP2(testInput, 60))
	assert.EqualValues(t, 905, runP2(input, mainBudget))
}
