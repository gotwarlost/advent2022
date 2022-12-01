package main

import (
	"log"
	"os"

	"cd.splunkdev.com/kanantheswaran/advent2022/dec1"
)

type puzzle func()

var puzzles = map[string]puzzle{
	"dec1": dec1.Run,
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalln("usage: advent2022 <puzzle-name>")
	}
	name := args[1]
	p := puzzles[name]
	if p == nil {
		log.Fatalf("no puzzle named %q available", name)
	}
	p()
}
