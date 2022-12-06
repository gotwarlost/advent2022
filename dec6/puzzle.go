package dec6

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type uniqueWindowChecker struct {
	size int
	seen map[byte]int
}

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
