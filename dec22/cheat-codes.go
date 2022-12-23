package dec22

/*
 test pattern
 	| | |A| |
	|B|C|D| |
	| | |E|F|
*/

type Edge int

const (
	eTop Edge = iota
	eBottom
	eLeft
	eRight
)

func (e Edge) String() string {
	switch e {
	case eTop:
		return "^"
	case eBottom:
		return "v"
	case eLeft:
		return "<"
	case eRight:
		return ">"
	}
	return "X"
}

type tmpRegionEdge struct {
	region string
	edge   Edge
}

var testLetterMap = map[string]string{
	"A": "02",
	"B": "10",
	"C": "11",
	"D": "12",
	"E": "22",
	"F": "23",
}

var testRoutes = map[tmpRegionEdge]tmpRegionEdge{
	{"A", eTop}:    {"B", eTop},
	{"A", eLeft}:   {"C", eTop},
	{"A", eRight}:  {"F", eRight},
	{"A", eBottom}: {"D", eTop},

	{"B", eLeft}:   {"F", eBottom},
	{"B", eRight}:  {"C", eLeft},
	{"B", eBottom}: {"E", eBottom},

	{"C", eRight}:  {"D", eLeft},
	{"C", eBottom}: {"E", eLeft},

	{"D", eBottom}: {"E", eTop},
	{"D", eRight}:  {"F", eTop},

	{"E", eRight}: {"F", eLeft},
}

/*
	 main pattern
	 	| |B|A|
		| |C| |
		|E|D| |
	    |F| | |
*/
var mainLetterMap = map[string]string{
	"A": "02",
	"B": "01",
	"C": "11",
	"D": "21",
	"E": "20",
	"F": "30",
}

var mainRoutes = map[tmpRegionEdge]tmpRegionEdge{
	{"A", eTop}:    {"F", eBottom},
	{"A", eLeft}:   {"B", eRight},
	{"A", eRight}:  {"D", eRight},
	{"A", eBottom}: {"C", eRight},

	{"B", eLeft}:   {"E", eLeft},
	{"B", eTop}:    {"F", eLeft},
	{"B", eBottom}: {"C", eTop},

	{"C", eLeft}:   {"E", eTop},
	{"C", eBottom}: {"D", eTop},

	{"D", eBottom}: {"F", eRight},
	{"D", eLeft}:   {"E", eRight},

	{"E", eBottom}: {"F", eTop},
}

type region struct {
	row, col int
}

type regionEdge struct {
	region
	edge Edge
}

func cheatRoutes(letterMap map[string]string, letterRoutes map[tmpRegionEdge]tmpRegionEdge) map[regionEdge]regionEdge {
	toGridRegion := func(s string) region {
		r, ok := letterMap[s]
		if !ok || len(r) != 2 {
			panic("internal")
		}
		row := int(r[0] - '0')
		col := int(r[1] - '0')
		return region{row, col}
	}
	ret := map[regionEdge]regionEdge{}
	for left, right := range letterRoutes {
		mLeft := regionEdge{region: toGridRegion(left.region), edge: left.edge}
		mRight := regionEdge{region: toGridRegion(right.region), edge: right.edge}
		ret[mLeft] = mRight
		ret[mRight] = mLeft
	}
	if len(ret) != 24 {
		panic("bad route count")
	}
	return ret
}

func testRouteMap() map[regionEdge]regionEdge {
	return cheatRoutes(testLetterMap, testRoutes)
}

func mainRouteMap() map[regionEdge]regionEdge {
	return cheatRoutes(mainLetterMap, mainRoutes)
}
