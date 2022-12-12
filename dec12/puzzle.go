package dec12

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
)

type position struct {
	row, col int
}

func (p position) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

type grid [][]int

func (g grid) rows() int {
	return len(g)
}

func (g grid) cols() int {
	return len(g[0])
}

func (g grid) value(pos position) int {
	return g[pos.row][pos.col]
}

func (g grid) String() string {
	var w bytes.Buffer
	for _, row := range g {
		for _, col := range row {
			x := fmt.Sprintf("%c ", byte(col))
			w.Write([]byte(x))
		}
		w.Write([]byte{'\n'})
	}
	return w.String()
}

type puzzleInput struct {
	grid       grid
	start, end position
}

//go:embed input.txt
var input string

func toGrid(in string) puzzleInput {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	var g grid
	var start, end position
	for i, row := range lines {
		cols := strings.Split(row, "")
		var colArray []int
		for j, c := range cols {
			val := c[0]
			if val == 'S' {
				start.row = i
				start.col = j
				val = 'a'
			}
			if val == 'E' {
				end.row = i
				end.col = j
				val = 'z'
			}
			colArray = append(colArray, int(val))
		}
		g = append(g, colArray)
	}
	return puzzleInput{
		grid:  g,
		start: start,
		end:   end,
	}
}

func possibleReverseSteps(g grid, current position) []position {
	val := g.value(current)
	var ret []position

	add := func(p position) {
		v2 := g.value(p)
		if val-v2 <= 1 {
			ret = append(ret, p)
		}
	}
	if current.row != 0 {
		add(position{current.row - 1, current.col})
	}
	if current.row != g.rows()-1 {
		add(position{current.row + 1, current.col})
	}
	if current.col != 0 {
		add(position{current.row, current.col - 1})
	}
	if current.col != g.cols()-1 {
		add(position{current.row, current.col + 1})
	}
	return ret
}

func computePositionScores(g grid, current position, scores map[position]int, score int) {
	// check if there is a path with an equal or lower score and abandon
	if prevScore, ok := scores[current]; ok {
		if prevScore <= score {
			return
		}
	}
	// set the current score
	scores[current] = score

	// compute for reachable steps
	possibles := possibleReverseSteps(g, current)
	for _, x := range possibles {
		computePositionScores(g, x, scores, score+1)
	}
}

func computeScores(in string) (puzzleInput, map[position]int) {
	pz := toGrid(in)
	scores := map[position]int{}
	computePositionScores(pz.grid, pz.end, scores, 0)
	return pz, scores
}

func runP1(in string) int {
	pz, scores := computeScores(in)
	return scores[pz.start]
}

func runP2(in string) int {
	pz, scores := computeScores(in)
	min := 100000
	for i, row := range pz.grid {
		for j := range row {
			p := position{i, j}
			if pz.grid.value(p) != int('a') {
				continue
			}
			if _, ok := scores[p]; !ok { // unreachable, no score set
				continue
			}
			if scores[p] < min {
				min = scores[p]
			}
		}
	}
	return min
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
