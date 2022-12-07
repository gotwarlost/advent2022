package dec7

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type fileEntry struct {
	Name string
	Size int
}

type dirEntry struct {
	name    string
	parent  *dirEntry
	SubDirs map[string]*dirEntry
	Files   []fileEntry
}

func (s *dirEntry) size() int {
	size := 0
	for _, f := range s.Files {
		size += f.Size
	}
	for _, d := range s.SubDirs {
		size += d.size()
	}
	return size
}

func (s *dirEntry) walk(fn func(name string, d *dirEntry)) {
	fn(s.name, s)
	for _, d := range s.SubDirs {
		d.walk(fn)
	}
}

func readTree(s string) *dirEntry {
	root := &dirEntry{
		name:    "/",
		SubDirs: map[string]*dirEntry{},
	}
	var currentDir *dirEntry
	lines := strings.Split(strings.TrimSpace(s), "\n")
	inList := false
	for _, line := range lines {
		if inList && !strings.HasPrefix(line, "$ ") {
			args := strings.SplitN(line, " ", 2)
			if len(args) != 2 {
				panic(line)
			}
			switch args[0] {
			case "dir":
				subdir := args[1]
				currentDir.SubDirs[subdir] = &dirEntry{
					name:    subdir,
					parent:  currentDir,
					SubDirs: map[string]*dirEntry{},
					Files:   nil,
				}
			default:
				size, err := strconv.Atoi(args[0])
				if err != nil {
					panic(err)
				}
				currentDir.Files = append(currentDir.Files, fileEntry{
					Name: args[1],
					Size: size,
				})
			}
			continue
		}
		inList = false
		command := line[2:]
		args := strings.SplitN(command, " ", 2)
		switch len(args) {
		case 1:
			if command != "ls" {
				panic(line)
			}
			inList = true
		case 2:
			if args[0] != "cd" {
				panic(line)
			}
			dir := args[1]
			switch dir {
			case "/":
				currentDir = root
			case "..":
				currentDir = currentDir.parent
				if currentDir == nil {
					panic("bad dir")
				}
			default:
				if currentDir.SubDirs[dir] == nil {
					panic("wtf")
				}
				currentDir = currentDir.SubDirs[dir]
			}
		}
	}
	return root
}

func run(s string, max int) (p1 int, p2 int) {
	root := readTree(s)
	var size int
	root.walk(func(name string, d *dirEntry) {
		if d.size() < max {
			size += d.size()
		}
	})

	var min int
	unused := 70000000 - root.size()
	required := 30000000
	want := required - unused

	root.walk(func(name string, d *dirEntry) {
		sz := d.size()
		if sz >= want {
			if min == 0 {
				min = sz
			} else {
				if sz < min {
					min = sz
				}
			}
		}
	})

	return size, min
}

func Run() {
	p1, p2 := run(input, 100000)
	fmt.Println("P1:", p1, "P2:", p2)
}
