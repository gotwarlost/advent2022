package dec24

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input string

func runP1(in string) int {
	return 9
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
