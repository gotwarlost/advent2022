package dec01

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test-input.txt
var testInput string

func TestTop(t *testing.T) {
	t1, t3 := run(testInput)
	assert.Equal(t, 24000, t1)
	assert.Equal(t, 45000, t3)
	t1, t3 = run(actualInput)
	assert.Equal(t, 75622, t1)
	assert.Equal(t, 213159, t3)
}
