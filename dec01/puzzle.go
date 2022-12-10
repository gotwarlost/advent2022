package dec01

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var actualInput string

type elf struct {
	number   int
	calories int
}

func findElvesByCaloriesDesc(elfBlocks []string) []elf {
	var elves []elf
	for i, b := range elfBlocks {
		vals := strings.Split(b, "\n")
		var cals int
		for _, v := range vals {
			n, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			cals += n
		}
		elves = append(elves, elf{number: i + 1, calories: cals})
	}
	sort.Slice(elves, func(i, j int) bool {
		return elves[i].calories > elves[j].calories
	})
	return elves
}

func topN(sortedElves []elf, n int) int {
	total := 0
	for i := 0; i < n; i++ {
		total += sortedElves[i].calories
	}
	return total
}

func run(input string) (top, top3 int) {
	input = strings.TrimSpace(input)
	elfBlocks := strings.Split(input, "\n\n")
	elves := findElvesByCaloriesDesc(elfBlocks)
	return topN(elves, 1), topN(elves, 3)
}

func Run() {
	t1, t3 := run(actualInput)
	fmt.Println("TOP 1:", t1)
	fmt.Println("TOP 3:", t3)
}
