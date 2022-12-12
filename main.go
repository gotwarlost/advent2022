package main

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gotwarlost/advent2022/dec01"
	"github.com/gotwarlost/advent2022/dec02"
	"github.com/gotwarlost/advent2022/dec03"
	"github.com/gotwarlost/advent2022/dec04"
	"github.com/gotwarlost/advent2022/dec05"
	"github.com/gotwarlost/advent2022/dec06"
	"github.com/gotwarlost/advent2022/dec07"
	"github.com/gotwarlost/advent2022/dec08"
	"github.com/gotwarlost/advent2022/dec09"
	"github.com/gotwarlost/advent2022/dec10"
	"github.com/gotwarlost/advent2022/dec11"
	"github.com/gotwarlost/advent2022/dec12"
)

type puzzle func()

var puzzles = map[string]puzzle{
	"dec1":     dec01.Run,
	"dec2-p1":  dec02.RunPart1,
	"dec2-p2":  dec02.RunPart2,
	"dec3-p1":  dec03.RunPart1,
	"dec3-p2":  dec03.RunPart2,
	"dec4-p1":  dec04.RunPart1,
	"dec4-p2":  dec04.RunPart2,
	"dec5-p1":  dec05.RunPart1,
	"dec5-p2":  dec05.RunPart2,
	"dec6-p1":  dec06.RunPart1,
	"dec6-p2":  dec06.RunPart2,
	"dec7":     dec07.Run,
	"dec8-p1":  dec08.RunP1,
	"dec8-p2":  dec08.RunP2,
	"dec9-p1":  dec09.RunP1,
	"dec9-p2":  dec09.RunP2,
	"dec10-p1": dec10.RunP1,
	"dec10-p2": dec10.RunP2,
	"dec11-p1": dec11.RunP1,
	"dec11-p2": dec11.RunP2,
	"dec12-p1": dec12.RunP1,
	"dec12-p2": dec12.RunP2,
}

func names() []string {
	var out []string
	for k := range puzzles {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalln("usage: advent2022 <puzzle-name>")
	}
	name := args[1]
	p := puzzles[name]
	if p == nil {
		log.Fatalf("no puzzle named %q available, available names are: \n\t%s", name, strings.Join(names(), "\n\t"))
	}
	p()
}
