package dec6

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

func TestPuzzleP1(t *testing.T) {
	assert.Equal(t, 7, nonRepeatingChars("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 4))
	assert.Equal(t, 5, nonRepeatingChars("bvwbjplbgvbhsrlpgdmjqwftvncz", 4))
	assert.Equal(t, 6, nonRepeatingChars("nppdvjthqldpwncqszvftbrmjlhg", 4))
	assert.Equal(t, 10, nonRepeatingChars("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 4))
	assert.Equal(t, 11, nonRepeatingChars("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 4))
}

func TestPuzzleP2(t *testing.T) {
	assert.Equal(t, 19, nonRepeatingChars("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 14))
	assert.Equal(t, 23, nonRepeatingChars("bvwbjplbgvbhsrlpgdmjqwftvncz", 14))
	assert.Equal(t, 23, nonRepeatingChars("nppdvjthqldpwncqszvftbrmjlhg", 14))
	assert.Equal(t, 29, nonRepeatingChars("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 14))
	assert.Equal(t, 26, nonRepeatingChars("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 14))
}
