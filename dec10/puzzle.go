package dec10

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

type strengths []int

func (s strengths) multAt(n int) int {
	return n * s[n-1]
}

func (s strengths) valueAt(n int) int {
	if n == 0 {
		return 1
	}
	return s[n-1]
}

type signal struct {
	xInc int
}

//go:embed input.txt
var input string

func toSignals(in string) []signal {
	var sigs []signal
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		switch line {
		case "noop":
			sigs = append(sigs, signal{})
		default:
			parts := strings.Split(line, " ")
			if len(parts) != 2 || parts[0] != "addx" {
				panic(line)
			}
			n, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			sigs = append(sigs, signal{}, signal{xInc: n})
		}
	}
	return sigs
}

func run(in string) strengths {
	sigs := toSignals(in)
	x := 1
	var afterValues []int
	for _, cycle := range sigs {
		x += cycle.xInc
		afterValues = append(afterValues, x)
	}
	return afterValues
}

func run2(in string) {
	ss := run(in)
	cycle := 0
	for row := 0; row < 6; row++ {
		for col := 0; col < 40; col++ {
			s := ss.valueAt(cycle)
			min := s - 1
			max := s + 1
			if col >= min && col <= max {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
			cycle++
		}
		fmt.Println("")
	}
}

func RunP1() {
	fmt.Println(run(input))
}

func RunP2() {
	run2(input)
}
