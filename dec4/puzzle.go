package dec4

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type assigment struct {
	min int
	max int
}

func (a assigment) contains(other assigment) bool {
	return a.min <= other.min && a.max >= other.max
}

func (a assigment) overlaps(other assigment) bool {
	if a.contains(other) {
		return true
	}
	return (a.min >= other.min && a.min <= other.max) ||
		(a.max >= other.min && a.max <= other.max)
}

func parsePart(s string) assigment {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		panic(fmt.Errorf("invalid part: %q", s))
	}
	min, err1 := strconv.Atoi(parts[0])
	max, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		panic(fmt.Errorf("invalid parts: %q", s))
	}
	return assigment{
		min: min,
		max: max,
	}
}

func runPart1(in string) int {
	count := 0
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid part: %q", line))
		}
		a1 := parsePart(parts[0])
		a2 := parsePart(parts[1])
		if a1.contains(a2) || a2.contains(a1) {
			count++
		}
	}
	return count
}

func runPart2(in string) int {
	count := 0
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid part: %q", line))
		}
		a1 := parsePart(parts[0])
		a2 := parsePart(parts[1])
		if a1.overlaps(a2) || a2.overlaps(a1) {
			count++
		}
	}
	return count
}

func RunPart1() {
	fmt.Println("Output:", runPart1(input))
}

func RunPart2() {
	fmt.Println("Output:", runPart2(input))
}
