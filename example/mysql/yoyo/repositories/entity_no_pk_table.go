package repositories

import (
	"database/sql"
	"fmt"

	
)

type NoPkTable struct {
	Col uint32

	

	persisted *NoPkTable
}

type NoPkTables struct {
	rs *sql.Rows
}

func (es *NoPkTables) Next() bool {
	return es.rs.Next()
}

func (es *NoPkTables) Scan(e *NoPkTable) (err error) {
	if e == nil {
		return fmt.Errorf("in NoPkTables.Scan: passed a nil entity")
	}

	err = es.rs.Scan(&e.Col)
	persisted := *e
	e.persisted = &persisted
	return err
}
