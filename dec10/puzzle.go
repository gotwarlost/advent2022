package dec10

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input string

func run(in string) int {
	return 0
}

func RunP1() {
	fmt.Println(run(input))
}

func RunP2() {
	fmt.Println(run(input))
}
