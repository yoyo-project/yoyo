package reverse

import (
	"fmt"

	"github.com/dotvezz/yoyo/internal/schema"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

// ReadDatabase uses the given Adapter to scan and write the database into a schema.Database
func InitDatabaseReader(loadAdapter AdapterLoader) DatabaseReader {
	return func(config yoyo.Config) (db schema.Database, err error) {
		var adapter Adapter
		adapter, err = loadAdapter(config.Schema.Dialect)
		if err != nil {
			return db, err
		}
		var tables, columns, indices, references []string
		if tables, err = adapter.ListTables(); err == nil {
			for _, tableName := range tables {
				table := schema.Table{}
				if columns, err = adapter.ListColumns(tableName); err == nil {
					for _, colName := range columns {
						if column, err := adapter.GetColumn(tableName, colName); err == nil {
							if table.Columns == nil {
								table.Columns = map[string]schema.Column{}
							}
							column.SetName(colName)
							table.Columns[colName] = column
						} else {
							return db, fmt.Errorf("%w in GetColumn", err)
						}
					}
				} else {
					return db, fmt.Errorf("%w in ListColumns", err)
				}
				if indices, err = adapter.ListIndices(tableName); err == nil {
					for _, indexName := range indices {
						if index, err := adapter.GetIndex(tableName, indexName); err == nil {
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
				if references, err = adapter.ListReferences(tableName); err == nil {
					for _, refName := range references {
						if reference, err := adapter.GetReference(tableName, refName); err == nil {
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
				table.SetName(tableName)
				db.Tables[tableName] = table
			}
		} else {
			err = fmt.Errorf("%w in ListTables", err)
		}

		return db, err
	}
}
