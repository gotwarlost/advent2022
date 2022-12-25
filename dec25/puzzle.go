package dec25

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

func snafuToDecimal(s string) int64 {
	sum := int64(0)
	mult := int64(1)
	for i := len(s) - 1; i >= 0; i-- {
		x := s[i]
		var val int
		switch x {
		case '=':
			val = -2
		case '-':
			val = -1
		case '0':
			val = 0
		case '1':
			val = 1
		case '2':
			val = 2
		default:
			panic(fmt.Errorf("invalid SNAFU digit: %v", x))
		}
		sum += int64(val) * mult
		mult *= 5
	}
	return sum
}

func decimalToSnafu(x int64) string {
	var ret []byte
	prepend := func(b byte) {
		ret = append([]byte{b}, ret...)
	}
	for {
		q := x / int64(5)
		r := x % int64(5)
		switch r {
		case 0:
			prepend('0')
		case 1:
			prepend('1')
		case 2:
			prepend('2')
		case 3:
			q++
			prepend('=')
		case 4:
			q++
			prepend('-')
		}
		x = q
		if x == 0 {
			break
		}
	}
	return string(ret)
}

func toSnafus(in string) []string {
	return strings.Split(strings.TrimSpace(in), "\n")
}

func runP1(in string) string {
	total := int64(0)
	for _, s := range toSnafus(in) {
		total += snafuToDecimal(s)
	}
	return decimalToSnafu(total)
}

func runP2(in string) int {
	return 9
}

func RunP1() {
	fmt.Println(runP1(input))
}

func RunP2() {
	fmt.Println(runP2(input))
}
