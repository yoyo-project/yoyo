package migration

import (
	"fmt"
	"io"

	"github.com/dotvezz/yoyo/internal/reverse"
	"github.com/dotvezz/yoyo/internal/schema"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

const (
	// AddAll is the option that tells some generator methods to not check if things exist and just add them all
	AddAll uint8 = 0x000
	// AddMissing is the option that tells some generator methods to check if things exist and skip over them if they do
	AddMissing uint8 = 0x001
)

// TableGenerator functions take a string, schema.Table, and io.StringWriter. Implementations will use them to generate
// SQL for creating or modifying tables
type TableGenerator func(name string, table schema.Table, sw io.StringWriter) error

// RefGenerator functions take a string, map[string]schema.Reference, and io.StringWriter. Implementations will use them
// to generate SQL for working with references
type RefGenerator func(localTable string, refs map[string]schema.Reference, sw io.StringWriter) error

// StringSearcher functions take a string and return true if the matching entity (table, column, etc) exists.
type StringSearcher func(string) (bool, error)

// Generator functions take a schema.Database and io.StringWriter, generating stuff for the Database to the io.StringWriter
type Generator func(db schema.Database, w io.StringWriter) error

type GeneratorLoader func(config yoyo.Config) (Generator, error)

// NewGenerator returns a function that generates a schema and writes it to the given io.StringWriter.
func NewGenerator(
	createTable TableGenerator,
	addMissingColumns TableGenerator,
	addMissingIndices TableGenerator,
	addAllIndices TableGenerator,
	hasTable StringSearcher,
	addMissingRefs RefGenerator,
	addAllRefs RefGenerator,
) Generator {
	return func(db schema.Database, w io.StringWriter) error {
		for n, t := range db.Tables {
			exists, err := hasTable(n)
			if err != nil {
				return fmt.Errorf("unable to check if table existss: %w", err)
			}

			if !exists {
				err = createTable(n, t, w)
				if err != nil {
					return fmt.Errorf("unable to unable to generate table create query: %w", err)
				}

				err = addAllIndices(n, t, w)
				if err != nil {
					return fmt.Errorf("unable to generate queries to add existingIndices: %w", err)
				}
			} else {
				err = addMissingColumns(n, t, w)
				if err != nil {
					return fmt.Errorf("unable to generate queries to add existingColumns: %w", err)
				}

				err = addMissingIndices(n, t, w)
				if err != nil {
					return fmt.Errorf("unable to generate queries to add existingIndices: %w", err)
				}
			}
		}

		// Generate References/Foreign Keys after generating all tables
		for n, t := range db.Tables {
			exists, err := hasTable(n)
			if err != nil {
				return fmt.Errorf("unable to check if table existss: %w", err)
			}

			if !exists {
				err = addAllRefs(n, t.References, w)
			} else {
				err = addMissingRefs(n, t.References, w)
			}
		}

		return nil
	}
}

// NewTableAdder returns a TableGenerator that adds a table
func NewTableAdder(
	a Adapter,
) TableGenerator {
	return func(tName string, t schema.Table, sw io.StringWriter) error {
		_, err := sw.WriteString(fmt.Sprintf("%s\n", a.CreateTable(tName, t)))
		if err != nil {
			return fmt.Errorf("unable to generate migration: %sw", err)
		}
		return nil
	}
}

// NewColumnAdder returns a TableGenerator that adds columns from a schema.Table.
func NewColumnAdder(
	a Adapter,
	options uint8,
	hasColumn reverse.TableSearcher,
) TableGenerator {
	return func(tName string, t schema.Table, sw io.StringWriter) error {
		for cName, c := range t.Columns {
			if options&AddMissing > 0 {

				if hasColumn(tName, cName) {
					continue
				}
			}

			_, err := sw.WriteString(fmt.Sprintf("%s\n", a.AddColumn(tName, cName, c)))
			if err != nil {
				return fmt.Errorf("unable to generate migration: %sw", err)
			}
		}
		return nil
	}
}

// NewIndexAdder returns a TableGenerator that adds indices from a schema.Table.
func NewIndexAdder(
	a Adapter,
	options uint8,
	hasIndex reverse.TableSearcher,
) TableGenerator {
	return func(tName string, t schema.Table, sw io.StringWriter) error {
		for iName, i := range t.Indices {
			if options&AddMissing > 0 {
				if hasIndex(tName, iName) {
					continue
				}
			}
			_, err := sw.WriteString(a.AddIndex(tName, iName, i))
			if err != nil {
				return fmt.Errorf("unable to generate migration: %sw", err)
			}
		}
		return nil
	}
}

// NewRefAdder returns a RefGenerator that adds references to a given table.
func NewRefAdder(
	a Adapter,
	db schema.Database,
	options uint8,
	hasReference reverse.TableSearcher,
) RefGenerator {
	return func(localTable string, refs map[string]schema.Reference, sw io.StringWriter) error {
		for foreignTable, ref := range refs {
			if ref.HasMany { // swap the tables if it's a HasMany
				localTable, foreignTable = foreignTable, localTable
			}
			ft, ok := db.Tables[foreignTable]

			if options&AddMissing > 0 {
				if hasReference(localTable, foreignTable) {
					continue
				}
			}

			if !ok { // This should technically be caught by validation, but still
				return fmt.Errorf("referenced table `%s` does not exist in dbms definition", foreignTable)
			}
			s := a.AddReference(localTable, foreignTable, ft, ref)
			_, err := sw.WriteString(s)
			if err != nil {
				return fmt.Errorf("unable to generate migration: %w", err)
			}
		}
		return nil
	}
}

func InitGeneratorLoader(
	initReverseAdapter func(dia string) (adapter reverse.Adapter, err error),
	initMigrationAdapter func(dia string) (a Adapter, err error),
	newGenerator func(TableGenerator, TableGenerator, TableGenerator, TableGenerator, StringSearcher, RefGenerator, RefGenerator) Generator,
) GeneratorLoader {
	return func(config yoyo.Config) (Generator, error) {
		var (
			err      error
			reverser reverse.Adapter
			migrator Adapter
		)
		reverser, err = initReverseAdapter(config.Schema.Dialect)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize reverse adapter: %w", err)
		}
		migrator, err = initMigrationAdapter(config.Schema.Dialect)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize migration adapter: %w", err)
		}

		return newGenerator(
			NewTableAdder(migrator),
			NewColumnAdder(migrator, AddMissing, reverse.InitHasColumn(reverser.GetColumn)),
			NewIndexAdder(migrator, AddMissing, reverse.InitHasIndex(reverser.GetIndex)),
			NewIndexAdder(migrator, AddAll, nil),
			reverse.InitHasTable(reverser.ListTables),
			NewRefAdder(migrator, config.Schema, AddMissing, reverse.InitHasReference(reverser.GetReference)),
			NewRefAdder(migrator, config.Schema, AddAll, nil),
		), nil
	}
}
