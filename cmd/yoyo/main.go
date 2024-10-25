package main

import (
	"fmt"
	"os"

	"github.com/yoyo-project/yoyo/internal/file"
	"github.com/yoyo-project/yoyo/internal/validation"

	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/cli"
	"github.com/dotvezz/lime/options"

	"github.com/yoyo-project/yoyo/cmd/yoyo/generate"
	"github.com/yoyo-project/yoyo/cmd/yoyo/usecases"
)

func main() {

	ucs := usecases.Init()

	c := cli.New()
	_ = c.SetOptions(options.NoInteractiveMode)
	_ = c.SetCommands(
		lime.Command{
			Keyword: "generate",
			Commands: []lime.Command{
				{
					Keyword: "migration",
					Func:    generate.Migrations(ucs.GetCurrentTime, ucs.LoadMigrationGenerator, file.CreateWithDirs, validation.ValidateDatabase),
				},
				{
					Keyword: "repos",
					Func:    generate.Repos(ucs.LoadRepositoryGenerator),
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
