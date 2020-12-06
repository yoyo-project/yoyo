package mysql

import (
	"database/sql"
	"fmt"
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/reverse"
	"github.com/dotvezz/yoyo/internal/schema"
	goMysql "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

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

func InitNewReverser(open func(driver, dsn string) (*sql.DB, error)) func(host, user, dbname, password, port string) (reverse.Reverser, error) {
	return func(host, user, dbname, password, port string) (reverse.Reverser, error) {
		reverser := reverser{}
		cnf := goMysql.NewConfig()

		cnf.User = user
		cnf.Passwd = password
		cnf.Net = "tcp"
		cnf.Addr = host
		if port != "" {
			cnf.Addr += port
		}
		cnf.DBName = dbname

		var err error
		reverser.db, err = open("mysql", cnf.FormatDSN())
		if err != nil {
			return nil, fmt.Errorf("unable to open database connection for mysql reverser: %w", err)
		}

		return &reverser, nil
	}
}

type reverser struct {
	db *sql.DB
}

func (d *reverser) ListTables() ([]string, error) {
	rs, err := d.db.Query("SHOW TABLES")
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

func (d *reverser) ListIndices(table string) ([]string, error) {
	query := fmt.Sprintf(listIndicesQuery, table)
	rs, err := d.db.Query(query)
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

func (d *reverser) ListColumns(table string) ([]string, error) {
	rs, err := d.db.Query(fmt.Sprintf(listColumnsQuery, table))
	if err != nil {
		return nil, fmt.Errorf("unable to list indices: %w", err)
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

func (d *reverser) ListReferences(table string) ([]string, error) {
	rs, err := d.db.Query(fmt.Sprintf(listReferencesQuery, table))
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

// TODO: Charset and Collation
func (d *reverser) GetColumn(tableName, colName string) (schema.Column, error) {
	var (
		dt         string
		nullable   string
		key        = new(string)
		defaultVal = new(string)
		extra      string
		col        schema.Column
	)
	rs, err := d.db.Query(fmt.Sprintf(getColumnQuery, tableName, colName))
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

	if strings.Contains(dts[0], "(") {
		dts = strings.Split(dts[0], "(")
		if len(dts) > 1 {
			dts = strings.Split(dts[1], ",")
			col.Precision, err = strconv.Atoi(strings.Trim(dts[0], ",()"))
			if err != nil {
				return col, fmt.Errorf("unable to determine precision/length for column `%s`.`%s`: %w", tableName, colName, err)
			}
		}
		if col.Datatype.IsNumeric() && len(dts) > 1 {
			col.Scale, err = strconv.Atoi(strings.Trim(dts[1], ")"))
			if err != nil {
				return col, fmt.Errorf("unable to determine scale for column `%s`.`%s`: %w", tableName, colName, err)
			}
		}
	}

	col.Default = defaultVal
	col.PrimaryKey = key != nil && *key == "PRI"
	col.AutoIncrement = strings.Contains(strings.ToLower(extra), "auto_increment")

	return col, nil
}

func (d *reverser) GetIndex(tableName, indexName string) (schema.Index, error) {
	var (
		tempColName string
		columns     []string
		index       schema.Index
	)

	rs, err := d.db.Query(fmt.Sprintf(getIndexQuery, tableName, indexName))
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

func (d *reverser) GetReference(tableName, referenceName string) (schema.Reference, error) {
	var (
		ref            schema.Reference
		constraintName string
		tempString     string
		columnNames    []string
	)

	rs, err := d.db.Query(fmt.Sprintf(fmt.Sprintf(getReferenceQuery, tableName, referenceName)))
	if err != nil {
		return ref, fmt.Errorf("unable to get reference information for table `%s` from table `%s`: %w", referenceName, tableName, err)
	}

	for rs.Next() {
		err = rs.Scan(&ref.OnUpdate, &ref.OnDelete, &constraintName)
	}
	_ = rs.Close()

	rs, err = d.db.Query(fmt.Sprintf(getReferenceColumnsQuery, tableName, referenceName, constraintName))
	if err != nil {
		return ref, fmt.Errorf("unable to get reference columns for table `%s` from table `%s`: %w", referenceName, tableName, err)
	}

	for rs.Next() {
		err = rs.Scan(&tempString, &ref.Optional)
		columnNames = append(columnNames, tempString)
	}
	_ = rs.Close()

	switch len(columnNames) {
	case 0:
		return ref, fmt.Errorf("unable to find any reference columns for table `%s` from table `%s`", referenceName, tableName)
	case 1:
		ref.ColumnName = columnNames[0]
	default:
		ref.ColumnNames = columnNames
	}

	return ref, err
}
