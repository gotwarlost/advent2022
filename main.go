package main

import (
	"log"
	"os"

	"cd.splunkdev.com/kanantheswaran/advent2022/dec1"
	"cd.splunkdev.com/kanantheswaran/advent2022/dec2"
	"cd.splunkdev.com/kanantheswaran/advent2022/dec3"
)

type puzzle func()

var puzzles = map[string]puzzle{
	"dec1":    dec1.Run,
	"dec2-p1": dec2.RunPart1,
	"dec2-p2": dec2.RunPart2,
	"dec3-p1": dec3.RunPart1,
	"dec3-p2": dec3.RunPart2,
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
