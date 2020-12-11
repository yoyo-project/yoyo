package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dotvezz/yoyo/internal/reverse"
	"github.com/dotvezz/yoyo/internal/schema"
)

// InitNewReverser returns a function that returns a PostgreSQL reverse.Reverser
func InitNewReverser(open func(driver, dsn string) (*sql.DB, error)) func(host, user, dbname, password, port string) (reverse.Reverser, error) {
	return func(host, user, dbname, password, port string) (reverse.Reverser, error) {
		reverser := reverser{}

		var err error
		reverser.db, err = open("postgresql", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
		if err != nil {
			return nil, fmt.Errorf("unable to open database connection for mysql reverser: %w", err)
		}

		return &reverser, nil
	}
}

type reverser struct {
	db *sql.DB
}

// ListTables returns a list of tables on the selected database.
func (r reverser) ListTables() ([]string, error) {
	panic("implement me")
}

// ListColumns returns a []string of column names for the given table
// It does NOT return any columns which are foreign key columns. These will instead come from ListReferences
func (r reverser) ListColumns(table string) ([]string, error) {
	panic("implement me")
}

// ListIndices returns a []string of index names for the given table.
// It will NOT return information referring to PrimaryKey or Foreign Keys, which will instead come from GetColumn and
// ListReferences respectively
func (r reverser) ListIndices(table string) ([]string, error) {
	panic("implement me")
}

// ListReferences returns a []string of tables referenced from the given table.
func (r reverser) ListReferences(table string) ([]string, error) {
	panic("implement me")
}

// GetColumn returns a schema.Column representing the given tableName and colName.
func (r reverser) GetColumn(table, column string) (schema.Column, error) {
	panic("implement me")
}

// GetIndex returns a schema.Index representing the given tableName and indexName.
func (r reverser) GetIndex(table, column string) (schema.Index, error) {
	panic("implement me")
}

// GetReference returns a schema.Reference representing the given tableName and indexName.
func (r reverser) GetReference(table, column string) (schema.Reference, error) {
	panic("implement me")
}
