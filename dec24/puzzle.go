package dec24

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed input.txt
var input string

type Direction byte

const (
	_ Direction = iota
	DirUp
	DirLeft
	DirDown
	DirRight
)

func (d Direction) String() string {
	switch d {
	case DirUp:
		return "^"
	case DirLeft:
		return "<"
	case DirDown:
		return "v"
	case DirRight:
		return ">"
	}
	return "X"
}

type Point struct {
	row, col int
}

type Blizzard struct {
	id    int
	point Point
	dir   Direction
}

func (b *Blizzard) state() []byte {
	return []byte(fmt.Sprintf("%d:%d:%v", b.point.row, b.point.col, b.dir))
}

type Grid struct {
	rows          int
	cols          int
	entrance      Point
	exit          Point
	blizzardBlock map[Point][]int
	wallBlock     map[Point]bool
	blizzards     map[int]Blizzard
}

func (g *Grid) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("rows: %d, cols: %d, entrance: %v, exit: %v\n", g.rows, g.cols, g.entrance, g.exit))
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			p := Point{i, j}
			switch {
			case p == g.entrance:
				b.WriteString("e")
			case p == g.exit:
				b.WriteString("x")
			case g.wallBlock[p]:
				b.WriteString("#")
			case len(g.blizzardBlock[p]) == 1:
				b.WriteString(fmt.Sprint(g.blizzards[g.blizzardBlock[p][0]].dir))
			case len(g.blizzardBlock[p]) > 0:
				b.WriteString(fmt.Sprint(len(g.blizzardBlock[p])))
			default:
				b.WriteString(".")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func toGrid(in string) *Grid {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	blizzardBlock := map[Point][]int{}
	wallBlock := map[Point]bool{}
	blizzards := map[int]Blizzard{}
	bid := 0
	for i, line := range lines {
		for j, ch := range []byte(line) {
			p := Point{i, j}
			switch ch {
			case '#':
				wallBlock[p] = true
			case '<':
				bid++
				blizzards[bid] = Blizzard{bid, p, DirLeft}
				blizzardBlock[p] = append(blizzardBlock[p], bid)
			case '^':
				bid++
				blizzards[bid] = Blizzard{bid, p, DirUp}
				blizzardBlock[p] = append(blizzardBlock[p], bid)
			case '>':
				bid++
				blizzards[bid] = Blizzard{bid, p, DirRight}
				blizzardBlock[p] = append(blizzardBlock[p], bid)
			case 'v':
				bid++
				blizzards[bid] = Blizzard{bid, p, DirDown}
				blizzardBlock[p] = append(blizzardBlock[p], bid)
			case '.':
			}
		}
	}
	g := &Grid{
		rows:          len(lines),
		cols:          len(lines[0]),
		blizzardBlock: blizzardBlock,
		wallBlock:     wallBlock,
		blizzards:     blizzards,
	}
	var e, x Point
	for i := 0; i < g.cols; i++ {
		if lines[0][i] == '.' {
			e = Point{0, i}
		}
		if lines[g.rows-1][i] == '.' {
			x = Point{g.rows - 1, i}
		}
	}
	g.entrance = e
	g.exit = x
	return g
}

func (g *Grid) isBlocked(p Point) bool {
	if p.row < 0 || p.col < 0 || p.row >= g.rows || p.col >= g.cols {
		return true
	}
	_, b1 := g.blizzardBlock[p]
	_, b2 := g.wallBlock[p]
	return b1 || b2
}

func (g *Grid) next() *Grid {
	blizzardBlock := map[Point][]int{}
	blizzards := map[int]Blizzard{}
	for bid, b := range g.blizzards {
		p := b.point
		switch b.dir {
		case DirLeft:
			p.col--
			if p.col < 1 {
				p.col = g.cols - 2
			}
		case DirRight:
			p.col++
			if p.col >= g.cols-1 {
				p.col = 1
			}
		case DirUp:
			p.row--
			if p.row < 1 {
				p.row = g.rows - 2
			}
		case DirDown:
			p.row++
			if p.row >= g.rows-1 {
				p.row = 1
			}
		}
		blizzards[bid] = Blizzard{id: bid, point: p, dir: b.dir}
		blizzardBlock[p] = append(blizzardBlock[p], bid)
	}
	return &Grid{
		rows:          g.rows,
		cols:          g.cols,
		entrance:      g.entrance,
		exit:          g.exit,
		blizzardBlock: blizzardBlock,
		wallBlock:     g.wallBlock,
		blizzards:     blizzards,
	}
}

func (g *Grid) state() string {
	var blizzardStates []string
	for _, b := range g.blizzards {
		blizzardStates = append(blizzardStates, string(b.state()))
	}
	sort.Strings(blizzardStates)
	return strings.Join(blizzardStates, "\x00")
}

type gridTime struct {
	time int
	g    *Grid
}

type GridMovie struct {
	grids map[int]*Grid
	last  int
}

var printed = false

func (gm *GridMovie) gridFor(time int) *Grid {
	if time <= gm.last {
		return gm.grids[time]
	}
	if time != gm.last+1 {
		panic(fmt.Errorf("need ascending requests: want %d, have %d", time, gm.last))
	}
	g := gm.grids[gm.last].next()
	gm.last = time
	gm.grids[time] = g
	if gm.last%1000 == 0 {
		fmt.Printf("G%d ...", gm.last)
	}
	return g
}

type Stats struct {
	bestTime  int
	seenMoves map[Move]bool
}

type Move struct {
	time  int
	point Point
}

var counter int

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m Move) next(gm *GridMovie, goal Point, stats *Stats) {
	if stats.seenMoves[m] {
		return
	}
	stats.seenMoves[m] = true
	counter++
	if counter%1000000 == 0 {
		fmt.Printf("M%d... ", counter/1000000)
	}
	g := gm.gridFor(m.time + 1)
	if m.point == goal {
		if stats.bestTime > m.time {
			stats.bestTime = m.time
		}
		return
	}
	movingDown := goal.row >= m.point.row && goal.col >= m.point.col
	incr := 1
	if !movingDown {
		incr = -1
	}
	colDiff := abs(goal.col - m.point.col)
	rowDiff := abs(goal.row - m.point.row)
	minSteps := colDiff + rowDiff
	if m.time+minSteps > stats.bestTime {
		return
	}
	t := m.time + 1
	var candidates []Move

	maybeAdd := func(p Point) {
		if !g.isBlocked(p) {
			candidates = append(candidates, Move{t, p})
		}
	}

	colPref := colDiff > rowDiff
	// prefer to go right or down
	nextRow, nextCol := Point{m.point.row + incr, m.point.col}, Point{m.point.row, m.point.col + incr}
	if colPref {
		maybeAdd(nextCol)
		maybeAdd(nextRow)
	} else {
		maybeAdd(nextRow)
		maybeAdd(nextCol)
	}
	// prefer to stay put
	maybeAdd(m.point)
	// backtrack if needed
	prevRow, prevCol := Point{m.point.row - incr, m.point.col}, Point{m.point.row, m.point.col - incr}
	if colPref {
		maybeAdd(prevRow)
		maybeAdd(prevCol)
	} else {
		maybeAdd(prevCol)
		maybeAdd(prevRow)
	}
	for _, c := range candidates {
		c.next(gm, goal, stats)
	}
}

func crossOnce(g *Grid, start, goal Point, budget int) (*GridMovie, int) {
	gm := &GridMovie{
		grids: map[int]*Grid{0: g},
		last:  0,
	}
	s := &Stats{bestTime: budget, seenMoves: map[Move]bool{}}
	m := Move{point: start}
	m.next(gm, goal, s)
	return gm, s.bestTime
}

func runP1(in string, budget int) int {
	log.SetFlags(0)
	g := toGrid(in)
	log.Println(g)
	_, best := crossOnce(g, g.entrance, g.exit, budget)
	return best
}

func runP2(in string, budget int) int {
	log.SetFlags(0)
	g := toGrid(in)
	gm, best1 := crossOnce(g, g.entrance, g.exit, budget)
	g = gm.gridFor(best1)
	log.Println("B1:", best1)
	gm, best2 := crossOnce(g, g.exit, g.entrance, budget)
	g = gm.gridFor(best2)
	log.Println("B2:", best2)
	gm, best3 := crossOnce(g, g.entrance, g.exit, budget)
	log.Println("B3:", best3)
	return best1 + best2 + best3
}

const mainBudget = 50000

func RunP1() {
	fmt.Println(runP1(input, mainBudget))
}

func RunP2() {
	fmt.Println(runP2(input, mainBudget))
}
