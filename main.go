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
	"github.com/gotwarlost/advent2022/dec13"
	"github.com/gotwarlost/advent2022/dec14"
	"github.com/gotwarlost/advent2022/dec15"
	"github.com/gotwarlost/advent2022/dec16"
	"github.com/gotwarlost/advent2022/dec17"
	"github.com/gotwarlost/advent2022/dec18"
	"github.com/gotwarlost/advent2022/dec19"
	"github.com/gotwarlost/advent2022/dec20"
	"github.com/gotwarlost/advent2022/dec21"
	"github.com/gotwarlost/advent2022/dec22"
	"github.com/gotwarlost/advent2022/dec23"
	"github.com/gotwarlost/advent2022/dec24"
	"github.com/gotwarlost/advent2022/dec25"
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
	"dec13-p1": dec13.RunP1,
	"dec13-p2": dec13.RunP2,
	"dec14-p1": dec14.RunP1,
	"dec14-p2": dec14.RunP2,
	"dec15-p1": dec15.RunP1,
	"dec15-p2": dec15.RunP2,
	"dec16-p1": dec16.RunP1,
	"dec16-p2": dec16.RunP2,
	"dec17-p1": dec17.RunP1,
	"dec17-p2": dec17.RunP2,
	"dec18-p1": dec18.RunP1,
	"dec18-p2": dec18.RunP2,
	"dec19-p1": dec19.RunP1,
	"dec19-p2": dec19.RunP2,
	"dec20-p1": dec20.RunP1,
	"dec20-p2": dec20.RunP2,
	"dec21-p1": dec21.RunP1,
	"dec21-p2": dec21.RunP2,
	"dec22-p1": dec22.RunP1,
	"dec22-p2": dec22.RunP2,
	"dec23-p1": dec23.RunP1,
	"dec23-p2": dec23.RunP2,
	"dec24-p1": dec24.RunP1,
	"dec24-p2": dec24.RunP2,
	"dec25-p1": dec25.RunP1,
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
