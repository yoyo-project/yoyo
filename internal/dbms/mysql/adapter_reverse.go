package mysql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/reverse"
	"github.com/yoyo-project/yoyo/internal/schema"
)

var paramIsolator = regexp.MustCompile("(^.*?(\\(|$)|[\"\\)\\s])")

const listColumnsQuery = `SELECT c.COLUMN_NAME FROM information_schema.COLUMNS c
    LEFT JOIN information_schema.KEY_COLUMN_USAGE kcu
        ON c.COLUMN_NAME = kcu.COLUMN_NAME
            AND c.TABLE_NAME = kcu.TABLE_NAME
            AND c.TABLE_SCHEMA = kcu.TABLE_SCHEMA
    WHERE c.TABLE_NAME = '%s'
        AND c.TABLE_SCHEMA = DATABASE()
        AND kcu.REFERENCED_COLUMN_NAME IS NULL`

const listIndicesQuery = `SELECT INDEX_NAME FROM information_schema.STATISTICS
    LEFT JOIN information_schema.KEY_COLUMN_USAGE ON CONSTRAINT_NAME = INDEX_NAME
    WHERE STATISTICS.TABLE_NAME = '%s' 
        AND STATISTICS.TABLE_SCHEMA = DATABASE()
        AND INDEX_NAME != 'PRIMARY'
        AND REFERENCED_TABLE_NAME IS NULL
    GROUP BY INDEX_NAME`

const listReferencesQuery = `SELECT REFERENCED_TABLE_NAME 
	FROM information_schema.KEY_COLUMN_USAGE 
	WHERE TABLE_NAME = %s 
		AND TABLE_SCHEMA = DATABASE()`

// TODO: Charset and Collation
const getColumnQuery = "SHOW COLUMNS IN `%s` WHERE Field = '%s'"

const getIndexQuery = `SELECT NOT NON_UNIQUE, COLUMN_NAME 
    FROM information_schema.STATISTICS
    WHERE TABLE_NAME = '%s'
        AND TABLE_SCHEMA = DATABASE()
        AND INDEX_NAME = '%s'
    ORDER BY SEQ_IN_INDEX`

const getReferenceQuery = `SELECT UPDATE_RULE, DELETE_RULE, CONSTRAINT_NAME
    FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE TABLE_NAME = '%s'
		AND CONSTRAINT_SCHEMA = DATABASE()
        AND REFERENCED_TABLE_NAME = '%s'`

const getReferenceColumnsQuery = `SELECT kcu.COLUMN_NAME, NOT c.IS_NULLABLE
    FROM information_schema.KEY_COLUMN_USAGE kcu
    LEFT JOIN information_schema.COLUMNS c
        ON c.COLUMN_NAME = kcu.COLUMN_NAME
            AND c.TABLE_NAME = kcu.TABLE_NAME
            AND c.TABLE_SCHEMA = kcu.TABLE_SCHEMA
    WHERE kcu.TABLE_NAME = '%s'
		AND kcu.TABLE_SCHEMA = DATABASE()
        AND kcu.REFERENCED_TABLE_NAME = '%s'
        AND kcu.CONSTRAINT_NAME = '%s'`

// InitReverserBuilder returns a `NewReverser` function, which returns a reverse.Adapter.
func InitReverserBuilder(open func(driver, dsn string) (*sql.DB, error)) func(dbURL string) (reverse.Adapter, error) {
	return func(dbURL string) (reverse.Adapter, error) {
		r := adapter{}

		var err error
		r.db, err = open("mysql", dbURL)
		if err != nil {
			return nil, fmt.Errorf("unable to open database connection for mysql r: %w", err)
		}

		return &r, nil
	}
}

// ListTables returns a list of tables on the selected database.
func (a *adapter) ListTables() ([]string, error) {
	rs, err := a.db.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("unable to list tables: %w", err)
	}

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
	_ = rs.Close()

	return tableNames, nil
}

// ListIndices returns a []string of index names for the given table.
// It will NOT return information referring to PrimaryKey or Foreign Keys, which will instead come from GetColumn and
// ListReferences respectively
func (a *adapter) ListIndices(table string) ([]string, error) {
	query := fmt.Sprintf(listIndicesQuery, table)
	rs, err := a.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to list indices: %w", err)
	}

	var (
		tempString string
		indexNames = make([]string, 0)
	)

	for rs.Next() {
		err := rs.Scan(&tempString)
		if err != nil {
			return nil, fmt.Errorf("unable to scan index list results: %w", err)
		}
		indexNames = append(indexNames, tempString)
	}

	return indexNames, nil
}

// ListColumns returns a []string of column names for the given table
// It does NOT return any columns which are foreign key columns. These will instead come from ListReferences
func (a *adapter) ListColumns(table string) ([]string, error) {
	rs, err := a.db.Query(fmt.Sprintf(listColumnsQuery, table))
	if err != nil {
		return nil, fmt.Errorf("unable to list columns: %w", err)
	}

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
	_ = rs.Close()

	return columnNames, nil
}

// ListReferences returns a []string of tables referenced from the given table.
func (a *adapter) ListReferences(table string) ([]string, error) {
	rs, err := a.db.Query(fmt.Sprintf(listReferencesQuery, table))
	if err != nil {
		return nil, fmt.Errorf("unable to list indices: %w", err)
	}

	tableNames := make([]string, 0)

	for rs.Next() {
		scanDestination := new(string)
		err := rs.Scan(&scanDestination)
		if err != nil {
			return nil, fmt.Errorf("unable to scan table results: %w", err)
		}
		if scanDestination != nil {
			tableNames = append(tableNames, *scanDestination)
		}
	}
	_ = rs.Close()

	return tableNames, nil
}

// GetColumn returns a schema.Column representing the given tableName and colName.
// TODO: Charset and Collation
func (a *adapter) GetColumn(tableName, colName string) (schema.Column, error) {
	var (
		dt         string
		nullable   string
		key        = new(string)
		defaultVal = new(string)
		extra      string
		col        schema.Column
	)
	rs, err := a.db.Query(fmt.Sprintf(getColumnQuery, tableName, colName))
	if err != nil {
		return col, fmt.Errorf("unable to get column information for `%s`.`%s`: %w`", tableName, colName, err)
	}

	if !rs.Next() {
		return col, fmt.Errorf("unable to get column, empty result")
	}
	err = rs.Scan(new(interface{}), &dt, &nullable, &key, &defaultVal, &extra)
	if err != nil {
		return col, fmt.Errorf("unable to scan result reading column `%s`.`%s`: %w", tableName, colName, err)
	}

	err = rs.Close()
	if err != nil {
		return col, fmt.Errorf("unable to close rows after getting column, too many rows: %w", err)
	}

	dts := strings.Split(strings.ToUpper(dt), " ")
	col.Datatype, err = datatype.FromString(dts[0])
	if err != nil {
		return col, fmt.Errorf("unable to determine datatype for column `%s`.`%s`: %w", tableName, colName, err)
	}

	if col.Datatype.IsSignable() && len(dts) > 1 && dts[1] == "UNSIGNED" {
		col.Unsigned = true
	}

	ps := strings.Split(paramIsolator.ReplaceAllString(dts[0], ""), ",")

	if len(ps) > 0 {
		col.Params = ps
	}

	col.Default = defaultVal
	col.PrimaryKey = key != nil && *key == "PRI"
	col.AutoIncrement = strings.Contains(strings.ToLower(extra), "auto_increment")

	return col, nil
}

// GetIndex returns a schema.Index representing the given tableName and indexName.
func (a *adapter) GetIndex(tableName, indexName string) (schema.Index, error) {
	var (
		tempColName string
		columns     []string
		index       schema.Index
	)

	rs, err := a.db.Query(fmt.Sprintf(getIndexQuery, tableName, indexName))
	if err != nil {
		return index, fmt.Errorf("unable to get information for index `%s` on table `%s`: %w`", indexName, tableName, err)
	}

	for rs.Next() {
		err = rs.Scan(&index.Unique, &tempColName)
		if err != nil {
			return index, fmt.Errorf("unable to scan result reading index `%s` on table `%s`: %w", indexName, tableName, err)
		}

		columns = append(columns, tempColName)
	}
	_ = rs.Close()

	index.Columns = columns
	return index, nil
}

// GetReference returns a schema.Reference representing the given tableName and referenceName.
func (a *adapter) GetReference(tableName, referenceName string) (schema.Reference, error) {
	var (
		ref            schema.Reference
		constraintName string
		tempString     string
		columnNames    []string
	)

	rs, err := a.db.Query(fmt.Sprintf(fmt.Sprintf(getReferenceQuery, tableName, referenceName)))
	if err != nil {
		return ref, fmt.Errorf("unable to get reference information for table `%s` from table `%s`: %w", referenceName, tableName, err)
	}

	for rs.Next() {
		err = rs.Scan(&ref.OnUpdate, &ref.OnDelete, &constraintName)
	}
	_ = rs.Close()

	rs, err = a.db.Query(fmt.Sprintf(getReferenceColumnsQuery, tableName, referenceName, constraintName))
	if err != nil {
		return ref, fmt.Errorf("unable to get reference columns for table `%s` from table `%s`: %w", referenceName, tableName, err)
	}

	for rs.Next() {
		err = rs.Scan(&tempString, &ref.Required)
		columnNames = append(columnNames, tempString)
	}
	_ = rs.Close()

	ref.ColumnNames = columnNames

	if len(ref.ColumnNames) == 0 {
		return ref, fmt.Errorf("unable to find any reference columns for table `%s` from table `%s`", referenceName, tableName)
	}

	return ref, err
}
