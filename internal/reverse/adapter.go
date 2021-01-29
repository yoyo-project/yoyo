package reverse

import (
	"fmt"

	"github.com/dotvezz/yoyo/env"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/schema"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

type AdapterLoader func(dia string) (adapter Adapter, err error)
type DatabaseReader func(config yoyo.Config) (db schema.Database, err error)
type AdapterBuilder func(host, userName, dbName, password, port string) (Adapter, error)

// Adapter is the yoyo interface for reverse-engineering databases to a schema.Database for creating diff migrations
// and for creating a yoyo.yml from an existing database.
type Adapter interface {
	ListTables() ([]string, error)

	// ListColumns returns a []string of column names for the given table
	// It MUST NOT return any columns which are foreign key columns. These will instead come from ListReferences
	ListColumns(table string) ([]string, error)

	// ListIndices returns a []string of index names for the given table.
	// It MUST NOT return information referring to PrimaryKey or Foreign Keys, which will instead come from GetColumn and
	// ListReferences respectively
	ListIndices(table string) ([]string, error)

	// ListReferences returns a []string of tables referenced from the given table.
	ListReferences(table string) ([]string, error)

	// GetColumn returns a schema.Column representing the given tableName and colName.
	GetColumn(table, column string) (schema.Column, error)

	// GetIndex returns a schema.Index representing the given tableName and indexName.
	GetIndex(table, column string) (schema.Index, error)

	// GetReference returns a schema.Reference representing the given tableName and indexName.
	GetReference(table, column string) (schema.Reference, error)
}

func InitAdapterSelector(newMysqlReverser, newPostgresReverser AdapterBuilder) func(dia string) (adapter Adapter, err error) {
	return func(dia string) (adapter Adapter, err error) {
		switch dia {
		case dialect.MySQL:
			adapter, err = newMysqlReverser(env.DBHost(), env.DBUser(), env.DBName(), env.DBPassword(), env.DBPort())
		case dialect.PostgreSQL:
			adapter, err = newPostgresReverser(env.DBHost(), env.DBUser(), env.DBName(), env.DBPassword(), env.DBPort())
		default:
			err = fmt.Errorf("unknown dialect `%s`", dia)
		}

		return adapter, err
	}
}
