package dec21

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type op string

const (
	_          op = ""
	opAdd      op = "+"
	opSubtract op = "-"
	opMultiply op = "*"
	opDivide   op = "/"
)

func (o op) apply(left, right int64) int64 {
	switch o {
	case opAdd:
		return left + right
	case opSubtract:
		return left - right
	case opMultiply:
		return left * right
	case opDivide:
		return left / right
	}
	panic("boo")
}

func (o op) reverseApply(result, operand int64, operandIsLeft bool) (ret int64) {
	switch o {
	case opAdd:
		// result = operand + x
		return result - operand
	case opSubtract:
		if operandIsLeft {
			// result = operand - x
			return operand - result
		} else {
			// result = x - operand
			return result + operand
		}
	case opMultiply:
		// result = operand *  x
		return result / operand
	case opDivide:
		if operandIsLeft {
			// result = operand / x
			return operand / result
		} else {
			// result = x / operand
			return result * operand
		}
	}
	panic("boo")
}

type monkey struct {
	name  string
	expr  bool
	val   int64
	left  *monkey
	right *monkey
	op
}

type tmpMonkey struct {
	name  string
	expr  bool
	val   int64
	left  string
	right string
	op
}

var (
	reNumMonkey = regexp.MustCompile(`^(\S+): (\d+)$`)
	reOpMonkey  = regexp.MustCompile(`^(\S+): (\S+) (.) (\S+)$`)
)

func toMonkeys(in string) map[string]*monkey {
	// m := map[string]*monkey{}
	lines := strings.Split(strings.TrimSpace(in), "\n")
	tmpMap := map[string]*tmpMonkey{}
	for _, line := range lines {
		var mk *tmpMonkey
		if m := reNumMonkey.FindStringSubmatch(line); len(m) > 0 {
			n, err := strconv.Atoi(m[2])
			if err != nil {
				panic(err)
			}
			mk = &tmpMonkey{
				name: m[1],
				val:  int64(n),
			}
		} else if m := reOpMonkey.FindStringSubmatch(line); len(m) > 0 {
			var oper op
			o := m[3]
			switch o {
			case "+":
				oper = opAdd
			case "-":
				oper = opSubtract
			case "*":
				oper = opMultiply
			case "/":
				oper = opDivide
			default:
				panic("bad op:" + m[3])
			}
			mk = &tmpMonkey{
				name:  m[1],
				expr:  true,
				left:  m[2],
				op:    oper,
				right: m[4],
			}
		} else {
			panic("no match for: " + line)
		}
		tmpMap[mk.name] = mk
	}

	monkeyMap := map[string]*monkey{}
	var tmpToReal func(s string)
	tmpToReal = func(name string) {
		_, ok := monkeyMap[name]
		if ok {
			return
		}
		x := tmpMap[name]
		if !x.expr {
			monkeyMap[name] = &monkey{
				name: name,
				val:  x.val,
			}
		} else {
			tmpToReal(x.left)
			tmpToReal(x.right)
			monkeyMap[name] = &monkey{
				name:  name,
				expr:  true,
				left:  monkeyMap[x.left],
				right: monkeyMap[x.right],
				op:    x.op,
			}
		}
	}

	for k := range tmpMap {
		tmpToReal(k)
	}
	return monkeyMap
}

func toNumber(m *monkey) int64 {
	if !m.expr {
		return m.val
	}
	left := toNumber(m.left)
	right := toNumber(m.right)
	return m.op.apply(left, right)
}

func runP1(in string) int64 {
	monkeys := toMonkeys(in)
	root := monkeys["root"]
	return toNumber(root)
}

func hasDep(root *monkey, name string) bool {
	if !root.expr {
		if root.name == name {
			return true
		}
		return false
	}
	if root.left.name == name || root.right.name == name {
		return true
	}
	return hasDep(root.left, name) || hasDep(root.right, name)
}

func reverseEval(root *monkey, result int64, name string) int64 {
	l := root.left
	r := root.right
	leftDep := hasDep(l, name)
	rightDep := hasDep(r, name)
	var leftKnown bool
	var eval *monkey
	var operand int64
	if leftDep && rightDep {
		panic("human all over the place!")
	} else if leftDep {
		eval = l
		operand = toNumber(r)
		leftKnown = false
	} else if rightDep {
		eval = r
		operand = toNumber(l)
		leftKnown = true
	} else {
		panic(fmt.Sprintf("no human at reverseEval"))
	}

	newResult := root.op.reverseApply(result, operand, leftKnown)
	switch {
	case eval.name == name:
		return newResult
	default:
		return reverseEval(eval, newResult, name)
	}
}

func runP2(in string) int64 {
	human := "humn"
	monkeys := toMonkeys(in)
	l := monkeys["root"].left
	r := monkeys["root"].right

	leftDep := hasDep(l, human)
	rightDep := hasDep(r, human)
	result := int64(0)
	var eval *monkey
	if leftDep && rightDep {
		panic("human everywhere!")
	}
	if leftDep {
		eval = l
		result = toNumber(r)
	} else if rightDep {
		eval = r
		result = toNumber(l)
	} else {
		panic("no human?")
	}
	return reverseEval(eval, result, human)
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
