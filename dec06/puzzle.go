package dec06

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

// uniqueWindowChecker maintains counts of letters seen in a map
// ensuring that when a count drops to 0 the entry is removed.
// Thus, checks for unique of all entries in flight is a simple
// comparison between window size and the size of the map.
type uniqueWindowChecker struct {
	size int
	seen map[byte]int
}

// newUniqueWindowChecker sets up a window checker with bytes from the initial window.
func newUniqueWindowChecker(initial string) *uniqueWindowChecker {
	seen := map[byte]int{}
	for _, b := range []byte(initial) {
		seen[b]++
	}
	return &uniqueWindowChecker{
		size: len(initial),
		seen: seen,
	}
}

// replace a character previously added with one the now needs to be added.
func (u *uniqueWindowChecker) replace(src, dest byte) {
	n := u.seen[src]
	if n == 0 {
		panic("precondition failed")
	}
	n--
	if n == 0 {
		delete(u.seen, src)
	} else {
		u.seen[src] = n
	}
	u.seen[dest]++
}

func (u *uniqueWindowChecker) isUnique() bool {
	return len(u.seen) == u.size
}

func nonRepeatingChars(in string, window int) int {
	in = strings.TrimSpace(in)
	initial := in[:window]
	uc := newUniqueWindowChecker(initial)
	if uc.isUnique() {
		return window
	}
	for i := window; i < len(in); i++ {
		uc.replace(in[i-window], in[i])
		if uc.isUnique() {
			return i + 1
		}
	}
	panic("not found")
}

func RunPart1() {
	fmt.Println("Output:", nonRepeatingChars(input, 4))
}

func RunPart2() {
	fmt.Println("Output:", nonRepeatingChars(input, 14))
}
