package reverse

import (
	"github.com/dotvezz/yoyo/internal/schema"
)

type Reverser interface {
	ListTables() ([]string, error)

	// ListColumns returns a []string of column names for the currently-connected database
	// It MUST NOT return any columns which are foreign key columns. These will instead be come from ListReferences
	ListColumns(table string) ([]string, error)

	// ListIndices returns a []string of index names for the given table.
	// It MUST NOT return information referring to PrimaryKey or Foreign Keys, which will instead come from GetColumn and
	// ListReferences respectively
	ListIndices(table string) ([]string, error)

	ListReferences(table string) ([]string, error)

	GetColumn(table, column string) (schema.Column, error)

	GetIndex(table, column string) (schema.Index, error)

	GetReference(table, column string) (schema.Reference, error)
}
