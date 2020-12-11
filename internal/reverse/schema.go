package reverse

import (
	"fmt"
	"github.com/dotvezz/yoyo/internal/schema"
)

// ReadDatabase uses the given Reverser to scan and write the database into a schema.Database
func ReadDatabase(reverser Reverser) (db schema.Database, err error) {
	var tables, columns, indices, references []string
	if tables, err = reverser.ListTables(); err == nil {
		for _, tableName := range tables {
			table := schema.Table{}
			if columns, err = reverser.ListColumns(tableName); err == nil {
				for _, colName := range columns {
					if column, err := reverser.GetColumn(tableName, colName); err == nil {
						if table.Columns == nil {
							table.Columns = map[string]schema.Column{}
						}
						table.Columns[colName] = column
					} else {
						return db, fmt.Errorf("%w in GetColumn", err)
					}
				}
			} else {
				return db, fmt.Errorf("%w in ListColumns", err)
			}
			if indices, err = reverser.ListIndices(tableName); err == nil {
				for _, indexName := range indices {
					if index, err := reverser.GetIndex(tableName, indexName); err == nil {
						if table.Indices == nil {
							table.Indices = map[string]schema.Index{}
						}
						table.Indices[indexName] = index
					} else {
						return db, fmt.Errorf("%w in GetIndex", err)
					}
				}
			} else {
				return db, fmt.Errorf("%w in ListIndices", err)
			}
			if references, err = reverser.ListReferences(tableName); err == nil {
				for _, refName := range references {
					if reference, err := reverser.GetReference(tableName, refName); err == nil {
						if table.References == nil {
							table.References = map[string]schema.Reference{}
						}
						table.References[refName] = reference
					} else {
						return db, fmt.Errorf("%w in GetReference", err)
					}
				}
			} else {
				return db, fmt.Errorf("%w in ListReferences", err)
			}
			if db.Tables == nil {
				db.Tables = map[string]schema.Table{}
			}
			db.Tables[tableName] = table
		}
	} else {
		err = fmt.Errorf("%w in ListTables", err)
	}

	return db, err
}
