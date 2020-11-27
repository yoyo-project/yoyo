package main

import (
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/cli"
	"github.com/dotvezz/lime/options"
	"os"
)

func main() {
	c := cli.New()
	_ = c.SetOptions(options.NoShell)
	_ = c.SetCommands(
		//lime.Command{
		//	Keyword: "generate",
		//	Commands: []lime.Command{
		//		{
		//			Keyword: "migration",
		//			Func:    generate.MigrationGenerator(time.Now()),
		//		},
		//	},
		//},
		lime.Command{
			Keyword: "reverse",
			Func:    initReverser(),
		},
	)
	err := c.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
