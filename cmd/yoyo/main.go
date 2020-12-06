package main

import (
	"database/sql"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/cli"
	"github.com/dotvezz/lime/options"
	"github.com/dotvezz/yoyo/internal/dbms/mysql"
	"github.com/dotvezz/yoyo/internal/dbms/postgres"
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
			Func:    initReverser(mysql.InitNewReverser(sql.Open), postgres.InitNewReverser(sql.Open)),
		},
	)
	err := c.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
