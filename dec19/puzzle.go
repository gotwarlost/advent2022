package dec19

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type Thing string

const (
	Nothing  Thing = ""
	Ore      Thing = "ore"
	Clay     Thing = "clay"
	Obsidian Thing = "obsidian"
	Geode    Thing = "geode"
)

func (t Thing) String() string {
	return string(t)
}

func (t Thing) Value() int {
	switch t {
	case Nothing:
		return 0
	case Ore:
		return 1
	case Clay:
		return 2
	case Obsidian:
		return 3
	default:
		return 4
	}
}

func (t Thing) toBot() Robot {
	return Robot(t.String())
}

type Robot string

func (t Robot) String() string {
	return string(t) + "-bot"
}

func (t Robot) Product() Thing {
	return Thing(t)
}

func (t Robot) Value() int {
	return Thing(t).Value()
}

var (
	NothingBot  = Nothing.toBot()
	OreBot      = Ore.toBot()
	ClayBot     = Clay.toBot()
	ObsidianBot = Obsidian.toBot()
	GeodeBot    = Geode.toBot()
)

type Need struct {
	Thing Thing `json:"thing"`
	Count int   `json:"count"`
}

type BotSpec struct {
	Bot   Robot         `json:"bot"`
	Needs map[Thing]int `json:"needs"`
}

type Blueprint struct {
	ID   int       `json:"id"`
	Bots []BotSpec `json:"bots"`
}

func (b Blueprint) Needs(t Robot) map[Thing]int {
	if t.Product() == Nothing {
		return map[Thing]int{}
	}
	for _, c := range b.Bots {
		if c.Bot == t {
			return c.Needs
		}
	}
	panic("must find dep")
}

type State struct {
	print     *Blueprint
	Remaining int           `json:"remaining"`
	Resources map[Thing]int `json:"resources"`
	Bots      map[Robot]int `json:"bots"`
}

func newState(print *Blueprint, remaining int) State {
	return State{
		print:     print,
		Remaining: remaining,
		Resources: map[Thing]int{},
		Bots:      map[Robot]int{},
	}
}

func (s State) clone() State {
	ret := newState(s.print, s.Remaining)
	for k, v := range s.Resources {
		ret.Resources[k] = v
	}
	for k, v := range s.Bots {
		ret.Bots[k] = v
	}
	return ret
}

type botErr struct {
	unsatisfiable Thing
}

func (b botErr) Error() string {
	return fmt.Sprintf("unsatisfiable %q", b.unsatisfiable)
}

func (s State) NextWith(bot Robot) (State, error) {
	ret := s.clone()
	needs := ret.print.Needs(bot)
	var uns []Thing
	for thing, want := range needs {
		have := ret.Resources[thing]
		if have < want {
			uns = append(uns, thing)
		}
		ret.Resources[thing] -= want
	}
	if len(uns) > 0 {
		sort.Slice(uns, func(i, j int) bool {
			return uns[i].Value() > uns[j].Value()
		})
		return ret, &botErr{unsatisfiable: uns[0]}
	}
	for k, v := range ret.Bots {
		ret.Resources[k.Product()] += v
	}
	if bot != NothingBot {
		ret.Bots[bot]++
	}
	ret.Remaining--
	return ret, nil
}

func (s State) String() string {
	return fmt.Sprintf("\tRemaining: %d\n\tResources: %v\n\tRobots: %v", s.Remaining, s.Resources, s.Bots)
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
			ID: mustInt(m[1]),
			Bots: []BotSpec{
				{
					Bot: OreBot,
					Needs: map[Thing]int{
						Ore: mustInt(m[2]),
					},
				},
				{
					Bot: ClayBot,
					Needs: map[Thing]int{
						Ore: mustInt(m[3]),
					},
				},
				{
					Bot: ObsidianBot,
					Needs: map[Thing]int{
						Ore:  mustInt(m[4]),
						Clay: mustInt(m[5]),
					},
				},
				{
					Bot: GeodeBot,
					Needs: map[Thing]int{
						Ore:      mustInt(m[6]),
						Obsidian: mustInt(m[7]),
					},
				},
			},
		}
		sort.Slice(b.Bots, func(i, j int) bool {
			left := b.Bots[i]
			right := b.Bots[j]
			return left.Bot.Value() > right.Bot.Value()
		})
		ret = append(ret, b)
	}
	return ret
}

var counter int

func bestPath(state State) State {
	if state.Remaining == 0 {
		return state
	}
	counter++
	if counter%1000000 == 0 {
		fmt.Print(counter/1000000, "... ")
	}

	nextOres := func() State {
		s2, _ := state.NextWith(NothingBot)
		s2 = bestPath(s2)
		next, err := state.NextWith(OreBot)
		if err != nil {
			return s2
		}
		s1 := bestPath(next)
		if s1.Resources[Geode] > s2.Resources[Geode] {
			return s1
		} else {
			return s2
		}
	}

	bot := GeodeBot
	for {
		next, err := state.NextWith(bot)
		if err == nil {
			return bestPath(next)
		}
		bot = err.(*botErr).unsatisfiable.toBot()
		if bot == OreBot {
			return nextOres()
		}
	}
}

func runBlueprint(b Blueprint, iterations int) int {
	counter = 0
	state := newState(&b, iterations)
	state.Bots = map[Robot]int{
		Ore.toBot(): 1,
	}

	log.Println("Blueprint", b.ID)
	log.Println("----------")
	best := bestPath(state)
	log.Println(best)
	log.Println("tried", counter, "possibilities")
	return best.Resources[Geode]
}

func runP1(in string) int {
	log.SetFlags(0)
	mins := 24
	prints := toBlueprints(in)
	x, _ := json.MarshalIndent(prints, "", "  ")
	log.Println(string(x))

	output := 0
	for _, p := range prints {
		geodes := runBlueprint(p, mins)
		fmt.Println()
		output += geodes * p.ID
		log.Printf("G: %d, I: %d, o: %d, O: %d", geodes, p.ID, geodes*p.ID, output)
	}
	return output
}

func runP2(in string) int {
	return 0
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
