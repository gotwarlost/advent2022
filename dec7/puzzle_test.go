package dec7

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test-input.txt
var testInput string

func TestAll(t *testing.T) {
	p1, p2 := run(testInput, 100000)
	assert.Equal(t, 95437, p1)
	assert.Equal(t, 24933642, p2)
}
