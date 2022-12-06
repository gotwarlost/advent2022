package dec6

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func run(in string, seq int) int {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	var last int
	for n, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		var input string
		var expected int
		var err error
		if len(parts) == 2 {
			input = parts[0]
			expected, err = strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
		} else {
			input = parts[0]
			expected = -1
		}
		out := -1
		for i := 0; i < len(input)-seq-1; i++ {
			s := input[i : i+seq]
			m := map[byte]bool{}
			for _, b := range []byte(s) {
				m[b] = true
			}
			if len(m) == seq {
				out = i + seq
				break
			}
		}
		if out == -1 {
			panic("not found")
		}
		if expected != -1 && out != expected {
			panic(fmt.Errorf("want %d, got %d", expected, out))
		}
		log.Println("line", n+1, ":", out)
		last = out
	}
	return last
}

func RunPart1() {
	fmt.Println("Output:", run(input, 4))
}

func RunPart2() {
	fmt.Println("Output:", run(input, 14))
}
