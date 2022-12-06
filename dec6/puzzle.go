package dec6

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

func hasDistinctChars(s string) bool {
	m := map[byte]bool{}
	for _, b := range []byte(s) {
		m[b] = true
	}
	return len(m) == len(s)
}

func nonrepeatingChars(in string, seq int) int {
	in = strings.TrimSpace(in)
	for i := 0; i < len(in)-seq-1; i++ {
		s := in[i : i+seq]
		if hasDistinctChars(s) {
			return i + seq
		}
	}
	panic("not found")
}

func RunPart1() {
	fmt.Println("Output:", nonrepeatingChars(input, 4))
}

func RunPart2() {
	fmt.Println("Output:", nonrepeatingChars(input, 14))
}
