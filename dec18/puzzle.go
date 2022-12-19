package dec18

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type plane int

const (
	_ plane = iota
	xy
	xz
	yz
)

func (p plane) String() string {
	switch p {
	case xy:
		return "xy"
	case xz:
		return "xz"
	default:
		return "yz"
	}
}

type surface struct {
	plane   plane
	x, y, z int
}

func (s surface) String() string {
	return fmt.Sprintf("%v (%d,%d,%d)", s.plane, s.x, s.y, s.z)
}

type cube struct {
	x, y, z int
}

func (c cube) String() string {
	return fmt.Sprintf("(%d,%d,%d)", c.x, c.y, c.z)
}

func (c cube) top() surface    { return surface{xz, c.x, c.y + 1, c.z} }
func (c cube) bottom() surface { return surface{xz, c.x, c.y, c.z} }
func (c cube) left() surface   { return surface{yz, c.x, c.y, c.z} }
func (c cube) right() surface  { return surface{yz, c.x + 1, c.y, c.z} }
func (c cube) front() surface  { return surface{xy, c.x, c.y, c.z} }
func (c cube) back() surface   { return surface{xy, c.x, c.y, c.z + 1} }

func (c cube) surfaces() []surface {
	return []surface{c.top(), c.bottom(), c.left(), c.right(), c.front(), c.back()}
}

func toCubes(in string) []cube {
	var ret []cube
	for _, line := range strings.Split(strings.TrimSpace(in), "\n") {
		xyz := strings.Split(line, ",")
		x, e1 := strconv.Atoi(xyz[0])
		y, e2 := strconv.Atoi(xyz[1])
		z, e3 := strconv.Atoi(xyz[2])
		if e1 != nil || e2 != nil || e3 != nil {
			panic("bad")
		}
		ret = append(ret, cube{x, y, z})
	}
	return ret
}

type segment struct {
	direction  plane   // direction is perpendicular to the plane
	begin, end surface // the surfaces at  begin and end
}

func (s segment) String() string {
	if s == notFound {
		return "<no-segment>"
	}
	return fmt.Sprintf("%v: %v -> %v", s.direction, s.begin, s.end)
}

var notFound = segment{}

func run(in string, checkEnclosed bool) int {
	minCube := cube{1000000, 1000000, 1000000}
	maxCube := cube{-1000000, -1000000, -1000000}

	cubes := toCubes(in)
	cubeLocations := map[cube]bool{}
	surfaceLocations := map[surface]bool{}
	seen := map[surface]int{}

	setMinMax := func(c cube) {
		if minCube.x > c.x {
			minCube.x = c.x
		}
		if maxCube.x < c.x {
			maxCube.x = c.x
		}
		if minCube.y > c.y {
			minCube.y = c.y
		}
		if maxCube.y < c.y {
			maxCube.y = c.y
		}
		if minCube.z > c.z {
			minCube.z = c.z
		}
		if maxCube.z < c.z {
			maxCube.z = c.z
		}
	}

	for _, c := range cubes {
		cubeLocations[c] = true
		setMinMax(c)
		for _, s := range c.surfaces() {
			surfaceLocations[s] = true
			seen[s]++
		}
	}
	// min bounds are the other faces of the cube
	minCube.x++
	minCube.y++
	minCube.z++
	area := 0

	for _, c := range cubes {
		count := 0
		for _, s := range c.surfaces() {
			if seen[s] == 2 {
				continue
			}
			area++
		}
		area += count
	}

	if !checkEnclosed {
		return area
	}

	segmentAlongDirection := func(c cube, p plane) (ret segment) {
		var start, end, increment int
		var ptr *int
		var fn func(c cube) surface
		test := c

		findSurface := func() (out *surface) {
			test = c
			i := start
		loop:
			for {
				*ptr = i
				s := fn(test)
				if surfaceLocations[s] {
					return &s
				}
				i += increment
				switch {
				case start > end:
					if i < end {
						break loop
					}
				case start < end:
					if i > end {
						break loop
					}
				default:
					break loop
				}
			}
			return nil
		}

		switch p {
		case xy:
			start, end, increment, ptr, fn = c.z, minCube.z, -1, &test.z, func(c cube) surface { return c.front() }
		case xz:
			start, end, increment, ptr, fn = c.y, minCube.y, -1, &test.y, func(c cube) surface { return c.bottom() }
		case yz:
			start, end, increment, ptr, fn = c.x, minCube.x, -1, &test.x, func(c cube) surface { return c.left() }
		}
		x := findSurface()
		if x == nil {
			return notFound
		}
		begin := *x
		switch p {
		case xy:
			start, end, increment, ptr, fn = c.z, maxCube.z, 1, &test.z, func(c cube) surface { return c.back() }
		case xz:
			start, end, increment, ptr, fn = c.y, maxCube.y, 1, &test.y, func(c cube) surface { return c.top() }
		case yz:
			start, end, increment, ptr, fn = c.x, maxCube.x, 1, &test.x, func(c cube) surface { return c.right() }
		}
		x = findSurface()
		if x == nil {
			return notFound
		}
		finish := *x
		ret = segment{direction: p, begin: begin, end: finish}
		return ret
	}

	basicEnclosed := map[cube][]segment{}
	for x := minCube.x; x < maxCube.x; x++ {
		for y := minCube.y; y < maxCube.y; y++ {
			for z := minCube.z; z < maxCube.z; z++ {
				test := cube{x, y, z}
				if cubeLocations[test] { // not empty at that position
					continue
				}
				xs, ys, zs := segmentAlongDirection(test, yz), segmentAlongDirection(test, xz), segmentAlongDirection(test, xy)
				if xs == notFound || ys == notFound || zs == notFound {
					continue
				}
				basicEnclosed[test] = []segment{xs, ys, zs}
			}
		}
	}

	surfacesToRemove := map[surface]bool{}
	for _, segs := range basicEnclosed {
		candidateSurfaces := map[surface]bool{}
		bottomLeft := cube{x: segs[0].begin.x, y: segs[1].begin.y, z: segs[2].begin.z}
		topRight := cube{x: segs[0].end.x, y: segs[1].end.y, z: segs[2].end.z}
		allEnclosed := true

	outermost:
		for x := bottomLeft.x; x < topRight.x; x++ {
			for y := bottomLeft.y; y < topRight.y; y++ {
				for z := bottomLeft.z; z < topRight.z; z++ {
					t2 := cube{x, y, z}
					if _, isCube := cubeLocations[t2]; isCube {
						continue
					}
					segs2, ok := basicEnclosed[t2]
					if !ok {
						allEnclosed = false
						break outermost
					}
					for _, seg := range segs2 {
						candidateSurfaces[seg.begin] = true
						candidateSurfaces[seg.end] = true
					}
				}
			}
		}
		if allEnclosed {
			for s := range candidateSurfaces {
				surfacesToRemove[s] = true
			}
		}
	}
	area -= len(surfacesToRemove)
	return area
}

func RunP1() {
	fmt.Println(run(input, false))
}

func RunP2() {
	fmt.Println(run(input, true))
}
