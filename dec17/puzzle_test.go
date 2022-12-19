package dec17

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

//go:embed test-input.txt
var testInput string

func TestP1(t *testing.T) {
	assert.Equal(t, int64(3068), run(testInput, 2022))
	assert.Equal(t, int64(3181), run(input, 2022))
}

func TestP2(t *testing.T) {
	assert.Equal(t, int64(1514285714288), run(testInput, gazillion))
	assert.Equal(t, int64(1570434782634), run(input, gazillion))
}
