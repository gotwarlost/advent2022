package main

import (
	"fmt"
	"log"

	"cd.splunkdev.com/kanantheswaran/advent2022/dec1"
	"github.com/spf13/cobra"
)

type puzzle func(testMode bool)

var puzzles = map[string]puzzle{
	"dec1": dec1.Run,
}

func main() {
	var testMode bool
	cmd := cobra.Command{
		Use: "advent2022",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Usage()
			}
			cmd.SilenceUsage = true
			p := puzzles[args[0]]
			if p == nil {
				return fmt.Errorf("no puzzle named $q available")
			}
			p(testMode)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&testMode, "test", "t", false, "run with test inputs")
	if err := cmd.Execute(); err != nil {
		log.Fatalln("FATAL:", err)
	}
}
