package dec23

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type Direction int

const (
	_ Direction = iota
	North
	South
	West
	East
)

func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case South:
		return "South"
	case West:
		return "West"
	case East:
		return "East"
	}
	return "XXX"
}

type dirList []Direction

func (d dirList) next() dirList {
	var ret []Direction
	for i := 1; i < len(d); i++ {
		ret = append(ret, d[i])
	}
	ret = append(ret, d[0])
	return ret
}

type Point struct {
	row, col int
}

func (p Point) offset(row, col int) Point {
	return Point{p.row + row, p.col + col}
}

type Elf struct {
	id  int
	pos Point
}

type Collection struct {
	elves     []Elf
	positions map[Point]int
}

func (c *Collection) next(dl dirList) (*Collection, int) {
	dups := map[Point][]int{}
	proposals := map[int]Point{}

	anyFound := func(points ...Point) bool {
		for _, p := range points {
			if _, ok := c.positions[p]; ok {
				return true
			}
		}
		return false
	}

	for _, elf := range c.elves {
		pos := elf.pos
		// orig := elf.pos
		if anyFound(
			pos.offset(-1, -1), pos.offset(-1, 0), pos.offset(-1, 1),
			pos.offset(0, -1), pos.offset(0, 1),
			pos.offset(1, -1), pos.offset(1, 0), pos.offset(1, 1),
		) {
		loop:
			for _, d := range dl {
				switch d {
				case North:
					if anyFound(pos.offset(-1, -1), pos.offset(-1, 0), pos.offset(-1, 1)) {
						continue loop
					}
					pos = pos.offset(-1, 0)
					break loop
				case South:
					if anyFound(pos.offset(1, -1), pos.offset(1, 0), pos.offset(1, 1)) {
						continue loop
					}
					pos = pos.offset(1, 0)
					break loop
				case West:
					if anyFound(pos.offset(-1, -1), pos.offset(0, -1), pos.offset(1, -1)) {
						continue loop
					}
					pos = pos.offset(0, -1)
					break loop
				case East:
					if anyFound(pos.offset(-1, 1), pos.offset(0, 1), pos.offset(1, 1)) {
						continue loop
					}
					pos = pos.offset(0, 1)
					break loop
				default:
					panic("bad dir")
				}
			}
		}
		proposals[elf.id] = pos
		dups[pos] = append(dups[pos], elf.id)
	}

	numMoves := 0
	var elves []Elf
	positions := map[Point]int{}
	for _, elf := range c.elves {
		old := elf.pos
		proposed := proposals[elf.id]
		var newElf Elf
		if len(dups[proposed]) > 1 {
			newElf = Elf{id: elf.id, pos: old}
		} else {
			newElf = Elf{id: elf.id, pos: proposed}
		}
		elves = append(elves, newElf)
		positions[newElf.pos] = newElf.id
		if old != newElf.pos {
			numMoves++
		}
	}
	return &Collection{
		elves:     elves,
		positions: positions,
	}, numMoves
}

func (c *Collection) boundingBox() (Point, Point) {
	minRow := 1000000
	minCol := 1000000
	maxRow := -1000000
	maxCol := -1000000
	for _, e := range c.elves {
		row, col := e.pos.row, e.pos.col
		if row > maxRow {
			maxRow = row
		}
		if row < minRow {
			minRow = row
		}
		if col > maxCol {
			maxCol = col
		}
		if col < minCol {
			minCol = col
		}
	}
	return Point{minRow, minCol}, Point{maxRow, maxCol}
}

func (c *Collection) print(title string) {
	fmt.Println(title)
	fmt.Println("-----")
	min, max := c.boundingBox()
	for i := min.row - 2; i <= max.row+5; i++ {
		for j := min.col - 2; j <= max.col+5; j++ {
			if _, ok := c.positions[Point{i, j}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func toCollection(in string) *Collection {
	var elves []Elf
	positions := map[Point]int{}
	id := 0
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for row, line := range lines {
		for col := 0; col < len(line); col++ {
			if line[col] == '#' {
				id++
				pos := Point{row, col}
				elf := Elf{id: id, pos: pos}
				positions[pos] = id
				elves = append(elves, elf)
			}
		}
	}
	return &Collection{
		elves:     elves,
		positions: positions,
	}
}

func runP1(in string) int {
	coll := toCollection(in)
	dl := dirList([]Direction{North, South, West, East})
	coll.print("start")
	for i := 0; i < 10; i++ {
		coll, _ = coll.next(dl)
		dl = dl.next()
		coll.print(fmt.Sprintf("i:%d", i))
	}
	min, max := coll.boundingBox()
	count := (max.row-min.row+1)*(max.col-min.col+1) - len(coll.elves)
	return count
}

func runP2(in string) int {
	coll := toCollection(in)
	dl := dirList([]Direction{North, South, West, East})
	count := 0
	var moves int
	for {
		count++
		coll, moves = coll.next(dl)
		if moves == 0 {
			break
		}
		dl = dl.next()
	}
	return count
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
