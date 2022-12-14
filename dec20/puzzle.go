package dec20

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type cell struct {
	value int64
	next  *cell
	prev  *cell
}

func (c *cell) setNext(next *cell) {
	c.next = next
	next.prev = c
}

type ring struct {
	cellMap   map[int]*cell
	zeroValue *cell
	length    int
}

func toCells(in string, mult int) *ring {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	var numbers []int
	for _, l := range lines {
		n, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		numbers = append(numbers, n)
	}

	ret := &ring{
		cellMap: map[int]*cell{},
		length:  len(numbers),
	}
	var head, tail *cell
	for i, num := range numbers {
		c := &cell{value: int64(num) * int64(mult)}
		ret.cellMap[i] = c
		if num == 0 {
			if ret.zeroValue != nil {
				panic("multi-zero")
			}
			ret.zeroValue = c
		}
		if i == 0 {
			head = c
		} else {
			tail.setNext(c)
		}
		tail = c
	}
	tail.setNext(head)
	return ret
}

func run(in string, mult int, times int) int64 {
	cells := toCells(in, mult)

	for x := 0; x < times; x++ {
		for i := 0; i < cells.length; i++ {
			c := cells.cellMap[i]
			val := c.value
			left := false
			if val < 0 {
				left = true
				val = -val
			}
			moves := val % int64(cells.length-1)
			for pos := int64(0); pos < moves; pos++ {
				oldPrev, oldNext := c.prev, c.next
				switch left {
				case true:
					// prev2 -> prev -> c -> next
					// prev2 -> c -> prev -> next
					prev2 := c.prev.prev
					prev2.setNext(c)
					c.setNext(oldPrev)

				default:
					// prev -> c -> next -> next2
					// prev -> next -> c -> next2
					next2 := c.next.next
					oldNext.setNext(c)
					c.setNext(next2)
				}
				oldPrev.setNext(oldNext)
			}
		}
	}

	curr := cells.zeroValue
	total := int64(0)
	for i := 0; i <= 3000; i++ {
		switch i {
		case 1000, 2000, 3000:
			total += curr.value
		}
		curr = curr.next
	}
	return total
}

func RunP1() {
	fmt.Println(run(input, 1, 1))
}

func RunP2() {
	fmt.Println(run(input, 811589153, 10))
}
