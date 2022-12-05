package dec5

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type stack struct {
	ids []byte
}

var moveRE = regexp.MustCompile(`^move\s+(\d+)\s+from\s+(\d+)\s+to\s+(\d+)$`)

func run(in string, moveFn func(source, dest *stack, num int)) string {
	max := 10
	stacks := make([]*stack, max)
	for i := 0; i < max; i++ {
		stacks[i] = &stack{}
	}

	lines := strings.Split(strings.Trim(in, "\n"), "\n")
	for _, line := range lines {
		switch {
		case strings.TrimSpace(line) == "":

		case strings.HasPrefix(strings.TrimSpace(line), "1"):
		// noop

		case strings.HasPrefix(strings.TrimSpace(line), "["):
			length := len(line)
			if (length+1)%4 != 0 {
				panic("invalid line:" + line + fmt.Sprintf("(length: %d", length))
			}
			for i := 0; i <= length; i += 4 {
				stackNum := i / 4
				ch := line[i+1]
				if ch != ' ' {
					stacks[stackNum].ids = append(stacks[stackNum].ids, ch)
				}
			}

		case strings.HasPrefix(line, "move"):
			m := moveRE.FindStringSubmatch(line)
			if m == nil {
				panic("line:" + line + ", no move match")
			}
			num, e1 := strconv.Atoi(m[1])
			source, e2 := strconv.Atoi(m[2])
			dest, e3 := strconv.Atoi(m[3])
			if e1 != nil || e2 != nil || e3 != nil {
				panic(fmt.Errorf("convert: %s, %v %v %v", line, e1, e2, e3))
			}
			// 0-index
			sourceStack := stacks[source-1]
			destStack := stacks[dest-1]
			if len(sourceStack.ids) < num {
				panic("too few ids to transfer")
			}
			moveFn(sourceStack, destStack, num)
		default:
			panic("unknown line:" + line)
		}
	}

	var out []byte
	for _, s := range stacks {
		if len(s.ids) > 0 {
			out = append(out, s.ids[0])
		}
	}
	return string(out)
}

func runPart1(in string) string {
	return run(in, func(source, dest *stack, num int) {
		for x := 0; x < num; x++ {
			crate := source.ids[0]
			source.ids = source.ids[1:]
			dest.ids = append([]byte{crate}, dest.ids...)
		}
	})
}

func runPart2(in string) string {
	return run(in, func(source, dest *stack, num int) {
		slice := source.ids[0:num]
		source.ids = source.ids[num:]
		for i := len(slice) - 1; i >= 0; i-- {
			dest.ids = append([]byte{slice[i]}, dest.ids...)
		}
	})
}

func RunPart1() {
	fmt.Println("Output:", runPart1(input))
}

func RunPart2() {
	fmt.Println("Output:", runPart2(input))
}
