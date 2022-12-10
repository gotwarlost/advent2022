package dec10

import (
	"bytes"
	_ "embed"
	"strings"
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
	assert.EqualValues(t, 3960, s.multAt(220))
	assert.Equal(t, 13140, s.multAt(20)+s.multAt(60)+s.multAt(100)+s.multAt(140)+s.multAt(180)+s.multAt(220))

	s = getStrengths(input)
	assert.Equal(t, 17940, s.multAt(20)+s.multAt(60)+s.multAt(100)+s.multAt(140)+s.multAt(180)+s.multAt(220))
}

func TestP2(t *testing.T) {
	var b bytes.Buffer
	printCRT(input, &b)
	expected := `
####..##..###...##....##.####...##.####.
...#.#..#.#..#.#..#....#.#.......#....#.
..#..#....###..#..#....#.###.....#...#..
.#...#....#..#.####....#.#.......#..#...
#....#..#.#..#.#..#.#..#.#....#..#.#....
####..##..###..#..#..##..#.....##..####.
`
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(b.String()))
}
