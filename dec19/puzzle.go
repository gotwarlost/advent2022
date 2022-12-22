package dec19

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type Resource string

const (
	Nothing  Resource = ""
	Geode    Resource = "geode"
	Obsidian Resource = "obsidian"
	Clay     Resource = "clay"
	Ore      Resource = "ore"
)

type Counts struct {
	Geodes   int `json:"geodes,omitempty"`
	Obsidian int `json:"obsidian,omitempty"`
	Clay     int `json:"clay,omitempty"`
	Ores     int `json:"ores,omitempty"`
}

type Skips struct {
	Geodes   bool `json:"geodes,omitempty"`
	Obsidian bool `json:"obsidian,omitempty"`
	Clay     bool `json:"clay,omitempty"`
	Ores     bool `json:"ores,omitempty"`
}

type Blueprint struct {
	ID           int               `json:"id"`
	GeodeCost    Counts            `json:"geodeCost"`
	ObsidianCost Counts            `json:"obsidianCost"`
	ClayCost     Counts            `json:"clayCost"`
	OreCost      Counts            `json:"oreCost"`
	MaxOres      int               `json:"maxOres"`
	moveOn       map[Resource]bool // heuristic: do not process other resources if these are satisfied
}

type State struct {
	bp        *Blueprint
	Remaining int    `json:"remaining"`
	Resources Counts `json:"resources"`
	Bots      Counts `json:"bots"`
}

func NewState(bp *Blueprint, remaining int) State {
	return State{
		bp:        bp,
		Bots:      Counts{Ores: 1},
		Remaining: remaining,
	}
}

var errUnsatisfiable = fmt.Errorf("unsatisfiable")

func (s State) HasEnough(r Resource) bool {
	switch r {
	case Geode, Nothing:
		return false
	case Ore:
		return s.Bots.Ores >= s.bp.MaxOres
	case Clay:
		return s.Bots.Clay >= s.bp.ObsidianCost.Clay
	case Obsidian:
		return s.Bots.Obsidian >= s.bp.GeodeCost.Obsidian
	}
	panic("internal error")
}

func (s State) NextWith(bot Resource) (State, error) {
	ret := s
	var dummy int
	var c Counts
	incr := &dummy
	switch bot {
	case Nothing:
	// pass
	case Ore:
		incr = &ret.Bots.Ores
		c = ret.bp.OreCost
		if c.Ores > ret.Resources.Ores {
			return ret, errUnsatisfiable
		}
		ret.Resources.Ores -= c.Ores
	case Clay:
		incr = &ret.Bots.Clay
		c = ret.bp.ClayCost
		if c.Ores > ret.Resources.Ores {
			return ret, errUnsatisfiable
		}
		ret.Resources.Ores -= c.Ores
	case Obsidian:
		incr = &ret.Bots.Obsidian
		c = ret.bp.ObsidianCost
		if c.Clay > ret.Resources.Clay {
			return ret, errUnsatisfiable
		}
		if c.Ores > ret.Resources.Ores {
			return ret, errUnsatisfiable
		}
		ret.Resources.Clay -= c.Clay
		ret.Resources.Ores -= c.Ores
	case Geode:
		incr = &ret.Bots.Geodes
		c = ret.bp.GeodeCost
		if c.Obsidian > ret.Resources.Obsidian {
			return ret, errUnsatisfiable
		}
		if c.Ores > s.Resources.Ores {
			return ret, errUnsatisfiable
		}
		ret.Resources.Obsidian -= c.Obsidian
		ret.Resources.Ores -= c.Ores
	}

	// increment counts for next iteration
	ret.Resources.Geodes += ret.Bots.Geodes
	ret.Resources.Obsidian += ret.Bots.Obsidian
	ret.Resources.Clay += ret.Bots.Clay
	ret.Resources.Ores += ret.Bots.Ores
	// add a bot if needed
	*incr++
	// reduce iteration
	ret.Remaining--
	return ret, nil
}

func (s State) Key() string {
	return fmt.Sprintf("%v", s)
}

var blueRE = regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)

func toBlueprints(in string) []Blueprint {
	var ret []Blueprint
	for _, line := range strings.Split(strings.TrimSpace(in), "\n") {
		m := blueRE.FindStringSubmatch(line)
		if m == nil {
			panic("no RE match for line:" + line)
		}
		mustInt := func(s string) int {
			n, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			return n
		}
		b := Blueprint{
			ID:           mustInt(m[1]),
			OreCost:      Counts{Ores: mustInt(m[2])},
			ClayCost:     Counts{Ores: mustInt(m[3])},
			ObsidianCost: Counts{Ores: mustInt(m[4]), Clay: mustInt(m[5])},
			GeodeCost:    Counts{Ores: mustInt(m[6]), Obsidian: mustInt(m[7])},
		}
		maxOres := 0
		for _, cost := range []Counts{b.OreCost, b.ClayCost, b.ObsidianCost, b.GeodeCost} {
			if cost.Ores > maxOres {
				maxOres = cost.Ores
			}
		}
		b.MaxOres = maxOres
		ret = append(ret, b)
	}
	return ret
}

var counter, hits int

func bestState(state State, cache map[string]*State) (ret State) {
	key := state.Key()
	if seen, ok := cache[key]; ok {
		hits++
		if hits%1000000 == 0 {
			fmt.Fprint(os.Stderr, fmt.Sprintf("H%d", hits/1000000), "... ")
		}
		return *seen
	}
	defer func() {
		cache[key] = &ret
	}()
	if state.Remaining == 0 {
		return state
	}
	counter++
	if counter%1000000 == 0 {
		fmt.Fprint(os.Stderr, counter/1000000, "... ")
	}

	var candidates []State
	for _, b := range []Resource{Geode, Obsidian, Clay, Ore, Nothing} {
		if state.HasEnough(b) {
			continue
		}
		next, err := state.NextWith(b)
		if err == nil {
			candidates = append(candidates, next)
			if state.bp.moveOn[b] { // no need to test other stuff since this is already a best path
				break
			}
		}
	}
	var outputStates []State
	for _, c := range candidates {
		best := bestState(c, cache)
		outputStates = append(outputStates, best)
	}
	sort.Slice(outputStates, func(i, j int) bool {
		return outputStates[i].Resources.Geodes > outputStates[j].Resources.Geodes
	})
	return outputStates[0]
}

func runBlueprint(b Blueprint, iterations int) int {
	counter = 0
	state := NewState(&b, iterations)

	log.Println("Blueprint", b.ID)
	log.Println("----------")
	x, _ := json.MarshalIndent(b, "", "  ")
	log.Println(string(x))
	best := bestState(state, map[string]*State{})
	log.Println()
	log.Println(best)
	log.Println("tried", counter, "possibilities")
	return best.Resources.Geodes
}

func runP1(in string) int {
	log.SetFlags(0)
	remaining := 24
	prints := toBlueprints(in)

	output := 0
	for _, p := range prints {
		p.moveOn = map[Resource]bool{Geode: true}
		geodes := runBlueprint(p, remaining)
		log.Println()
		output += geodes * p.ID
		log.Printf("G: %d, I: %d, o: %d, O: %d", geodes, p.ID, geodes*p.ID, output)
	}
	return output
}

func runP2(in string) int {
	log.SetFlags(0)
	remaining := 32
	prints := toBlueprints(in)[:3]
	output := 1
	for _, p := range prints {
		p.moveOn = map[Resource]bool{Geode: true, Obsidian: true}
		geodes := runBlueprint(p, remaining)
		log.Println()
		output *= geodes
		log.Printf("G: %d, I: %d,  O: %d", geodes, p.ID, output)
	}
	return output
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
