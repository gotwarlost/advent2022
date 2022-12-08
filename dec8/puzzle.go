package dec8

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type grid [][]int

func (g grid) rows() int {
	return len(g)
}

func (g grid) cols() int {
	return len(g[0])
}

func clone(in []int) []int {
	l := len(in)
	x := make([]int, l)
	copy(x, in)
	return x
}

func reverse(s []int) []int {
	l := len(s)
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s // not a new slice but make it easy to chain
}

func (g grid) horizontalSlice(row, start, end int) []int {
	ret := clone(g[row][start:end])
	return ret
}

func (g grid) verticalSlice(col, start, end int) []int {
	var ret []int
	for i := start; i < end; i++ {
		ret = append(ret, g[i][col])
	}
	return ret
}

func makeGrid(s string) grid {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var heights [][]int
	for _, line := range lines {
		els := strings.Split(line, "")
		var nums []int
		for _, el := range els {
			n, err := strconv.Atoi(el)
			if err != nil {
				panic(err)
			}
			nums = append(nums, n)
		}
		heights = append(heights, nums)
	}

	g := grid(heights)
	return g
}

func run(s string) int {
	g := makeGrid(s)
	checkVisibility := func(i, j int) (out bool) {
		h := g[i][j]
		isVisible := func(s []int) bool {
			visible := true
			for _, val := range s {
				if val >= h {
					visible = false
					break
				}
			}
			return visible
		}

		left, right, top, bottom := g.horizontalSlice(i, 0, j),
			reverse(g.horizontalSlice(i, j+1, g.cols())),
			g.verticalSlice(j, 0, i),
			reverse(g.verticalSlice(j, i+1, g.rows()))

		return isVisible(left) || isVisible(right) || isVisible(top) || isVisible(bottom)
	}
	visible := 0
	for i := 0; i < g.rows(); i++ {
		for j := 0; j < g.cols(); j++ {
			if checkVisibility(i, j) {
				visible++
			}
		}
	}
	return visible
}

func run2(s string) int {
	g := makeGrid(s)
	getScore := func(i, j int) int {
		h := g[i][j]
		scoreFor := func(s []int) int {
			count := 0
			for _, val := range s {
				count++
				if val >= h {
					break
				}
			}
			return count
		}

		left, right, top, bottom := reverse(g.horizontalSlice(i, 0, j)),
			g.horizontalSlice(i, j+1, g.cols()),
			reverse(g.verticalSlice(j, 0, i)),
			g.verticalSlice(j, i+1, g.rows())

		return scoreFor(left) * scoreFor(right) * scoreFor(top) * scoreFor(bottom)
	}
	max := 0
	for i := 0; i < g.rows(); i++ {
		for j := 0; j < g.cols(); j++ {
			score := getScore(i, j)
			if score > max {
				max = score
			}
		}
	}
	return max
}

func RunP1() {
	fmt.Println(run(input))
}

func RunP2() {
	fmt.Println(run2(input))
}
