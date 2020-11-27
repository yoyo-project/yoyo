package main

import (
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/yoyo/env"
	"github.com/dotvezz/yoyo/internal/reverse"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

func initReverser() lime.Func {
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

		if err != nil {
			return fmt.Errorf("unable to initialize: %w", err)
		}

		reverser, err = reverse.LoadReverser(dia, env.DBHost(), env.DBUser(), env.DBName(), env.DBPassword(), env.DBPort())
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
