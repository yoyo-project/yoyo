package main

import (
	"fmt"
	"os"

	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/cli"
	"github.com/dotvezz/lime/options"
	"github.com/dotvezz/yoyo/cmd/yoyo/generate"
	"github.com/dotvezz/yoyo/cmd/yoyo/usecases"
)

func main() {

	ucs := usecases.Init()

	c := cli.New()
	_ = c.SetOptions(options.NoShell)
	_ = c.SetCommands(
		lime.Command{
			Keyword: "generate",
			Commands: []lime.Command{
				{
					Keyword: "migration",
					Func:    generate.BuildMigrationGeneratorFunc(ucs.GetCurrentTime, ucs.LoadMigrationGenerator, os.Create),
				},
			},
		},
		lime.Command{
			Keyword: "reverse",
			Func:    newReverser(ucs.ReadDatabase),
		},
	)
	err := c.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
