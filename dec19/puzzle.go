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
	Ores  int   `json:"ores"`
}

type BotSpec struct {
	Bot   Robot `json:"bot"`
	Needs Need  `json:"needs"`
}

type BlueprintX struct {
	ID      int       `json:"id"`
	Bots    []BotSpec `json:"bots"`
	MaxOres int       `json:"maxOres"`
}

func (b BlueprintX) Needs(t Robot) Need {
	if t.Product() == Nothing {
		return Need{}
	}
	for _, c := range b.Bots {
		if c.Bot == t {
			return c.Needs
		}
	}
	panic("must find dep")
}

type StateX struct {
	bp        *BlueprintX
	Remaining int            `json:"remaining"`
	Resources map[Thing]int  `json:"resources"`
	Bots      map[Robot]int  `json:"bots"`
	Skipped   map[Robot]bool `json:"skipped"`
}

func newState(bp *BlueprintX, remaining int) StateX {
	return StateX{
		bp:        bp,
		Remaining: remaining,
		Resources: map[Thing]int{},
		Bots:      map[Robot]int{},
		Skipped:   map[Robot]bool{},
	}
}

func (s StateX) clone() StateX {
	ret := newState(s.bp, s.Remaining)
	for k, v := range s.Resources {
		ret.Resources[k] = v
	}
	for k, v := range s.Bots {
		ret.Bots[k] = v
	}
	for k, v := range s.Skipped {
		ret.Skipped[k] = v
	}
	return ret
}

func (s StateX) ShouldSkip(bot Robot) bool {
	// never skip these
	if bot == NothingBot || bot == GeodeBot {
		return false
	}
	if s.Skipped[bot] { // skipped once, skip forever
		return true
	}
	switch bot {
	case OreBot:
		return s.Bots[OreBot] >= s.bp.MaxOres
	case ClayBot:
		return s.Bots[ClayBot] >= s.bp.Needs(ObsidianBot).Count
	case ObsidianBot:
		return s.Bots[ObsidianBot] >= s.bp.Needs(GeodeBot).Count
	}
	panic("unreachable")
}

type cacheKey struct {
	remaining    int
	geodes       int
	obsidian     int
	clay         int
	ores         int
	geodeBots    int
	obsidianBots int
	clayBots     int
	oreBots      int
}

func (s StateX) key() cacheKey {
	return cacheKey{
		remaining:    s.Remaining,
		geodes:       s.Resources[Geode],
		obsidian:     s.Resources[Obsidian],
		clay:         s.Resources[Clay],
		ores:         s.Resources[Ore],
		geodeBots:    s.Bots[GeodeBot],
		obsidianBots: s.Bots[ObsidianBot],
		clayBots:     s.Bots[ClayBot],
		oreBots:      s.Bots[OreBot],
	}
}

type botErr struct {
	unsatisfiable Thing
}

func (b botErr) Error() string {
	return fmt.Sprintf("unsatisfiable %q", b.unsatisfiable)
}

func (s StateX) NextWith(bot Robot) (StateX, error) {
	ret := s.clone()
	needs := ret.bp.Needs(bot)
	var uns []Thing
	thing := needs.Thing
	if thing != Nothing {
		want := needs.Count
		have := ret.Resources[thing]
		if have < want {
			uns = append(uns, thing)
		}
		ret.Resources[thing] -= want
	}
	if needs.Ores > 0 {
		have := ret.Resources[Ore]
		if have < needs.Ores {
			uns = append(uns, thing)
		}
		ret.Resources[Ore] -= needs.Ores
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

func (s StateX) String() string {
	return fmt.Sprintf("\tRemaining: %d\n\tResources: %v\n\tRobots: %v", s.Remaining, s.Resources, s.Bots)
}

var blueRE = regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)

func toBlueprints(in string) []BlueprintX {
	var ret []BlueprintX
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
		b := BlueprintX{
			ID: mustInt(m[1]),
			Bots: []BotSpec{
				{
					Bot: OreBot,
					Needs: Need{
						Ores: mustInt(m[2]),
					},
				},
				{
					Bot: ClayBot,
					Needs: Need{
						Ores: mustInt(m[3]),
					},
				},
				{
					Bot: ObsidianBot,
					Needs: Need{
						Ores:  mustInt(m[4]),
						Thing: Clay,
						Count: mustInt(m[5]),
					},
				},
				{
					Bot: GeodeBot,
					Needs: Need{
						Ores:  mustInt(m[6]),
						Thing: Obsidian,
						Count: mustInt(m[7]),
					},
				},
			},
		}
		sort.Slice(b.Bots, func(i, j int) bool {
			left := b.Bots[i]
			right := b.Bots[j]
			return left.Bot.Value() > right.Bot.Value()
		})
		maxOres := 0
		for _, bot := range b.Bots {
			if bot.Needs.Ores > maxOres {
				maxOres = bot.Needs.Ores
			}
		}
		b.MaxOres = maxOres
		ret = append(ret, b)
	}
	return ret
}

var counter, hits int

func bestState(state StateX, cache map[cacheKey]StateX) StateX {
	addToCache := func(s StateX) StateX {
		cache[s.key()] = s
		return s
	}
	if state.Remaining == 0 {
		return addToCache(state)
	}
	cs := state.key()
	if s, ok := cache[cs]; ok {
		hits++
		if hits%1000 == 0 {
			fmt.Print("H", hits/1000, "... ")
		}
		return s
	}
	counter++
	if counter%1000000 == 0 {
		fmt.Fprint(os.Stderr, counter/1000000, "... ")
	}

	clone := state.clone()
	for _, c := range []Robot{ObsidianBot, ClayBot, OreBot} {
		if state.ShouldSkip(c) {
			clone.Skipped[c] = true
		}
	}
	state = clone

	bots := []Robot{GeodeBot, ObsidianBot, ClayBot, OreBot, NothingBot}
	var candidates []StateX
	for _, b := range bots {
		if state.Skipped[b] {
			continue
		}
		next, err := state.NextWith(b)
		if err == nil {
			if b == GeodeBot {
				return bestState(next, cache)
			}
			candidates = append(candidates, next)
		}
	}
	var outputStates []StateX
	for _, c := range candidates {
		best := addToCache(bestState(c, cache))
		outputStates = append(outputStates, best)
	}
	sort.Slice(outputStates, func(i, j int) bool {
		return outputStates[i].Resources[Geode] > outputStates[j].Resources[Geode]
	})
	return outputStates[0]
}

func runBlueprint(b BlueprintX, iterations int) int {
	counter = 0
	state := newState(&b, iterations)
	state.Bots = map[Robot]int{
		Ore.toBot(): 1,
	}

	log.Println("Blueprint", b.ID)
	log.Println("----------")
	x, _ := json.MarshalIndent(b, "", "  ")
	log.Println(string(x))
	best := bestState(state, map[cacheKey]StateX{})
	log.Println()
	log.Println(best)
	log.Println("tried", counter, "possibilities")
	return best.Resources[Geode]
}

func runP1(in string) int {
	log.SetFlags(0)
	mins := 24
	prints := toBlueprints(in)

	output := 0
	for _, p := range prints {
		geodes := runBlueprint(p, mins)
		log.Println()
		output += geodes * p.ID
		log.Printf("G: %d, I: %d, o: %d, O: %d", geodes, p.ID, geodes*p.ID, output)
	}
	return output
}

func runP2(in string) int {
	return 0
}

func RunP1() {
	solve()
}

func RunP2() {
	fmt.Println(runP2(input))
}
