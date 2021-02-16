package repository

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/yoyo-project/yoyo/internal/yoyo"

	"github.com/yoyo-project/yoyo/internal/schema"
)

type Generator func(db schema.Database, repositoriesPath string) error
type GeneratorLoader func(config yoyo.Config) Generator
type WriteGenerator func(db schema.Database, w io.StringWriter) error
type SimpleWriteGenerator func(w io.StringWriter) error
type EntityGenerator func(t schema.Table, w io.StringWriter) error
type FileOpener func(string) (*os.File, error)

func NewGenerator(
	generateEntity EntityGenerator,
	generateRepository EntityGenerator,
	generateQueryFile EntityGenerator,
	generateRepositoriesFile WriteGenerator,
	generateQueryNodeFile SimpleWriteGenerator,
	create FileOpener,
) Generator {
	return func(db schema.Database, repositoriesPath string) error {
		for _, t := range db.Tables {
			err := func() error {
				fName := filepath.Join(repositoriesPath, fmt.Sprintf("entity_%s.go", t.QueryPackageName()))
				f, err := create(fName)
				defer func() {
					if f != nil {
						_ = f.Close()
					}
				}()
				if err != nil {
					return fmt.Errorf("unable to create entity file %s for %s: %w", fName, t.QueryPackageName(), err)
				}

				err = generateEntity(t, f)
				if err != nil {
					return fmt.Errorf("unable to write to entity file %s for %s: %w", fName, t.QueryPackageName(), err)
				}
				return nil
			}()
			if err != nil {
				return err
			}

			err = func() error {
				fName := filepath.Join(repositoriesPath, fmt.Sprintf("repository_%s.go", t.QueryPackageName()))
				f, err := create(fName)
				defer func() {
					if f != nil {
						_ = f.Close()
					}
				}()
				if err != nil {
					return fmt.Errorf("unable to create repository file %s for %s: %w", fName, t.QueryPackageName(), err)
				}

				err = generateRepository(t, f)
				if err != nil {
					return fmt.Errorf("unable to write to repository file %s for %s: %w", fName, t.QueryPackageName(), err)
				}
				return nil
			}()
			if err != nil {
				return err
			}

			err = func() error {
				fName := filepath.Join(repositoriesPath, "query", t.QueryPackageName(), "query.go")
				f, err := create(fName)
				defer func() {
					if f != nil {
						_ = f.Close()
					}
				}()
				if err != nil {
					return fmt.Errorf("unable to create query file %s for %s: %w", fName, t.QueryPackageName(), err)
				}

				err = generateQueryFile(t, f)
				if err != nil {
					return fmt.Errorf("unable to write to query file %s for %s: %w", fName, t.QueryPackageName(), err)
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}

		err := func() error {
			fName := filepath.Join(repositoriesPath, "repositories.go")
			f, err := create(fName)
			defer func() {
				if f != nil {
					_ = f.Close()
				}
			}()
			defer func() { _ = f.Close() }()
			if err != nil {
				return fmt.Errorf("unable to create query file %s: %w", fName, err)
			}

			err = generateRepositoriesFile(db, f)
			if err != nil {
				return fmt.Errorf("unable to write to query file %s %w", fName, err)
			}
			return nil
		}()
		if err != nil {
			return err
		}

		err = func() error {
			fName := filepath.Join(repositoriesPath, "/query/node.go")
			f, err := create(fName)
			defer func() {
				if f != nil {
					_ = f.Close()
				}
			}()
			defer func() { _ = f.Close() }()
			if err != nil {
				return fmt.Errorf("unable to create query file %s: %w", fName, err)
			}

			err = generateQueryNodeFile(f)
			if err != nil {
				return fmt.Errorf("unable to write to query file %s: %w", fName, err)
			}
			return nil
		}()
		if err != nil {
			return err
		}
		return nil
	}
}

func InitGeneratorLoader(
	newGenerator func(EntityGenerator, EntityGenerator, EntityGenerator, WriteGenerator, SimpleWriteGenerator, FileOpener) Generator,
	loadAdapter AdapterLoader,
) GeneratorLoader {
	return func(config yoyo.Config) Generator {
		adapter, _ := loadAdapter(config.Schema.Dialect)
		return newGenerator(
			NewEntityGenerator(config.Schema.Tables),
			NewEntityRepositoryGenerator(config, adapter),
			NewQueryFileGenerator(config),
			NewRepositoriesGenerator(),
			NewQueryNodeGenerator(),
			os.Create,
		)
	}
}
