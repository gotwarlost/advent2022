package dec6

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

func TestPuzzleP1(t *testing.T) {
	assert.Equal(t, 7, nonrepeatingChars("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 4))
	assert.Equal(t, 5, nonrepeatingChars("bvwbjplbgvbhsrlpgdmjqwftvncz", 4))
	assert.Equal(t, 6, nonrepeatingChars("nppdvjthqldpwncqszvftbrmjlhg", 4))
	assert.Equal(t, 10, nonrepeatingChars("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 4))
	assert.Equal(t, 11, nonrepeatingChars("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 4))
}

func TestPuzzleP2(t *testing.T) {
	assert.Equal(t, 19, nonrepeatingChars("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 14))
	assert.Equal(t, 23, nonrepeatingChars("bvwbjplbgvbhsrlpgdmjqwftvncz", 14))
	assert.Equal(t, 23, nonrepeatingChars("nppdvjthqldpwncqszvftbrmjlhg", 14))
	assert.Equal(t, 29, nonrepeatingChars("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 14))
	assert.Equal(t, 26, nonrepeatingChars("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 14))
}
