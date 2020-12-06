package main

import (
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/yoyo/env"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/reverse"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

func initReverser(newMysqlReverser, newPostgresReverser func(host, userName, dbName, password, port string) (reverse.Reverser, error)) lime.Func {
	return func(args []string) error {
		var (
			config   yoyo.Config
			reverser reverse.Reverser
			dia      string
			err      error
		)

		if len(args) > 0 {
			dia = args[0]
		} else {
			dia = config.Schema.Dialect
		}

		if dia == "" {
			return fmt.Errorf("no dialect given")
		}

		switch dia {
		case dialect.MySQL:
			reverser, err = newMysqlReverser(env.DBHost(), env.DBUser(), env.DBName(), env.DBPassword(), env.DBPort())
		case dialect.PostgreSQL:
			reverser, err = newPostgresReverser(env.DBHost(), env.DBUser(), env.DBName(), env.DBPassword(), env.DBPort())
		default:
			err = fmt.Errorf("unknown dialect `%s`", dia)
		}

		if err != nil {
			return fmt.Errorf("unable to initialize: %w", err)
		}

		config.Schema, err = reverse.ReadSchema(reverser)
		if err != nil {
			return fmt.Errorf("unable to reverse-engineer schema: %w", err)
		}

		return nil
	}
}
