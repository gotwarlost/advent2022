package dec13

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

// listishRE is a regular expression that matches 0 or more left brackets, 0 or one digit sequence, and 0 or more right brackets
var listishRE = regexp.MustCompile(`^(\[*)(\d*)(]*)$`)

// Comparison is the result of a comparison.
type Comparison int

const (
	Equal Comparison = iota
	Lesser
	Greater
)

// Kind is the kind of thing we are dealing with
type Kind int

// types of things
const (
	KindList Kind = iota
	KindInteger
)

// listOrInt is either a list of listOrInt values or a single integer
type listOrInt struct {
	parent  *listOrInt   // the listOrInt that contains this one, nil for root
	kind    Kind         // kind of thing
	n       int          // integer value if int
	l       []*listOrInt // list of stuff if list
	special bool         // a marker for pre-specified values
}

func (l *listOrInt) writeValue(w io.Writer) {
	write := func(s string) {
		_, _ = w.Write([]byte(s))
	}
	if l.kind == KindInteger {
		write(fmt.Sprint(l.n))
		write(",")
		return
	}
	write("[ ")
	for _, x := range l.l {
		x.writeValue(w)
	}
	write(" ]")
}

func (l *listOrInt) String() string {
	var x bytes.Buffer
	l.writeValue(&x)
	return x.String()
}

// cmp compares this listOrInt with another
func (l *listOrInt) cmp(other *listOrInt) Comparison {
	switch {
	// if both integers, do a simple compare
	case l.kind == KindInteger && other.kind == KindInteger:
		switch {
		case l.n > other.n:
			return Greater
		case l.n < other.n:
			return Lesser
		default:
			return Equal
		}

	// if both lists, compare an element at a time. If everything equal and this is longer return Greater.
	// If everything equal but other is longer return Lesser.
	case l.kind == KindList && other.kind == KindList:
		for i := 0; i < len(l.l); i++ {
			if i >= len(other.l) {
				return Greater
			}
			c := l.l[i].cmp(other.l[i])
			switch c {
			case Lesser:
				return Lesser
			case Greater:
				return Greater
			default:
				// fall through
			}
		}
		if len(other.l) > len(l.l) {
			return Lesser
		}
		return Equal

	// mismatched types; promote integer to a list and resume comparison
	case l.kind == KindInteger:
		s := &listOrInt{kind: KindList, l: []*listOrInt{{kind: KindInteger, n: l.n}}}
		return s.cmp(other)
	default:
		s := &listOrInt{kind: KindList, l: []*listOrInt{{kind: KindInteger, n: other.n}}}
		return l.cmp(s)
	}
}

// parseLine parses a single line into a list or int
func parseLine(line string) *listOrInt {
	var ret *listOrInt

	if !(strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")) {
		panic(fmt.Errorf("bad line: %s", line))
	}

	parts := strings.Split(line, ",")
	var curr *listOrInt
	for _, p := range parts {
		m := listishRE.FindStringSubmatch(p)
		if m == nil {
			panic(fmt.Errorf("precondition violated: %s", p))
		}
		// process the [[[[ by pushing lists into current
		for i := 0; i < len(m[1]); i++ {
			if curr == nil {
				curr = &listOrInt{kind: KindList}
				ret = curr
			} else {
				if curr.kind != KindList {
					panic("oops")
				}
				next := &listOrInt{kind: KindList, parent: curr}
				curr.l = append(curr.l, next)
				curr = next
			}
		}

		// process the numbers by pushing them into current list
		if len(m[2]) > 0 {
			n, err := strconv.Atoi(m[2])
			if err != nil {
				panic(err)
			}
			if curr.kind != KindList {
				panic("oops")
			}
			curr.l = append(curr.l, &listOrInt{kind: KindInteger, n: n})
		}

		// process the ]]]] by popping them off current
		for i := 0; i < len(m[3]); i++ {
			if curr == nil {
				panic("too many pops")
			}
			if curr.kind != KindList {
				panic("precondition failed")
			}
			curr = curr.parent
		}
	}
	if curr != nil {
		panic("not everything was popped")
	}
	if ret == nil {
		panic("ret was nil")
	}
	return ret
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
		comparison := x.cmp(y)
		oneBased := index + 1
		if comparison != Greater {
			correct = append(correct, oneBased)
			sum += oneBased
		}
	}
	return sum
}

func runP2(in string) int {
	s1, s2 := parseLine("[[2]]"), parseLine("[[6]]")
	s1.special = true
	s2.special = true
	signals := []*listOrInt{s1, s2}
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		signals = append(signals, parseLine(line))
	}
	// sort the buggers by using the comparator
	sort.Slice(signals, func(i, j int) bool {
		left := signals[i]
		right := signals[j]
		if left.cmp(right) == Greater {
			return false
		}
		return true
	})

	// find the specials
	var specials []int
	for i, s := range signals {
		if s.special {
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
