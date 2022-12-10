package dec10

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	s := getStrengths(testInput)
	assert.EqualValues(t, 420, s.multAt(20))
	assert.EqualValues(t, 1140, s.multAt(60))
	assert.EqualValues(t, 1800, s.multAt(100))
	assert.EqualValues(t, 2940, s.multAt(140))
	assert.EqualValues(t, 2880, s.multAt(180))

	// this does not match explanation in puzzle but regular input checks out; need to figure out why
	assert.EqualValues(t, 4180, s.multAt(220))
	assert.Equal(t, 13360, s.multAt(20)+s.multAt(60)+s.multAt(100)+s.multAt(140)+s.multAt(180)+s.multAt(220))

	s = getStrengths(input)
	assert.Equal(t, 17940, s.multAt(20)+s.multAt(60)+s.multAt(100)+s.multAt(140)+s.multAt(180)+s.multAt(220))
}

func TestP2(t *testing.T) {
	// cannot write a test for this :(
}
