package dec6

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

func TestPuzzleP1(t *testing.T) {
	tests := []struct {
		input    string
		expect4  int
		expect14 int
	}{
		{
			input:    "mjqjpqmgbljsphdztnvjfqwrcgsmlb",
			expect4:  7,
			expect14: 19,
		},
		{
			input:    "bvwbjplbgvbhsrlpgdmjqwftvncz",
			expect4:  5,
			expect14: 23,
		},
		{
			input:    "nppdvjthqldpwncqszvftbrmjlhg",
			expect4:  6,
			expect14: 23,
		},
		{
			input:    "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg",
			expect4:  10,
			expect14: 29,
		},
		{
			input:    "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
			expect4:  11,
			expect14: 26,
		},
	}
	for _, test := range tests {
		t.Run(test.input[:5]+"-4", func(t *testing.T) {
			assert.Equal(t, test.expect4, nonRepeatingChars(test.input, 4))
		})
		t.Run(test.input[:5]+"-14", func(t *testing.T) {
			assert.Equal(t, test.expect14, nonRepeatingChars(test.input, 14))
		})
	}
}
