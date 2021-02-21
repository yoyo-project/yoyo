package generate

import (
	"fmt"

	"github.com/dotvezz/lime"
	"github.com/yoyo-project/yoyo/internal/repository"
	"github.com/yoyo-project/yoyo/internal/yoyo"
)

func Repos(
	loadGenerator repository.GeneratorLoader,
) lime.Func {
	return func(args []string) error {
		config, err := yoyo.LoadConfig()
		if err != nil {
			return fmt.Errorf("unable to load config: %w", err)
		}

		generate := loadGenerator(config)

		err = generate(config.Schema, config.Paths.Repositories)
		if err != nil {
			return fmt.Errorf("unable to generate: %w", err)
		}
		return nil
	}
}
