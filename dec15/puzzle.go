package dec15

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type minMax struct {
	min int
	max int
}

//go:embed input.txt
var input string

var reBeacon = regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)

type position struct {
	x, y int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p position) manhattanDistance(other position) int {
	xDiff := abs(p.x - other.x)
	yDiff := abs(p.y - other.y)
	return xDiff + yDiff
}

type pair struct {
	sensor position
	beacon position
}

var blank = minMax{-1, -1}

// minMax returns a min and max X coordinate that cannot possibly be beacons for a given y coordinate
// It returns blank if x coordinates cannot be eliminated
func (p pair) minMax(y int) minMax {
	d := p.sensor.manhattanDistance(p.beacon)
	yDiff := abs(p.sensor.y - y)
	d = d - yDiff
	if d < 0 {
		return blank
	}
	a, b := p.sensor.x-d, p.sensor.x+d
	return minMax{min: a, max: b}
}

type puzzleInput struct {
	pairs                  []pair
	minX, minY, maxX, maxY int
	beaconPositions        map[position]bool
}

func parseInput(in string) puzzleInput {
	var (
		minX = 1000000
		minY = 1000000
		maxX = -1000000
		maxY = -1000000
	)
	var pairs []pair
	positions := map[position]bool{}
	for _, line := range strings.Split(strings.TrimSpace(in), "\n") {
		m := reBeacon.FindStringSubmatch(line)
		if len(m) == 0 {
			panic("no match:" + line)
		}
		var sensor, beacon position
		var e1, e2, e3, e4 error
		sensor.x, e1 = strconv.Atoi(m[1])
		sensor.y, e2 = strconv.Atoi(m[2])
		beacon.x, e3 = strconv.Atoi(m[3])
		beacon.y, e4 = strconv.Atoi(m[4])
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			panic("bad line:" + line)
		}
		pairs = append(pairs, pair{
			sensor: sensor,
			beacon: beacon,
		})
		positions[beacon] = true

		if sensor.x < minX {
			minX = sensor.x
		}
		if beacon.x < minX {
			minX = beacon.x
		}

		if sensor.y < minY {
			minY = sensor.y
		}
		if beacon.y < minY {
			minY = beacon.y
		}

		if sensor.x > maxX {
			maxX = sensor.x
		}
		if beacon.x > maxX {
			maxX = beacon.x
		}

		if sensor.y > maxY {
			maxY = sensor.y
		}
		if beacon.y > maxY {
			maxY = beacon.y
		}
	}
	return puzzleInput{
		pairs:           pairs,
		beaconPositions: positions,
		minX:            minX,
		minY:            minY,
		maxX:            maxX,
		maxY:            maxY,
	}
}

func impossibleBeacon(pin puzzleInput, pos position) bool {
	for _, p := range pin.pairs {
		if pin.beaconPositions[pos] {
			return false
		}
		baseD := p.sensor.manhattanDistance(p.beacon)
		currD := p.sensor.manhattanDistance(pos)
		if currD <= baseD { // not possible
			return true
		}
	}
	return false
}

func possibleHiddenBeacon(pin puzzleInput, pos position) bool {
	for _, p := range pin.pairs {
		if pin.beaconPositions[pos] {
			return false
		}
		baseD := p.sensor.manhattanDistance(p.beacon)
		currD := p.sensor.manhattanDistance(pos)
		if currD <= baseD { // not possible
			return false
		}
	}
	return true
}

func runP1(in string, y int) int {
	pin := parseInput(in)
	count := 0
	for x := pin.minX - 2000000; x <= pin.maxX+2000000; x++ {
		pos := position{x: x, y: y}
		if impossibleBeacon(pin, pos) {
			count++
		}
	}
	return count
}

const max = 4000000

func runP2(in string) int {
	pin := parseInput(in)

	var maxX, maxY int
	for _, p := range pin.pairs {
		if p.sensor.x > maxX {
			maxX = p.sensor.x
		}
		if p.sensor.y > maxY {
			maxY = p.sensor.y
		}
	}

outer:
	// fix this crazy shit
	for y := 0; y <= maxY; y++ {
		var mmList []minMax
		for _, p := range pin.pairs {
			mm := p.minMax(y)
			if mm != blank {
				if mm.min < 0 {
					mm.min = 0
				}
				if mm.max > maxX {
					mm.max = maxX
				}
				mmList = append(mmList, mm)
			}
		}
		// log.Println("Y=", y, "MMLIST=", mmList)
		sort.Slice(mmList, func(i, j int) bool {
			l := mmList[i]
			r := mmList[j]
			return l.min < r.min
		})
		// log.Println("Y:", y, "MMLIST:", mmList)
		x := 0

		for {
			move := func(x int) int {
				for {
					prevX := x
					for _, mm := range mmList {
						if x >= mm.min && x <= mm.max {
							x = mm.max + 1
						}
					}
					if prevX == x {
						return x
					}
				}
			}
			x = move(x)
			if x > maxX {
				continue outer
			}
			if possibleHiddenBeacon(pin, position{x: x, y: y}) {
				log.Println("X:", x, "Y:", y)
				return x*4000000 + y
			}
			x++
		}
	}
	panic("nf")
}

func RunP1() {
	fmt.Println(runP1(input, 1000000))
}

func RunP2() {
	fmt.Println(runP2(input))
}
