package dec25

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

var snafu2Dec = map[byte]int64{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

var rem2Snafu = map[int64]byte{
	0: '0',
	1: '1',
	2: '2',
	3: '=',
	4: '-',
}

func snafuToDecimal(s string) int64 {
	sum := int64(0)
	mult := int64(1)
	for i := len(s) - 1; i >= 0; i-- {
		val := snafu2Dec[s[i]]
		sum += val * mult
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
		val := rem2Snafu[r]
		prepend(val)
		if val == '-' || val == '=' {
			q++
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

func RunP1() {
	fmt.Println(runP1(input))
}
