package migration

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
	"github.com/yoyo-project/yoyo/internal/dbms/mysql"
	"github.com/yoyo-project/yoyo/internal/dbms/postgres"
	"github.com/yoyo-project/yoyo/internal/schema"
)

type AdapterLoader func(name string) (d Adapter, err error)

// Adapter describes an internal dialect to be used by Yoyo for generating Migration code for a given DBMS.
type Adapter interface {
	// CreateTable returns a string query which creates a full table with columns columns and primary key
	CreateTable(table string, t schema.Table) string

	// AddColumn returns a string query which adds the specified column to a table
	AddColumn(table, column string, c schema.Column) string

	// AddIndex returns a string query which adds the specified index to a table
	AddIndex(table, index string, i schema.Index) string

	// AddReference returns a string query which adds the specified index to a table
	AddReference(table string, dt schema.Table, i schema.Reference) string
}

// LoadAdapter loads and returns an implementation of Adapter corresponding to the given name string
func LoadAdapter(name string) (a Adapter, err error) {
	switch name {
	case dialect.MySQL:
		a = mysql.NewAdapter()
	case dialect.PostgreSQL:
		a = postgres.NewAdapter()
	case dialect.SQLite:
		a = postgres.NewAdapter()
	default:
		err = fmt.Errorf("unknown dialect `%s`", name)
	}
	return a, err
}
