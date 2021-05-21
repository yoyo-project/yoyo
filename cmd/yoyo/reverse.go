package main

import (
	"fmt"
	"io"

	"github.com/dotvezz/lime"
	"github.com/yoyo-project/yoyo/internal/reverse"
	"github.com/yoyo-project/yoyo/internal/yoyo"
)

func newReverser(readDatabase reverse.DatabaseReader) lime.Func {
	return func(args []string, w io.Writer) (err error) {
		var config yoyo.Config

		config.Schema, err = readDatabase(config)
		if err != nil {
			return fmt.Errorf("unable to reverse-engineer schema: %w", err)
		}

		return nil
	}
}
