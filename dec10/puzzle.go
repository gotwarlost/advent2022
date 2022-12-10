package dec10

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

type strengths []int

func (s strengths) multAt(n int) int {
	return n * s.valueAt(n)
}

func (s strengths) valueAt(n int) int {
	return s[n]
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

func getStrengths(in string) strengths {
	sigs := toSignals(in)
	x := 1
	beforeValues := []int{1}
	for _, cycle := range sigs {
		x += cycle.xInc
		beforeValues = append(beforeValues, x)
	}
	return beforeValues
}

func printCRT(in string) {
	ss := getStrengths(in)
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
	s := getStrengths(input)
	fmt.Println(s.multAt(20) + s.multAt(60) + s.multAt(100) + s.multAt(140) + s.multAt(180) + s.multAt(220))
}

func RunP2() {
	printCRT(input)
}
