package dec22

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type Cell int

const (
	cellVoid Cell = iota
	cellWall
	cellOpen
)

func (c Cell) String() string {
	switch c {
	case cellOpen:
		return "."
	case cellWall:
		return "#"
	default:
		return "x"
	}
}

type Rotation int

const (
	_ Rotation = iota
	RotLeft
	RotRight
)

func (r Rotation) String() string {
	switch r {
	case RotLeft:
		return "L"
	case RotRight:
		return "R"
	default:
		return ""
	}
}

type Direction int

const (
	dirRight Direction = iota
	dirDown
	dirLeft
	dirUp
)

func (d Direction) String() string {
	switch d {
	case dirRight:
		return ">"
	case dirLeft:
		return "<"
	case dirUp:
		return "^"
	case dirDown:
		return "v"
	default:
		return ""
	}
}

func (d Direction) Rotate(rot Rotation) Direction {
	switch rot {
	case RotLeft:
		switch d {
		case dirRight:
			return dirUp
		case dirUp:
			return dirLeft
		case dirLeft:
			return dirDown
		case dirDown:
			return dirRight
		default:
			panic("bad dir")
		}
	case RotRight:
		switch d {
		case dirRight:
			return dirDown
		case dirUp:
			return dirRight
		case dirLeft:
			return dirUp
		case dirDown:
			return dirLeft
		default:
			panic("bad dir")
		}
	default:
		panic("bad rot")
	}
}

type Point struct {
	row, col int
}

func (p Point) next(dir Direction) Point {
	switch dir {
	case dirRight:
		p.col++
	case dirLeft:
		p.col--
	case dirDown:
		p.row++
	case dirUp:
		p.row--
	default:
		panic("bad dir for next")
	}
	return p
}

type MoveType int

const (
	_ MoveType = iota
	MoveWalk
	MoveTurn
)

func (m MoveType) String() string {
	switch m {
	case MoveWalk:
		return "Walk"
	case MoveTurn:
		return "Turn"
	default:
		return ""
	}
}

type Move struct {
	Type     MoveType
	Steps    int
	Rotation Rotation
}

type Grid struct {
	cols          int
	rows          int
	positions     [][]Cell
	walkLocation  Point
	walkDirection Direction
}

func (g *Grid) String() string {
	var b bytes.Buffer
	for _, row := range g.positions {
		b.WriteString("|")
		for _, col := range row {
			b.WriteString(col.String())
		}
		b.WriteString("|\n")
	}
	return b.String()
}

func (g *Grid) Value(pos Point) Cell {
	return g.positions[pos.row][pos.col]
}

func (g *Grid) firstNonEmptyCol(row int) int {
	line := g.positions[row]
	for i := 0; i < g.cols; i++ {
		if line[i] != cellVoid {
			return i
		}
	}
	panic("all empty")
}

func (g *Grid) lastNonEmptyCol(row int) int {
	line := g.positions[row]
	for i := g.cols - 1; i >= 0; i-- {
		if line[i] != cellVoid {
			return i
		}
	}
	panic("all empty")
}

func (g *Grid) firstNonEmptyRow(col int) int {
	for i := 0; i < g.rows; i++ {
		if g.positions[i][col] != cellVoid {
			return i
		}
	}
	panic("all empty")
}

func (g *Grid) lastNonEmptyRow(col int) int {
	for i := g.rows - 1; i >= 0; i-- {
		if g.positions[i][col] != cellVoid {
			return i
		}
	}
	panic("all empty")
}

func (g *Grid) next(pos Point, dir Direction) (Point, bool) {
	p := pos.next(dir)
	switch dir {
	case dirRight:
		if p.col >= g.cols || g.Value(p) == cellVoid {
			p.col = g.firstNonEmptyCol(pos.row)
		}
	case dirLeft:
		if p.col < 0 || g.Value(p) == cellVoid {
			p.col = g.lastNonEmptyCol(pos.row)
		}
	case dirUp:
		if p.row < 0 || g.Value(p) == cellVoid {
			p.row = g.lastNonEmptyRow(pos.col)
		}
	case dirDown:
		if p.row >= g.rows || g.Value(p) == cellVoid {
			p.row = g.firstNonEmptyRow(pos.col)
		}
	}
	if g.Value(p) == cellWall {
		return pos, false
	}
	return p, true
}

func (g *Grid) apply(m Move) {
	if m.Type == MoveTurn {
		g.walkDirection = g.walkDirection.Rotate(m.Rotation)
		return
	}
	for i := 0; i < m.Steps; i++ {
		nextPos, found := g.next(g.walkLocation, g.walkDirection)
		if !found {
			break
		}
		g.walkLocation = nextPos
	}
}

type PuzzleInput struct {
	grid  *Grid
	moves []Move
}

func toPuzzleInput(in string) *PuzzleInput {
	grid := &Grid{}
	in = strings.TrimRight(in, "\n")
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if line == "" {
			break
		}
		l := len(line)
		if l > grid.cols {
			extendBy := l - grid.cols
			var extension []Cell
			for e := 0; e < extendBy; e++ {
				extension = append(extension, cellVoid)
			}
			for i := range grid.positions {
				grid.positions[i] = append(grid.positions[i], extension...)
			}
			grid.cols = l
		}
		var row []Cell
		for _, b := range []byte(line) {
			switch b {
			case ' ':
				row = append(row, cellVoid)
			case '.':
				row = append(row, cellOpen)
			case '#':
				row = append(row, cellWall)
			}
		}
		if len(row) < grid.cols {
			extendBy := grid.cols - len(row)
			for i := 0; i < extendBy; i++ {
				row = append(row, cellVoid)
			}
		}
		grid.positions = append(grid.positions, row)
	}
	grid.rows = len(grid.positions)

	instructions := lines[len(lines)-1]

	var moves []Move
	i := 0
	for {
		b := instructions[i]
		switch b {
		case 'L':
			moves = append(moves, Move{Type: MoveTurn, Rotation: RotLeft})
			i++
		case 'R':
			moves = append(moves, Move{Type: MoveTurn, Rotation: RotRight})
			i++
		default:
			var stepsStr []byte
			for pos := i; pos < len(instructions); pos++ {
				c := instructions[pos]
				if c < '0' || c > '9' {
					break
				}
				stepsStr = append(stepsStr, c)
			}
			if len(stepsStr) == 0 {
				panic("what happened here?")
			}
			n, err := strconv.Atoi(string(stepsStr))
			if err != nil {
				panic(err)
			}
			moves = append(moves, Move{Type: MoveWalk, Steps: n})
			i += len(stepsStr)
		}
		if i >= len(instructions) {
			break
		}
	}
	return &PuzzleInput{
		grid:  grid,
		moves: moves,
	}
}

func runP1(in string) int {
	log.SetFlags(0)
	pin := toPuzzleInput(in)
	g := pin.grid
	m := pin.moves
	firstCol := -1

	log.Println(g.rows, g.cols)
	for i, c := range g.positions[0] {
		if c == cellOpen {
			firstCol = i
			break
		}
	}
	if firstCol == -1 {
		panic("no start position")
	}

	g.walkLocation = Point{row: 0, col: firstCol}
	g.walkDirection = dirRight

	for _, move := range m {
		g.apply(move)
	}

	return 1000*(g.walkLocation.row+1) + 4*(g.walkLocation.col+1) + int(g.walkDirection)
}

func runP2(in string) int64 {
	return 0
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
