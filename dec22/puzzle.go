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
	v := int(d)
	switch rot {
	case RotLeft:
		v--
		if v < 0 {
			v = int(dirUp)
		}
	case RotRight:
		v++
		if v > int(dirUp) {
			v = 0
		}
	default:
		panic("bad rot")
	}
	return Direction(v)
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
	regionSize    int
	routes        map[regionEdge]regionEdge
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

func (g *Grid) nextP1(pos Point, dir Direction) (Point, bool) {
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

func (g *Grid) applyP1(m Move) {
	if m.Type == MoveTurn {
		g.walkDirection = g.walkDirection.Rotate(m.Rotation)
		return
	}
	for i := 0; i < m.Steps; i++ {
		nextPos, found := g.nextP1(g.walkLocation, g.walkDirection)
		if !found {
			break
		}
		g.walkLocation = nextPos
	}
}

func (g *Grid) regionFor(point Point) (region, error) {
	rRows := g.rows / g.regionSize
	rCols := g.cols / g.regionSize
	for rRow := 0; rRow < rRows; rRow++ {
		for rCol := 0; rCol < rCols; rCol++ {
			rowStart := rRow * g.regionSize
			colStart := rCol * g.regionSize
			rowEnd := rowStart + g.regionSize
			colEnd := colStart + g.regionSize
			if point.row >= rowStart && point.row < rowEnd && point.col >= colStart && point.col < colEnd {
				return region{rRow, rCol}, nil
			}
		}
	}
	return region{}, fmt.Errorf("point %v not in any region", point)
}

func (g *Grid) toRegionCoords(gridPoint Point) Point {
	return Point{gridPoint.row % g.regionSize, gridPoint.col % g.regionSize}
}

func (g *Grid) toGridCoords(r region, regionPoint Point) Point {
	return Point{r.row*g.regionSize + regionPoint.row, r.col*g.regionSize + regionPoint.col}
}

func (g *Grid) moveRegion(current region, from Point, dir Direction) (region, Point, Direction) {
	var targetEdge regionEdge
	var outDir Direction

	localInPoint := g.toRegionCoords(from)
	switch dir {
	case dirLeft:
		targetEdge = g.routes[regionEdge{current, eLeft}]
	case dirRight:
		targetEdge = g.routes[regionEdge{current, eRight}]
	case dirUp:
		targetEdge = g.routes[regionEdge{current, eTop}]
	case dirDown:
		targetEdge = g.routes[regionEdge{current, eBottom}]
	}

	flip := func(i int) int {
		return g.regionSize - 1 - i
	}

	// implement rotation the hard way :(
	var row, col int
	switch targetEdge.edge {
	case eLeft:
		col = 0
		outDir = dirRight
		switch dir {
		case dirLeft:
			row = flip(localInPoint.row)
		case dirUp:
			row = localInPoint.col
		case dirRight:
			row = localInPoint.row
		case dirDown:
			row = flip(localInPoint.col)
		default:
			panic("1")
		}
	case eTop:
		row = 0
		outDir = dirDown
		switch dir {
		case dirLeft:
			col = localInPoint.row
		case dirUp:
			col = flip(localInPoint.col)
		case dirRight:
			col = flip(localInPoint.row)
		case dirDown:
			col = localInPoint.col
		default:
			panic("2")
		}
	case eRight:
		col = g.regionSize - 1
		outDir = dirLeft
		switch dir {
		case dirLeft:
			row = localInPoint.row
		case dirUp:
			row = flip(localInPoint.col)
		case dirRight:
			row = flip(localInPoint.row)
		case dirDown:
			row = localInPoint.col
		default:
			panic("3")
		}
	case eBottom:
		row = g.regionSize - 1
		outDir = dirUp
		switch dir {
		case dirLeft:
			col = flip(localInPoint.row)
		case dirUp:
			col = localInPoint.col
		case dirRight:
			col = localInPoint.row
		case dirDown:
			col = flip(localInPoint.col)
		default:
			panic("4")
		}
	}
	return targetEdge.region, g.toGridCoords(targetEdge.region, Point{row, col}), outDir
}

func (g *Grid) regionStartRow(r region) int { return g.regionSize * r.row }
func (g *Grid) regionEndRow(r region) int   { return g.regionSize*r.row + g.regionSize - 1 }
func (g *Grid) regionStartCol(r region) int { return g.regionSize * r.col }
func (g *Grid) regionEndCol(r region) int   { return g.regionSize*r.col + g.regionSize - 1 }

func (g *Grid) nextP2(walk Move) {
	currentRegion, err := g.regionFor(g.walkLocation)
	if err != nil {
		panic(err)
	}
	for i := 0; i < walk.Steps; i++ {
		pos := g.walkLocation
		dir := g.walkDirection
		p := pos.next(dir)

		r, err := g.regionFor(p)
		if err != nil || currentRegion != r {
			r, p, dir = g.moveRegion(currentRegion, pos, dir)
		}
		if g.Value(p) == cellVoid {
			panic("no void expected")
		}
		if g.Value(p) == cellWall {
			return // leave everything as-is
		}
		currentRegion = r
		g.walkLocation = p
		g.walkDirection = dir
	}
}

func (g *Grid) applyP2(m Move) {
	if m.Type == MoveTurn {
		g.walkDirection = g.walkDirection.Rotate(m.Rotation)
		return
	}
	g.nextP2(m)
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
	leading := grid.rows
	if grid.cols > leading {
		leading = grid.cols
	}
	if leading%4 != 0 {
		panic("4X3 matrix expected")
	}
	grid.regionSize = leading / 4
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
		g.applyP1(move)
	}

	return 1000*(g.walkLocation.row+1) + 4*(g.walkLocation.col+1) + int(g.walkDirection)
}

func runP2(in string, routes map[regionEdge]regionEdge) int {
	log.SetFlags(0)
	pin := toPuzzleInput(in)
	g := pin.grid
	g.routes = routes
	m := pin.moves
	firstCol := -1

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
	log.Println("WALKLOC:", g.walkLocation, g.toRegionCoords(g.walkLocation), g.walkDirection)
	log.Println(g.regionFor(g.walkLocation))
	for _, move := range m {
		g.applyP2(move)
	}

	return 1000*(g.walkLocation.row+1) + 4*(g.walkLocation.col+1) + int(g.walkDirection)
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input, mainRouteMap()))
}
