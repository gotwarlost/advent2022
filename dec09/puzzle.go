package dec09

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type point struct {
	x, y int
}

type move byte

const (
	left  move = 'L'
	right move = 'R'
	down  move = 'D'
	up    move = 'U'
)

func toMoves(in string) []move {
	var moves []move
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(line)
		}
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		switch parts[0] {
		case "D", "U", "L", "R":
			// ok
		default:
			panic(fmt.Errorf("line: %s, bad move", line))
		}
		for i := 0; i < count; i++ {
			moves = append(moves, move(parts[0][0]))
		}
	}
	return moves
}

func absDiff(a, b int) int {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d
}

func moveTail(head, tail point) point {
	switch {
	case absDiff(head.x, tail.x) < 2 && absDiff(head.y, tail.y) < 2:
		return tail
	case head.x == tail.x && head.y != tail.y:
		if head.y > tail.y {
			tail.y++
		} else {
			tail.y--
		}
	case head.x != tail.x && head.y == tail.y:
		if head.x > tail.x {
			tail.x++
		} else {
			tail.x--
		}
	default: // diagonal
		if head.y > tail.y {
			tail.y++
		} else {
			tail.y--
		}
		if head.x > tail.x {
			tail.x++
		} else {
			tail.x--
		}
	}
	return tail
}

func countVisited(in string, knots int) int {
	moves := toMoves(in)
	seen := map[point]bool{{0, 0}: true}
	rope := make([]point, knots)

	for _, m := range moves {
		// move the head
		head := rope[0]
		switch m {
		case down:
			head.y--
		case up:
			head.y++
		case left:
			head.x--
		case right:
			head.x++
		}
		rope[0] = head
		for k := 1; k < knots; k++ {
			prev := rope[k-1]
			prevTail := rope[k]
			tail := moveTail(prev, prevTail)
			if tail == prevTail { // this hasn't moved so nothing else below it should move
				break
			}
			rope[k] = tail
		}
		seen[rope[knots-1]] = true
	}
	return len(seen)
}

func RunP1() {
	fmt.Println(countVisited(input, 2))
}

func RunP2() {
	fmt.Println(countVisited(input, 10))
}
