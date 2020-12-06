package migration

import (
	"fmt"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/dbms/mysql"
	"github.com/dotvezz/yoyo/internal/dbms/postgres"
	"github.com/dotvezz/yoyo/internal/schema"
)

// Dialect describes an internal dialect to be used by Yoyo for generating Migration code for a given DBMS.
type Dialect interface {
	// CreateTable returns a string query which creates a full table with columns columns and primary key
	CreateTable(table string, t schema.Table) string

	// AddColumn returns a string query which adds the specified column to a table
	AddColumn(table, column string, c schema.Column) string

	// AddIndex returns a string query which adds the specified index to a table
	AddIndex(table, index string, i schema.Index) string

	// AddReference returns a string query which adds the specified index to a table
	AddReference(table, referencedTable string, dt schema.Table, i schema.Reference) (string, error)
}

func LoadDialect(name string) (d Dialect, err error) {
	switch name {
	case dialect.MySQL:
		d = mysql.NewMigrator()
	case dialect.PostgreSQL:
		d = postgres.NewMigrator()
	default:
		err = fmt.Errorf("unknown dialect `%s`", name)
	}
	return d, err
}
