package dec13

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

// comparison is the result of a comparison.
type comparison int

const (
	eq comparison = iota
	lesser
	greater
)

// nodeKind is the kind of node
type nodeKind int

// kinds of nodes
const (
	kindContainer nodeKind = iota
	kindLeaf
)

// node is either a container or a leaf
type node struct {
	kind     nodeKind // kind of node
	value    int      // integer value if int
	children []*node  // children if list
}

func (n *node) writeValue(w io.Writer) {
	write := func(s string) {
		_, _ = w.Write([]byte(s))
	}
	if n.kind == kindLeaf {
		write(fmt.Sprint(n.value))
		return
	}
	write("[")
	for i, x := range n.children {
		if i > 0 {
			write(",")
		}
		x.writeValue(w)
	}
	write("]")
}

func (n *node) String() string {
	var x bytes.Buffer
	n.writeValue(&x)
	return x.String()
}

// cmp compares this node with another
func (n *node) cmp(other *node) comparison {
	switch {
	// if both leaves, do a simple compare
	case n.kind == kindLeaf && other.kind == kindLeaf:
		switch {
		case n.value > other.value:
			return greater
		case n.value < other.value:
			return lesser
		default:
			return eq
		}

	// if both containers, compare an element at a time. If everything equal and this is longer return greater.
	// If everything equal but other is longer return lesser.
	case n.kind == kindContainer && other.kind == kindContainer:
		for i := 0; i < len(n.children); i++ {
			if i >= len(other.children) {
				return greater
			}
			c := n.children[i].cmp(other.children[i])
			switch c {
			case lesser, greater:
				return c
			default:
				// fall through
			}
		}
		if len(other.children) > len(n.children) {
			return lesser
		}
		return eq

	// mismatched types; promote leaf to container and resume comparison
	case n.kind == kindLeaf:
		return createContainer(n).cmp(other)
	default:
		return n.cmp(createContainer(other))
	}
}

func (n *node) addChild(child *node) {
	n.children = append(n.children, child)
}

func createContainer(children ...*node) *node {
	ret := &node{kind: kindContainer}
	for _, c := range children {
		ret.addChild(c)
	}
	return ret
}

type tokenKind int

const (
	_ tokenKind = iota
	tokOpen
	tokNumber
	tokClose
)

type token struct {
	kind  tokenKind
	value int
}

func tokenize(line []byte) []token {
	isNum := func(b byte) bool {
		return b >= '0' && b <= '9'
	}

	var ret []token
	pos := 0
	for {
		if pos == len(line) {
			break
		}
		ch := line[pos]
		switch {
		case ch == '[':
			ret = append(ret, token{kind: tokOpen})
			pos++
		case ch == ']':
			ret = append(ret, token{kind: tokClose})
			pos++
		case isNum(ch):
			var endPos int
			for endPos = pos + 1; endPos < len(line); endPos++ {
				if !isNum(line[endPos]) {
					break
				}
			}
			numStr := string(line[pos:endPos])
			n, err := strconv.Atoi(numStr)
			if err != nil {
				panic(err)
			}
			ret = append(ret, token{kind: tokNumber, value: n})
			pos = endPos
		default:
			pos++
		}
	}
	return ret
}

type nodeStack struct {
	nodes []*node
}

func (s *nodeStack) isEmpty() bool   { return len(s.nodes) == 0 }
func (s *nodeStack) push(node *node) { s.nodes = append(s.nodes, node) }
func (s *nodeStack) peek() *node     { return s.nodes[len(s.nodes)-1] }

func (s *nodeStack) pop() *node {
	top := len(s.nodes) - 1
	n := s.nodes[top]
	s.nodes = s.nodes[:top]
	return n
}

func parseLine(line string) *node {
	var root *node
	s := nodeStack{}
	tokenValues := tokenize([]byte(line))

	for _, tok := range tokenValues {
		switch tok.kind {
		case tokOpen:
			c := createContainer()
			if s.isEmpty() {
				root = c
			} else {
				s.peek().addChild(c)
			}
			s.push(c)
		case tokClose:
			s.pop()
		default:
			s.peek().addChild(&node{kind: kindLeaf, value: tok.value})
		}
	}
	if !s.isEmpty() {
		panic("stack not fully popped")
	}
	return root
}

func runP1(in string) int {
	blocks := strings.Split(strings.TrimSpace(in), "\n\n")
	var correct []int
	sum := 0
	for index, block := range blocks {
		lines := strings.Split(block, "\n")
		if len(lines) != 2 {
			panic(fmt.Errorf("bad lines, %v", block))
		}
		x, y := parseLine(lines[0]), parseLine(lines[1])
		c := x.cmp(y)
		oneBased := index + 1
		if c != greater {
			correct = append(correct, oneBased)
			sum += oneBased
		}
	}
	return sum
}

func runP2(in string) int {
	s1, s2 := parseLine("[[2]]"), parseLine("[[6]]")
	signals := []*node{s1, s2}
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		signals = append(signals, parseLine(line))
	}
	// sort the signals by using the comparator
	sort.Slice(signals, func(i, j int) bool {
		left := signals[i]
		right := signals[j]
		if left.cmp(right) == greater {
			return false
		}
		return true
	})

	// find the specials
	var specials []int
	for i, s := range signals {
		if s == s1 || s == s2 {
			specials = append(specials, i+1)
		}
	}
	if len(specials) != 2 {
		panic("bad specials")
	}
	return specials[0] * specials[1]
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
