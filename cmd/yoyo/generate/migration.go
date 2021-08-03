package generate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotvezz/lime"

	"github.com/yoyo-project/yoyo/internal/migration"
	"github.com/yoyo-project/yoyo/internal/schema"
	"github.com/yoyo-project/yoyo/internal/yoyo"
)

type FileOpener func(string) (*os.File, error)
type DatabaseValidator func(database schema.Database) error

func Migrations(
	now func() time.Time,
	loadGenerator migration.GeneratorLoader,
	create FileOpener,
	validate DatabaseValidator,
) lime.Func {
	return func(args []string, w io.Writer) error {
		config, err := yoyo.LoadConfig()
		if err != nil {
			return fmt.Errorf("unable to load config: %w", err)
		}

		if err := validate(config.Schema); err != nil {
			return err
		}

		var generate migration.Generator
		generate, err = loadGenerator(config)
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
		f, err = create(filepath.Join(config.Paths.Migrations, fmt.Sprintf("%s%s.sql", now().Format("20060102150405"), name)))
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
