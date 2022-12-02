package dec2

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type thing int

const (
	_ thing = iota
	rock
	paper
	scissors
)

func (t thing) score() int {
	switch t {
	case rock:
		return 1
	case paper:
		return 2
	case scissors:
		return 3
	default:
		panic(fmt.Errorf("invalid thing %d", t))
	}
}

func (t thing) String() string {
	switch t {
	case rock:
		return "rock"
	case paper:
		return "paper"
	case scissors:
		return "scissors"
	}
	return "unknown"
}

type outcome int

const (
	_ outcome = iota
	win
	loss
	draw
)

func (o outcome) score() int {
	switch o {
	case win:
		return 6
	case loss:
		return 0
	case draw:
		return 3
	default:
		panic(fmt.Errorf("invalid outcome %d", o))
	}
}

func (o outcome) String() string {
	switch o {
	case win:
		return "win"
	case loss:
		return "loss"
	case draw:
		return "draw"
	}
	return "unknown"
}

var wins = map[thing]thing{
	rock:     scissors,
	scissors: paper,
	paper:    rock,
}

var loses = map[thing]thing{}

func init() {
	for k, v := range wins {
		loses[v] = k
	}
}

func (t thing) outcome(other thing) outcome {
	if t == other {
		return draw
	}
	ret := loss
	if wins[t] == other {
		ret = win
	}
	return ret
}

func (t thing) opponentFromOutcome(o outcome) thing {
	switch o {
	case draw:
		return t
	case win:
		return loses[t]
	case loss:
		return wins[t]
	default:
		panic(fmt.Errorf("invalid outcome %d", o))
	}
}

var opponentMap = map[string]thing{
	"A": rock,
	"B": paper,
	"C": scissors,
}

var myMap = map[string]thing{
	"X": rock,
	"Y": paper,
	"Z": scissors,
}

var outcomeMap = map[string]outcome{
	"X": loss,
	"Y": draw,
	"Z": win,
}

type turn struct {
	opponent string
	second   string
}

func RunPart1() {
	fmt.Println(runPart1(input))
}

func toTurns(in string) []turn {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	var turns []turn
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid line: %q", line))
		}
		turns = append(turns, turn{opponent: parts[0], second: parts[1]})
	}
	return turns
}

func runPart1(in string) int {
	turns := toTurns(in)
	var score int
	for _, t := range turns {
		theirs := opponentMap[t.opponent]
		ours := myMap[t.second]
		if theirs == 0 || ours == 0 {
			panic(fmt.Errorf("problem: them %q, me %q", t.opponent, t.second))
		}
		o := ours.outcome(theirs)
		score += o.score() + ours.score()
	}
	return score
}

func runPart2(in string) int {
	turns := toTurns(in)
	var score int
	for _, t := range turns {
		theirs := opponentMap[t.opponent]
		o := outcomeMap[t.second]
		if theirs == 0 || o == 0 {
			panic(fmt.Errorf("problem: them %q, outcome %q", t.opponent, t.second))
		}
		us := theirs.opponentFromOutcome(o)
		score += o.score() + us.score()
	}
	return score
}

func RunPart2() {
	fmt.Println(runPart2(input))
}
