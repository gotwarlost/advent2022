package dec03

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

func toSet(word string) map[byte]bool {
	ret := map[byte]bool{}
	for i := 0; i < len(word); i++ {
		ret[word[i]] = true
	}
	return ret
}

func intersection(s1, s2 map[byte]bool) map[byte]bool {
	out := map[byte]bool{}
	for k := range s1 {
		if s2[k] {
			out[k] = true
		}
	}
	if len(out) == 0 {
		panic("no intersection found")
	}
	return out
}

func extractSingleValue(s map[byte]bool) byte {
	if len(s) > 1 {
		panic("more than one element found in intersection")
	}
	var out byte
	for k := range s {
		out = k
	}
	return out
}

func score(s byte) int {
	if s >= 'a' && s <= 'z' {
		return 1 + int(s-'a')
	}
	return 27 + int(s-'A')
}

func runPart1(in string) int {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	output := 0
	for _, line := range lines {
		l := len(line)
		if l%2 != 0 {
			panic(fmt.Errorf("line %q has odd number of characters", l))
		}
		first := line[:l/2]
		second := line[l/2:]
		s1 := toSet(first)
		s2 := toSet(second)
		output += score(extractSingleValue(intersection(s1, s2)))
	}
	return output
}

func runPart2(in string) int {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	output := 0
	if len(lines)%3 != 0 {
		panic(fmt.Errorf("number of entries not a multiple of 3, got: %d", len(lines)))
	}
	for i := 0; i < len(lines); i += 3 {
		s1 := toSet(lines[i])
		s2 := toSet(lines[i+1])
		s3 := toSet(lines[i+2])
		output += score(extractSingleValue(intersection(intersection(s1, s2), s3)))
	}
	return output
}

func RunPart1() {
	fmt.Println("Score:", runPart1(input))
}

func RunPart2() {
	fmt.Println("Score:", runPart2(input))
}
