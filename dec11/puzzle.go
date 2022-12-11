package dec11

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type monkey struct {
	items         []int64
	operation     func(x int64) int64
	divisor       int64
	divTarget     int
	defaultTarget int
	inspections   int64
}

func (m *monkey) String() string {
	return fmt.Sprintf("items: %v, inspections: %d", m.items, m.inspections)
}

func (m *monkey) doRound(all []*monkey, worryDivisor int64) {
	for _, item := range m.items {
		m.inspections++
		worryBase := m.operation(item)
		worry := worryBase / worryDivisor

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

func toMonkeys(in string) []*monkey {
	var ret []*monkey
	lines := strings.Split(strings.TrimSpace(in), "\n")
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
			var operand int64
			if what != "old" {
				n, err := strconv.ParseInt(what, 10, 64)
				if err != nil {
					panic(err)
				}
				operand = n
			}
			currMonkey.operation = func(x int64) int64 {
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
			d, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				panic(err)
			}
			currMonkey.divisor = d
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
				n, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
				if err != nil {
					panic(err)
				}
				currMonkey.items = append(currMonkey.items, n)
			}
		default:
			panic(fmt.Errorf("no match for line: %q", line))
		}
	}
	return ret
}

func run(in string, divisor int64, rounds int) []*monkey {
	monkeys := toMonkeys(in)
	for round := 0; round < rounds; round++ {
		for _, m := range monkeys {
			m.doRound(monkeys, divisor)
		}
		print := func(i int) {
			log.Println("ROUNDS:", i)
			for _, m := range monkeys {
				log.Printf("%+v\n", m)
			}
		}
		switch round {
		case 0, 14, 19, 999, 1999, 2999:
			print(round)
		}
	}
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspections > monkeys[j].inspections
	})
	return monkeys
}

func runP1(in string) int64 {
	list := run(in, 3, 20)
	return list[0].inspections * list[1].inspections
}

func runP2(in string) int64 {
	list := run(in, 1, 10000)
	return list[0].inspections * list[1].inspections
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP1(input))
}
