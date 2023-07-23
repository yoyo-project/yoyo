package repositories

import (
	"database/sql"
	"fmt"

	
)

type NoPkTable struct {
	Col int32
	Col2 int32

	

	persisted *NoPkTable
}

type NoPkTables struct {
	// If we're not in a transaction, then NoPkTable saves memory by wrapping a *sql.Rows to scan from the connection
	// buffer on-demand.
	// This uses less application memory but more connections to the DBMS.
	rs *sql.Rows

	// If we are in a transaction, then NoPkTable reads the entire result set to memory to clear the buffer and allow
	// other queries to run on the goroutine.
	// This uses more application memory but fewer connections to the DBMS.
	i  int
	es []NoPkTable
}

func (es *NoPkTables) Next() bool {
	if es.rs != nil {
		return es.rs.Next()
	} else {
		es.i++
		return es.i < len(es.es)
	}
}

func (es *NoPkTables) Scan(e *NoPkTable) (err error) {
	if e == nil {
		return fmt.Errorf("in NoPkTables.Scan: passed a nil entity")
	}

	if es.rs != nil {
		return es.scan(e)
	}

	return es.point(e)
}

func (es *NoPkTables) scan(e *NoPkTable) (err error) {
	err = es.rs.Scan(&e.Col, &e.Col2)
	persisted := *e
	e.persisted = &persisted
	return err
}

func (es *NoPkTables) point(e *NoPkTable) (err error) {
	if es.i >= len(es.es) || es.i < 0 {
		return fmt.Errorf("in NoPkTables.point: out of range")
	}
	*e = es.es[es.i]
	persisted := *e
	e.persisted = &persisted
	return nil
}
