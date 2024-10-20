package postgres

import (
	"database/sql"
	"fmt"
	"github.com/yoyo-project/yoyo/internal/datatype"
	"strings"

	"github.com/yoyo-project/yoyo/internal/reverse"
	"github.com/yoyo-project/yoyo/internal/schema"
)

const listTablesQuery = `SELECT CONCAT(schemaname, '.', tablename) 
    FROM pg_catalog.pg_tables 
    WHERE schemaname NOT IN (
                             'information_schema',
                             'hstore',
                             'pglogical',
                             'citext',
                             'tiger',
                             'topology'
                            )
		AND schemaname NOT LIKE 'pg_%'`

const listColumnsQuery = `SELECT kcu.column_name
FROM information_schema.key_column_usage kcu
    LEFT JOIN information_schema.table_constraints tc
        ON kcu.constraint_name = tc.constraint_name
            AND kcu.table_schema = tc.table_schema
WHERE kcu.table_schema = '%s'
    AND kcu.table_name = '%s'
    AND tc.constraint_type != 'FOREIGN KEY'`

const listIndicesQuery = `SELECT
    indexname,
    indexdef
FROM
    pg_indexes
WHERE
    schemaname = '%s'
        AND tablename = '%s'`

const listReferencesQuery = `SELECT
    constraint_name
FROM
    information_schema.table_constraints
WHERE
    constraint_type = 'FOREIGN KEY'
  AND table_schema = '%s'
  AND table_name = '%s';`

const getColumnQuery = `SELECT c.column_name, c.data_type, c.column_default, c.is_nullable::boolean, coalesce(c.character_maximum_length, c.numeric_precision), c.numeric_scale, tc.constraint_type = 'PRIMARY KEY'
FROM information_schema.columns c
LEFT JOIN
    information_schema.key_column_usage kcu
ON c.table_schema = kcu.table_schema
    AND c.table_name = kcu.table_name
    AND c.column_name = kcu.column_name
LEFT JOIN
    information_schema.table_constraints tc
    ON kcu.constraint_name = tc.constraint_name
    AND kcu.table_schema = tc.table_schema
    AND kcu.table_name = tc.table_name
WHERE table_schema = '%s'
    AND table_name = '%s'
    AND column_name = '%s'`

const getIndexQuery = `SELECT
    ix.indisunique,
    a.attname,
    a.attnum
FROM pg_index ix
        JOIN pg_class i ON i.oid = ix.indexrelid
        JOIN pg_attribute a ON a.attnum = ANY(ix.indkey) AND a.attrelid = ix.indrelid
WHERE
    i.relname = '%s'
ORDER BY a.attnum ASC`

const getReferenceQuery = `SELECT
    ccu.column_name AS referenced_column,
    rc.update_rule,
    rc.delete_rule
FROM information_schema.table_constraints AS tc
    JOIN information_schema.key_column_usage AS kcu
        ON tc.constraint_name = kcu.constraint_name
            AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
        ON ccu.constraint_name = tc.constraint_name
            AND ccu.table_schema = tc.table_schema
    JOIN information_schema.referential_constraints AS rc
         ON tc.constraint_name = rc.constraint_name
             AND tc.table_schema = rc.constraint_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_name = '%s'
    AND ccu.table_name = '%s'`

// InitReverserBuilder returns a function that returns a PostgreSQL reverse.Adapter
func InitReverserBuilder(open func(driver, dsn string) (*sql.DB, error)) func(dbURL string) (reverse.Adapter, error) {
	return func(dbURL string) (reverse.Adapter, error) {
		reverser := reverser{}

		var err error
		reverser.db, err = open("postgresql", dbURL)
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
	rs, err := r.db.Query(listTablesQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to list tables: %w", err)
	}
	defer rs.Close()

	var (
		tableNames []string
		tempString string
	)

	for rs.Next() {
		err := rs.Scan(&tempString)
		if err != nil {
			return nil, fmt.Errorf("unable to scan table results: %w", err)
		}
		tableNames = append(tableNames, tempString)
	}

	return tableNames, nil
}

// ListColumns returns a []string of column names for the given table
// It does NOT return any columns which are foreign key columns. These will instead come from ListReferences
func (r reverser) ListColumns(table string) ([]string, error) {
	tableParts := schemaAndTable(table)

	rs, err := r.db.Query(fmt.Sprintf(listColumnsQuery, tableParts[0], tableParts[1]))
	if err != nil {
		return nil, fmt.Errorf("unable to list columns: %w", err)
	}
	defer rs.Close()

	var (
		columnNames []string
		tempString  string
	)

	for rs.Next() {
		err := rs.Scan(&tempString)
		if err != nil {
			return nil, fmt.Errorf("unable to scan table results: %w", err)
		}
		columnNames = append(columnNames, tempString)
	}

	return columnNames, nil
}

// ListIndices returns a []string of index names for the given table.
// It will NOT return information referring to PrimaryKey or Foreign Keys, which will instead come from GetColumn and
// ListReferences respectively
func (r reverser) ListIndices(table string) ([]string, error) {
	tableParts := schemaAndTable(table)
	rs, err := r.db.Query(fmt.Sprintf(listIndicesQuery, tableParts[0], tableParts[1]))
	if err != nil {
		return nil, fmt.Errorf("unable to list indices: %w", err)
	}
	defer rs.Close()

	var (
		indexNames []string
		tempString string
	)

	for rs.Next() {
		err := rs.Scan(&tempString)
		if err != nil {
			return nil, fmt.Errorf("unable to scan table results: %w", err)
		}
		indexNames = append(indexNames, tempString)
	}

	return indexNames, nil
}

// ListReferences returns a []string of constraint names for references from the given table.
func (r reverser) ListReferences(table string) ([]string, error) {
	tableParts := schemaAndTable(table)
	rs, err := r.db.Query(fmt.Sprintf(listReferencesQuery, tableParts[0], tableParts[1]))
	if err != nil {
		return nil, fmt.Errorf("unable to list indices: %w", err)
	}
	defer rs.Close()

	var (
		constraintNames []string
		tempString      string
	)

	for rs.Next() {
		err := rs.Scan(&tempString)
		if err != nil {
			return nil, fmt.Errorf("unable to scan table results: %w", err)
		}
		constraintNames = append(constraintNames, tempString)
	}

	return constraintNames, nil
}

// GetColumn returns a schema.Column representing the given tableName and colName.
func (r reverser) GetColumn(table, column string) (schema.Column, error) {
	col := schema.Column{}
	tableParts := schemaAndTable(table)
	rs, err := r.db.Query(fmt.Sprintf(getColumnQuery, tableParts[0], tableParts[1], column))
	if err != nil {
		return col, fmt.Errorf("unable to get column: %w", err)
	}
	defer rs.Close()

	if !rs.Next() {
		return col, fmt.Errorf("column %s not found", column)
	}

	var dt, param1, param2 string

	err = rs.Scan(&col.Name, &dt, &col.Default, &col.Nullable, &param1, &param2, &col.PrimaryKey)
	if err != nil {
		return col, fmt.Errorf("unable scan coolumn: %w", err)
	}

	col.Datatype, err = datatype.FromString(dt)
	if err != nil {
		return col, fmt.Errorf("unable to parse column type: %w", err)
	}

	col.Params = []string{param1, param2}

	//TODO: Autoincrememt/serial?

	return col, nil
}

// GetIndex returns a schema.Index representing the given tableName and indexName.
func (r reverser) GetIndex(_, indexName string) (schema.Index, error) {
	rs, err := r.db.Query(fmt.Sprintf(getIndexQuery, indexName))
	if err != nil {
		return schema.Index{}, fmt.Errorf("unable to query for index: %w", err)
	}
	defer rs.Close()

	ind := schema.Index{Name: indexName}

	var colName string
	var i int
	for rs.Next() {
		err = rs.Scan(&ind.Unique, &colName, &i)
		if err != nil {
			return ind, fmt.Errorf("unable to scan index row: %w", err)
		}
		ind.Columns = append(ind.Columns, colName)
	}

	return ind, nil
}

// GetReference returns a schema.Reference representing the given tableName and columnName.
// TODO: Confirm that order for multi-column keys is correct
// TODO: hasOne vs hasMany?
func (r reverser) GetReference(table, referencedTable string) (schema.Reference, error) {
	rs, err := r.db.Query(fmt.Sprintf(getReferenceQuery, table, referencedTable))
	if err != nil {
		return schema.Reference{}, fmt.Errorf("unable to query for reference: %w", err)
	}
	defer rs.Close()

	ref := schema.Reference{TableName: referencedTable}

	colName := ""

	for rs.Next() {
		err = rs.Scan(
			&colName, &ref.OnUpdate, &ref.OnDelete,
		)
		if err != nil {
			return ref, fmt.Errorf("unable to scan reference row: %w", err)
		}
		ref.ColumnNames = append(ref.ColumnNames, colName)
	}

	return ref, nil
}

func schemaAndTable(in string) [2]string {
	tableParts := strings.Split(in, ".")
	if len(tableParts) == 1 {
		tableParts = []string{"public", in}
	}

	return [2]string{tableParts[0], tableParts[1]}
}
