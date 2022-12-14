package dec14

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

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
	// log.Println("V", pos)
	return g[pos.row][pos.col]
}

func (g grid) fill(pos position) {
	// log.Println("fill:", pos)
	if g[pos.row][pos.col] == 1 {
		panic(fmt.Errorf("already filled: %v", pos))
	}
	g[pos.row][pos.col] = 1
}

func newGrid(rows, cols int) grid {
	var ret [][]int
	for i := 0; i < rows; i++ {
		row := make([]int, cols)
		ret = append(ret, row)
	}
	return ret
}

func toPositions(in string) *puzzleInput {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	minRow := 0 // always 0
	maxRow := 0
	minCol := 1000000
	maxCol := 0

	var ret [][]position

	for _, line := range lines {
		var linePos []position
		parts := strings.Split(line, "->")
		for _, p := range parts {
			xy := strings.Split(strings.TrimSpace(p), ",")
			if len(xy) != 2 {
				panic(fmt.Errorf("bad line xy: %s (%q)", line, p))
			}
			col, err1 := strconv.Atoi(xy[0])
			row, err2 := strconv.Atoi(xy[1])
			if err1 != nil || err2 != nil {
				panic(fmt.Errorf("XY: %v %v", err1, err2))
			}
			if row < minRow {
				minRow = row
			}
			if row > maxRow {
				maxRow = row
			}
			if col < minCol {
				minCol = col
			}
			if col > maxCol {
				maxCol = col
			}
			linePos = append(linePos, position{row, col})
			if row > maxRow {
				maxRow = row
			}
		}
		ret = append(ret, linePos)
	}
	return &puzzleInput{
		posLines: ret,
		minRow:   minRow,
		minCol:   minCol,
		maxRow:   maxRow,
		maxCol:   maxCol,
	}
}

func minMax(i, j int) (int, int) {
	if i > j {
		return j, i
	}
	return i, j
}

type puzzleInput struct {
	g        grid
	posLines [][]position
	minRow   int
	minCol   int
	maxRow   int
	maxCol   int
}

func fillGrid(pin *puzzleInput) {
	g := pin.g
	posLines := pin.posLines
	maybeFill := func(row, col int) {
		p := position{row, col}
		if g.value(p) == 0 {
			g.fill(p)
		}
	}
	for _, line := range posLines {
		prev := line[0]
		for i := 1; i < len(line); i++ {
			curr := line[i]
			switch {
			case curr.row == prev.row:
				a, b := minMax(curr.col, prev.col)
				for index := a; index <= b; index++ {
					maybeFill(curr.row, index)
				}
			case curr.col == prev.col:
				a, b := minMax(curr.row, prev.row)
				for index := a; index <= b; index++ {
					maybeFill(index, curr.col)
				}
			default:
				panic(fmt.Errorf("no straight line: %v, %v", curr, prev))
			}
			prev = curr
		}
	}
}

var chasm = position{-1, -1}

func fallDown(pin *puzzleInput, p position, isFloor bool) (landed position) {
	if p.row >= pin.maxRow {
		if !isFloor {
			return chasm
		} else {
			return p
		}
	}
	if !isFloor {
		if p.col <= pin.minCol {
			return chasm
		}
		if p.col >= pin.maxCol {
			return chasm
		}
	}
	g := pin.g

	test := position{p.row + 1, p.col}
	if g.value(test) == 0 {
		return fallDown(pin, test, isFloor)
	}

	testLeft := position{p.row + 1, p.col - 1}
	if g.value(testLeft) == 0 {
		return fallDown(pin, testLeft, isFloor)
	}
	testRight := position{p.row + 1, p.col + 1}
	if g.value(testRight) == 0 {
		return fallDown(pin, testRight, isFloor)
	}
	return p
}

func runP1(in string) int {
	pin := toPositions(in)
	pin.g = newGrid(pin.maxRow+1, pin.maxCol+1)
	fillGrid(pin)
	count := 0
	for {
		pos := fallDown(pin, position{0, 500}, false)
		if pos == chasm {
			break
		}
		count++
		pin.g.fill(pos)
	}
	return count
}

func runP2(in string) int {
	pin := toPositions(in)
	pin.maxRow += 1
	pin.g = newGrid(pin.maxRow+1, 2*pin.maxCol)
	fillGrid(pin)
	count := 0
	for {
		orig := position{0, 500}
		pos := fallDown(pin, orig, true)
		if pos == chasm {
			panic("chasm not expected in P2")
		}
		count++
		if pos == orig {
			break
		}
		pin.g.fill(pos)
	}
	return count
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
