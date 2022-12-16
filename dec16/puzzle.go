package dec16

import (
	_ "embed"
	"fmt"
	"log"
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

func visit(valves []*valve, order []int, r routeMap, initialCost int) (pressure int) {
	pos := order[0]
	v := valves[pos]
	prev := v

	timeLeft := 30 - initialCost
	pressure = v.flow * (timeLeft - 1) // valve flows one minute later

	for i, valveIndex := range order {
		if i == 0 {
			continue
		}
		v = valves[valveIndex]
		cost, ok := r.getCost(prev, v)
		if !ok {
			panic(fmt.Errorf("no cost from %s to %s", prev.name, v.name))
		}
		timeLeft -= cost + 1
		if timeLeft <= 0 {
			break
		}
		pressure += v.flow * (timeLeft - 1)
		prev = v
	}
	return pressure
}

func runP1(in string) int {
	valves := toValves(in)

	var aaValve *valve
	flowingValves := map[string]*valve{}
	var flowingValveIndices []int
	for i, v := range valves {
		if v.name == "AA" {
			aaValve = v
		}
		if v.flow > 0 {
			flowingValves[v.name] = v
			flowingValveIndices = append(flowingValveIndices, i)
		}
	}

	routes := routeMap{}

	addRoute := func(current *valve, other *valve, score int) {
		if current.name == other.name {
			return
		}
		if flowingValves[other.name] == nil {
			return
		}
		routes.setCost(current, other, score)
	}

	aaScores := map[string]int{}
	computePathScores(aaValve, aaScores, 0)

	for _, v1 := range valves {
		if flowingValves[v1.name] == nil {
			continue
		}
		scores := map[string]int{}
		computePathScores(v1, scores, 0)
		for _, v2 := range valves {
			score, ok := scores[v2.name]
			if ok {
				addRoute(v1, v2, score)
			}
		}
	}

	for r, c := range routes {
		log.Println(r, "=>", c)
	}

	ch := itertools.PermutationsInt(flowingValveIndices, len(flowingValveIndices))
	max := 0
	for order := range ch {
		first := valves[order[0]].name
		aaCost, ok := aaScores[first]
		if !ok {
			panic("no route from AA to " + first)
		}
		pressure := visit(valves, order, routes, aaCost)
		if pressure > max {
			max = pressure
		}
	}

	return max
}

func runP2(in string) int {
	return 3
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
