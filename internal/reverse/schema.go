package reverse

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/schema"
	"github.com/yoyo-project/yoyo/internal/yoyo"
)

// InitDatabaseReader uses the given AdapterLoader to scan and write the database into a schema.Database
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
							column.Name = colName
							table.Columns = append(table.Columns, column)
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
							index.Name = indexName
							table.Indices = append(table.Indices, index)
						} else {
							return db, fmt.Errorf("%w in GetIndex", err)
						}
					}
				} else {
					return db, fmt.Errorf("%w in ListIndices", err)
				}
				if references, err = adapter.ListReferences(tableName); err == nil {
					for _, ftName := range references {
						if reference, err := adapter.GetReference(tableName, ftName); err == nil {
							reference.TableName = ftName
							table.References = append(table.References, reference)
						} else {
							return db, fmt.Errorf("%w in GetReference", err)
						}
					}
				} else {
					return db, fmt.Errorf("%w in ListReferences", err)
				}

				table.Name = tableName
				db.Tables = append(db.Tables, table)
			}
		} else {
			err = fmt.Errorf("%w in ListTables", err)
		}

		return db, err
	}
}
