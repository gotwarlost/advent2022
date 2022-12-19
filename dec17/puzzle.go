package dec17

import (
	"crypto/sha1"
	_ "embed"
	"encoding/hex"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type direction int

const (
	_ direction = iota
	dirLeft
	dirRight
	dirDown
)

type rock struct {
	_w    int
	rows  [][]int
	next  *rock
	index int
}

func (r *rock) width() int {
	return r._w
}

func (r *rock) height() int {
	return len(r.rows)
}

func (r *rock) leftAt(row int) int {
	if row < 0 || row >= r.height() {
		panic("bad left offset")
	}
	for i := 0; i < r.width(); i++ {
		if r.rows[row][i] == 1 {
			return i
		}
	}
	panic("internal leftAt")
}

func (r *rock) rightAt(row int) int {
	if row < 0 || row >= r.height() {
		panic("bad right offset")
	}
	for i := r.width() - 1; i >= 0; i-- {
		if r.rows[row][i] == 1 {
			return i
		}
	}
	panic("internal rightAt")
}

func (r *rock) topAt(col int) int {
	if col < 0 || col >= r.width() {
		panic("bad top col")
	}
	for i := 0; i < r.height(); i++ {
		if r.rows[i][col] == 1 {
			return i
		}
	}
	panic("internal topAt")
}

func (r *rock) makeLines(left, towerWidth int) [][]int {
	if left < 0 || left >= towerWidth {
		panic("bad left for makeLines")
	}
	var ret [][]int
	for y := 0; y < r.height(); y++ {
		var line []int
		for i := 0; i < left; i++ {
			line = append(line, 0)
		}
		for _, x := range r.rows[y] {
			line = append(line, x)
		}
		for i := left + r.width(); i < towerWidth; i++ {
			line = append(line, 0)
		}
		ret = append(ret, line)
	}
	return ret
}

func (r *rock) prettyPrint(left, towerWidth int) {
	fmt.Println("===")
	lines := r.makeLines(left, towerWidth)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		for _, x := range l {
			if x == 1 {
				fmt.Print("o")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

var lastRock *rock

func init() {
	// note: rocks are rendered bottom to top because we are moving "up"
	rocks := []*rock{
		{
			_w: 4,
			rows: [][]int{
				{1, 1, 1, 1},
			},
		},
		{
			_w: 3,
			rows: [][]int{
				{0, 1, 0},
				{1, 1, 1},
				{0, 1, 0},
			},
		},
		{
			_w: 3,
			rows: [][]int{
				{1, 1, 1},
				{0, 0, 1},
				{0, 0, 1},
			},
		},
		{
			_w: 1,
			rows: [][]int{
				{1},
				{1},
				{1},
				{1},
			},
		},
		{
			_w: 2,
			rows: [][]int{
				{1, 1},
				{1, 1},
			},
		},
	}

	l := len(rocks)
	for i := 0; i < l; i++ {
		rocks[i].index = i
		if i == l-1 {
			rocks[i].next = rocks[0]
		} else {
			rocks[i].next = rocks[i+1]
		}
	}
	lastRock = rocks[len(rocks)-1]
}

type push struct {
	dir   direction
	next  *push
	index int
}

func toPushes(in string) *push {
	vals := strings.Split(strings.TrimSpace(in), "")
	var seq []*push
	for i, v := range vals {
		switch v {
		case "<":
			seq = append(seq, &push{dir: dirLeft, index: i})
		case ">":
			seq = append(seq, &push{dir: dirRight, index: i})
		default:
			panic("bad push:" + v)
		}
	}
	l := len(seq)
	for i := 0; i < l; i++ {
		if i == l-1 {
			seq[i].next = seq[0]
		} else {
			seq[i].next = seq[i+1]
		}
	}
	return seq[len(seq)-1]
}

func prettyPrint(t [][]int, lastn int) {
	if lastn > 0 {
		t = t[len(t)-lastn:]
	}
	fmt.Println("---")
	for row := len(t) - 1; row >= 0; row-- {
		col := t[row]
		for _, x := range col {
			if x == 1 {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

const noValue = -1

type foundState struct {
	rockNum int64
	height  int64
}

func run(in string, numRocks int64) int64 {
	theRock := lastRock
	p := toPushes(in)
	towerWidth := 7
	var tower [][]int

	computeState := func() string {
		lookBehind := 10
		if len(tower) < lookBehind {
			return ""
		}
		h := sha1.New()
		write := func(x interface{}) {
			_, _ = h.Write([]byte(fmt.Sprint(x)))
		}
		write(theRock.index)
		write(p.index)
		for _, line := range tower[len(tower)-lookBehind:] {
			for _, x := range line {
				write(x)
			}
		}
		return hex.EncodeToString(h.Sum(nil))
	}

	valueAtIndex := func(row, col int) int {
		if row < 0 {
			return noValue
		}
		if col < 0 {
			return noValue
		}
		if row >= len(tower) {
			return noValue
		}
		if col >= towerWidth {
			return noValue
		}
		return tower[row][col]
	}

	seen := map[string]foundState{}

	addHeight := int64(0)
	rockNum := int64(0)
	checking := true
	for rockNum = 0; rockNum < numRocks; rockNum++ {
		theRock = theRock.next
		top := len(tower) + 3
		left := 2

		state := computeState()
		if checking && state != "" {
			if prev, ok := seen[state]; ok {
				rockCycle := rockNum - prev.rockNum
				heightForInterval := int64(len(tower)) - prev.height
				remainingRocks := numRocks - rockNum
				div := remainingRocks / rockCycle
				rem := remainingRocks % rockCycle
				addHeight = div * heightForInterval
				rockNum = numRocks - rem
				checking = false
			} else {
				seen[state] = foundState{
					rockNum: rockNum,
					height:  int64(len(tower)),
				}
			}
		}
	rockFall:
		for {
			p = p.next

			switch p.dir {
			case dirLeft:
				canLeft := left > 0
				if canLeft {
					for r := 0; r < theRock.height(); r++ {
						rLeft := theRock.leftAt(r)
						if valueAtIndex(top+r, left+rLeft-1) == 1 {
							canLeft = false
						}
					}
				}
				if canLeft {
					left--
				}
			case dirRight:
				right := left + theRock.width() - 1
				canRight := right < towerWidth-1
				if canRight {
					for r := 0; r < theRock.height(); r++ {
						rRight := theRock.rightAt(r)
						if valueAtIndex(top+r, left+rRight+1) == 1 {
							canRight = false
						}
					}
				}
				if canRight {
					left++
				}
			default:
				panic("bad dir")
			}

			canUp := top > 0
			if canUp {
				for t := 0; t < theRock.width(); t++ {
					rTop := theRock.topAt(t)
					y, x := top+rTop-1, left+t
					if valueAtIndex(y, x) == 1 {
						canUp = false
						break
					}
				}
			}

			if canUp {
				top--
				continue rockFall
			}

			union := func(line1, line2 []int) []int {
				if len(line1) != len(line2) {
					panic("wtf")
				}
				var ret []int
				for i := range line1 {
					lhs := line1[i]
					rhs := line2[i]
					if lhs+rhs > 1 {
						panic(fmt.Errorf("union mismatch: %v, %v", line1, line2))
					}
					if lhs == 1 || rhs == 1 {
						ret = append(ret, 1)
					} else {
						ret = append(ret, 0)
					}
				}
				return ret
			}

			lines := theRock.makeLines(left, towerWidth)
			for l := 0; l < len(lines); l++ {
				theTop := top + l
				if theTop > len(tower) {
					panic("unexpected top")
				}
				if theTop < len(tower) {
					tower[theTop] = union(tower[theTop], lines[l])
				} else {
					tower = append(tower, lines[l])
				}
			}
			break rockFall
		}
	}
	return int64(len(tower)) + addHeight
}

func RunP1() {
	fmt.Println(run(input, 2022))
}

const gazillion int64 = 1000000000000

func RunP2() {
	fmt.Println(run(input, gazillion))
}
