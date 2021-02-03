package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotvezz/lime"
	"github.com/dotvezz/yoyo/internal/migration"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

type FileOpener func(string) (*os.File, error)

func Migrations(
	now func() time.Time,
	loadGenerator migration.GeneratorLoader,
	create FileOpener,
) lime.Func {
	return func(args []string) error {
		config, err := yoyo.LoadConfig()
		if err != nil {
			return fmt.Errorf("unable to load config: %w", err)
		}

		generate, err := loadGenerator(config)
		if err != nil {
			return fmt.Errorf("unable to initialize migration generator: %w", err)
		}

		sb := strings.Builder{}

		err = generate(config.Schema, &sb)

		if err != nil {
			return fmt.Errorf("unable to generate migration: %w", err)
		}

		var name string
		if len(args) > 0 {
			name = fmt.Sprintf("_%s", strings.ToLower(strings.Join(args, "-")))
		}

		var f *os.File
		defer func() { _ = f.Close() }()
		f, err = create(filepath.Join(config.Paths.Migrations, "%s%s.sql", now().Format("20060102150405"), name))
		if err != nil {
			return fmt.Errorf("cannot create migration file '%s': %w", config.Paths.Migrations, err)
		}

		_, err = f.WriteString(sb.String())
		if err != nil {
			return fmt.Errorf("cannot write to migration file '%s': %w", config.Paths.Migrations, err)
		}

		return nil
	}
}
