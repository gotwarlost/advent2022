package dec16

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ernestosuarez/itertools"
)

//go:embed input.txt
var input string

var reSpec = regexp.MustCompile(`^Valve (\S+) has flow rate=(\d+); tunnel(s?) lead(s?) to valve(s?) (.+)$`)

type valve struct {
	name  string
	flow  int
	links []*valve
}

func toValves(in string) []*valve {
	outlets := map[string][]string{}
	var valves []*valve
	valvesByName := map[string]*valve{}
	for _, line := range strings.Split(strings.TrimSpace(in), "\n") {
		m := reSpec.FindStringSubmatch(line)
		if len(m) == 0 {
			panic("bad line: " + line)
		}
		name := m[1]
		flow, err := strconv.Atoi(m[2])
		if err != nil {
			panic(err)
		}
		rest := m[6]
		parts := strings.Split(rest, ",")
		for _, p := range parts {
			outlets[name] = append(outlets[name], strings.TrimSpace(p))
		}
		v := &valve{
			name: name,
			flow: flow,
		}
		valves = append(valves, v)
		valvesByName[name] = v
	}
	for name, out := range outlets {
		v := valvesByName[name]
		for _, dest := range out {
			v2 := valvesByName[dest]
			if v2 == nil {
				panic("no valve:" + dest)
			}
			v.links = append(v.links, v2)
		}
	}
	return valves
}

func computePathScores(current *valve, scores map[string]int, score int) {
	// check if there is a path with an equal or lower score and abandon
	if prevScore, ok := scores[current.name]; ok {
		if prevScore <= score {
			return
		}
	}
	// set the current score
	scores[current.name] = score

	// compute for reachable steps
	for _, x := range current.links {
		computePathScores(x, scores, score+1)
	}
}

type route struct {
	source string
	target string
}

type routeMap map[route]int

func (r routeMap) setCost(v1, v2 *valve, cost int) {
	n1, n2 := v1.name, v2.name
	if n1 > n2 {
		n1, n2 = n2, n1
	}
	r[route{n1, n2}] = cost
}

func (r routeMap) getCost(v1, v2 *valve) (int, bool) {
	n1, n2 := v1.name, v2.name
	if n1 > n2 {
		n1, n2 = n2, n1
	}
	x, ok := r[route{n1, n2}]
	return x, ok
}

func (r routeMap) mustGetCost(v1, v2 *valve) int {
	c, ok := r.getCost(v1, v2)
	if !ok {
		panic(fmt.Errorf("no cost from %s to %s", v1.name, v2.name))
	}
	return c
}

func visit(start *valve, valves []*valve, order []int, r routeMap, totalTime int) (pressure int) {
	timeLeft := totalTime
	for _, valveIndex := range order {
		v := valves[valveIndex]
		cost, ok := r.getCost(start, v)
		if !ok {
			panic(fmt.Errorf("no cost from %s to %s", start.name, v.name))
		}
		timeLeft -= cost + 1
		if timeLeft <= 0 {
			break
		}
		extraPressure := v.flow * timeLeft
		pressure += extraPressure
		start = v
	}
	return pressure
}

func without(x map[string]*valve, name string) map[string]*valve {
	ret := make(map[string]*valve, len(x)-1)
	for k, v := range x {
		if k == name {
			continue
		}
		ret[k] = v
	}
	return ret
}

func computeMaxInternal(start *valve, remaining map[string]*valve, routes routeMap, pressure, timeLeft int) (outputPressure int) {
	max := pressure
	for _, v := range remaining {
		p := pressure
		tl := timeLeft
		cost := routes.mustGetCost(start, v) + 1
		tl -= cost
		if tl <= 0 {
			continue
		}
		p += tl * v.flow
		outPressure := computeMaxInternal(v, without(remaining, v.name), routes, p, tl)
		if outPressure > max {
			max = outPressure
		}
	}
	return max
}

func computeMax(start *valve, remaining map[string]*valve, routes routeMap, totalTime int) int {
	return computeMaxInternal(start, remaining, routes, 0, totalTime)
}

type puzzleContext struct {
	valves  []*valve
	first   *valve
	nonZero map[string]*valve
	routes  routeMap
}

func toPuzzleContext(in string) *puzzleContext {
	valves := toValves(in)

	var firstValve *valve
	// interesting means non-zero flows plus start valve
	interestingValves := map[string]*valve{}
	for _, v := range valves {
		if v.name == "AA" || v.flow > 0 {
			interestingValves[v.name] = v
		}
	}
	firstValve = interestingValves["AA"]

	routes := routeMap{}
	for _, v1 := range valves {
		if interestingValves[v1.name] == nil {
			continue
		}
		scores := map[string]int{}
		computePathScores(v1, scores, 0)
		for _, v2 := range valves {
			if v1.name == v2.name {
				continue
			}
			if interestingValves[v2.name] == nil {
				continue
			}
			score, ok := scores[v2.name]
			if ok {
				routes.setCost(v1, v2, score)
			}
		}
	}
	return &puzzleContext{
		valves:  valves,
		first:   firstValve,
		nonZero: without(interestingValves, "AA"),
		routes:  routes,
	}
}

func runP1(in string) int {
	c := toPuzzleContext(in)
	max := computeMax(c.first, c.nonZero, c.routes, 30)
	return max
}

func split(valves map[string]*valve, indices []string) (left, right map[string]*valve) {
	left = map[string]*valve{}
	right = map[string]*valve{}
	seen := map[string]bool{}
	for _, k := range indices {
		left[k] = valves[k]
		seen[k] = true
	}
	for k, v := range valves {
		if seen[k] {
			continue
		}
		right[k] = v
	}
	return left, right
}

func runP2(in string) int {
	c := toPuzzleContext(in)
	maxLen := len(c.nonZero)/2 + 1 // Bernard's nifty trick of taking advantage of mirror images
	var names []string
	for k := range c.nonZero {
		names = append(names, k)
	}
	max := 0
	for i := 1; i < maxLen; i++ {
		combinations := itertools.CombinationsStr(names, i)
		for slice := range combinations {
			left, right := split(c.nonZero, slice)
			m1 := computeMax(c.first, left, c.routes, 26)
			m2 := computeMax(c.first, right, c.routes, 26)
			if m1+m2 > max {
				max = m1 + m2
			}
		}
	}
	return max
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
