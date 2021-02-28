package postgres

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/dbms/base"
	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
	"github.com/yoyo-project/yoyo/internal/schema"
)

// NewAdapter returns an implementation of migration.Dialect for PostgreSQL
func NewAdapter() *adapter {
	return &adapter{
		Base: base.Base{
			Dialect: dialect.PostgreSQL,
		},
		validator: validator{},
	}
}

type adapter struct {
	base.Base
	validator
}

// TypeString returns the string representation of a given datatype.Datatype for PostgreSQL
// An error will be returned if the datatype.Datatype is invalid or not supported by PostgreSQL
func (a *adapter) TypeString(dt datatype.Datatype) (s string, err error) {
	if !a.SupportsDatatype(dt) {
		return "", fmt.Errorf("datatype %s is not supported in postgresql", dt)
	}
	switch dt {
	default:
		s, err = a.Base.TypeString(dt)
	}
	return s, err
}

func (a *adapter) PreparedStatementPlaceholders(count int) []string {
	panic("implement me")
}

// CreateTable generates a query to create a given table.
func (a *adapter) CreateTable(table string, t schema.Table) string {
	panic("implement me")
}

// AddColumn generates a query that adds a column to an existing table
func (a *adapter) AddColumn(table, column string, c schema.Column) string {
	panic("implement me")
}

// AddIndex returns a string query which adds the specified index to a table
func (a *adapter) AddIndex(table, index string, i schema.Index) string {
	panic("implement me")
}

// AddReference generates a query that adds columns and foreign keys for the given table, foreign table, and schema.Reference
func (a *adapter) AddReference(table string, dt schema.Table, i schema.Reference) string {
	panic("implement me")
}

// ListTables returns a list of tables on the selected database.
func (a *adapter) ListTables() ([]string, error) {
	panic("implement me")
}

// ListColumns returns a []string of column names for the given table
// It does NOT return any columns which are foreign key columns. These will instead come from ListReferences
func (a *adapter) ListColumns(table string) ([]string, error) {
	panic("implement me")
}

// ListIndices returns a []string of index names for the given table.
// It will NOT return information referring to PrimaryKey or Foreign Keys, which will instead come from GetColumn and
// ListReferences respectively
func (a *adapter) ListIndices(table string) ([]string, error) {
	panic("implement me")
}

// ListReferences returns a []string of tables referenced from the given table.
func (a *adapter) ListReferences(table string) ([]string, error) {
	panic("implement me")
}

// GetColumn returns a schema.Column representing the given tableName and colName.
func (a *adapter) GetColumn(table, column string) (schema.Column, error) {
	panic("implement me")
}

// GetIndex returns a schema.Index representing the given tableName and indexName.
func (a *adapter) GetIndex(table, column string) (schema.Index, error) {
	panic("implement me")
}

// GetReference returns a schema.Reference representing the given tableName and indexName.
func (a *adapter) GetReference(table, column string) (schema.Reference, error) {
	panic("implement me")
}
