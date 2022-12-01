package dec1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTop(t *testing.T) {
	t1, t3 := run(testInput)
	assert.Equal(t, 24000, t1)
	assert.Equal(t, 45000, t3)
}
