package dec11

import (
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type monkey struct {
	items         []int
	operation     func(x int) int
	divisor       int
	divTarget     int
	defaultTarget int
	inspections   int
}

func (m *monkey) String() string {
	return fmt.Sprintf("items: %v, inspections: %d", m.items, m.inspections)
}

func (m *monkey) doRound(all []*monkey, worryDivisor, commonDivisor int) {
	for _, item := range m.items {
		m.inspections++
		worryBase := m.operation(item)
		worry := worryBase / worryDivisor

		worry = worry % commonDivisor

		testOutcome := worry % m.divisor
		var target int
		if testOutcome == 0 {
			target = m.divTarget
		} else {
			target = m.defaultTarget
		}
		all[target].items = append(all[target].items, worry)
	}
	m.items = nil
}

var (
	monkeyRE     = regexp.MustCompile(`^Monkey\s+(\d+):$`)
	startItemsRE = regexp.MustCompile(`\s+Starting items: (.+)$`)
	opRE         = regexp.MustCompile(`^\s+Operation: new = old ([+*])\s+(.+)$`)
	testRE       = regexp.MustCompile(`^\s+Test: divisible by (\d+)$`)
	ifRE         = regexp.MustCompile(`^\s+If (true|false): throw to monkey (\d+)`)
)

func toMonkeys(in string) ([]*monkey, int) {
	var ret []*monkey
	lines := strings.Split(strings.TrimSpace(in), "\n")
	commonDivisor := 1
	var currMonkey *monkey
	for _, line := range lines {
		switch {
		case line == "":

		case monkeyRE.MatchString(line):
			currMonkey = &monkey{}
			ret = append(ret, currMonkey)
		case opRE.MatchString(line):
			m := opRE.FindStringSubmatch(line)
			operator := m[1]
			what := m[2]
			var operand int
			if what != "old" {
				n, err := strconv.Atoi(what)
				if err != nil {
					panic(err)
				}
				operand = n
			}
			currMonkey.operation = func(x int) int {
				target := operand
				if what == "old" {
					target = x
				}
				if operator == "+" {
					x += target
				} else {
					x *= target
				}
				return x
			}
		case testRE.MatchString(line):
			m := testRE.FindStringSubmatch(line)
			d, err := strconv.Atoi(m[1])
			if err != nil {
				panic(err)
			}
			currMonkey.divisor = d
			commonDivisor = commonDivisor * d
		case ifRE.MatchString(line):
			m := ifRE.FindStringSubmatch(line)
			target, err := strconv.Atoi(m[2])
			if err != nil {
				panic(err)
			}
			if m[1] == "true" {
				currMonkey.divTarget = target
			} else {
				currMonkey.defaultTarget = target
			}
		case startItemsRE.MatchString(line):
			m := startItemsRE.FindStringSubmatch(line)
			items := strings.Split(strings.TrimSpace(m[1]), ",")
			for _, item := range items {
				n, err := strconv.Atoi(strings.TrimSpace(item))
				if err != nil {
					panic(err)
				}
				currMonkey.items = append(currMonkey.items, n)
			}
		default:
			panic(fmt.Errorf("no match for line: %q", line))
		}
	}
	return ret, commonDivisor
}

func run(in string, divisor int, rounds int) []*monkey {
	monkeys, commonDivisor := toMonkeys(in)
	for round := 0; round < rounds; round++ {
		for _, m := range monkeys {
			m.doRound(monkeys, divisor, commonDivisor)
		}
	}
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspections > monkeys[j].inspections
	})
	return monkeys
}

func runP1(in string) int {
	list := run(in, 3, 20)
	return list[0].inspections * list[1].inspections
}

func runP2(in string) int {
	list := run(in, 1, 10000)
	return list[0].inspections * list[1].inspections
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
