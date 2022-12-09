package dec9

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

func absDiff(a, b int) int {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d
}

func moveTail(head, tail point) point {
	if absDiff(head.x, tail.x) < 2 && absDiff(head.y, tail.y) < 2 {
		return tail
	}
	switch {
	case head.x == tail.x || head.y == tail.y:
		if head.x != tail.x {
			if head.x > tail.x {
				tail.x++
			} else {
				tail.x--
			}
		}
		if head.y != tail.y {
			if head.y > tail.y {
				tail.y++
			} else {
				tail.y--
			}
		}
	default: // diagonal
		if head.y > tail.y {
			tail.y++
		} else {
			tail.y--
		}
		if head.x < tail.x {
			tail.x--
		} else {
			tail.x++
		}
	}
	return tail
}

func run(in string, knots int) int {
	seen := map[point]bool{{0, 0}: true}
	rope := make([]point, knots)

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
		for i := 0; i < count; i++ {
			for k := 0; k < knots-1; k++ {
				head := rope[k]
				tail := rope[k+1]
				if k == 0 { // move head
					switch parts[0] {
					case "D":
						head.y--
					case "U":
						head.y++
					case "L":
						head.x--
					case "R":
						head.x++
					default:
						panic(parts[0])
					}
				}
				tail = moveTail(head, tail)
				rope[k] = head
				rope[k+1] = tail
				if k == knots-2 {
					seen[tail] = true
				}
			}
		}
	}
	return len(seen)
}

func RunP1() {
	fmt.Println(run(input, 2))
}

func RunP2() {
	fmt.Println(run(input, 10))
}
