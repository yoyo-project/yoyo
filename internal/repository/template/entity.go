package template

const (
	EntityFields    = "$ENTITY_FIELDS$"
	ScanFields      = "$SCAN_FIELDS$"
	ReferenceFields = "$REFERENCE_FIELDS$"
)

const EntityFile = `package ` + PackageName + `

import (
	"database/sql"
	"fmt"

	` + Imports + `
)

type ` + EntityName + ` struct {
	` + EntityFields + `

	` + ReferenceFields + `

	persisted *` + EntityName + `
}

type ` + EntityName + `s struct {
	// If we're not in a transaction, then ` + EntityName + ` saves memory by wrapping a *sql.Rows to scan from the connection
	// buffer on-demand.
	// This uses less application memory but more connections to the DBMS.
	rs *sql.Rows

	// If we are in a transaction, then ` + EntityName + ` reads the entire result set to memory to clear the buffer and allow
	// other queries to run on the goroutine.
	// This uses more application memory but fewer connections to the DBMS.
	i  int
	es []` + EntityName + `
}

func (es *` + EntityName + `s) Next() bool {
	if es.rs != nil {
		return es.rs.Next()
	} else {
		es.i++
		return es.i < len(es.es)
	}
}

func (es *` + EntityName + `s) Scan(e *` + EntityName + `) (err error) {
	if e == nil {
		return fmt.Errorf("in ` + EntityName + `s.Scan: passed a nil entity")
	}

	if es.rs != nil {
		return es.scan(e)
	}

	return es.point(e)
}

func (es *` + EntityName + `s) scan(e *` + EntityName + `) (err error) {
	err = es.rs.Scan(` + ScanFields + `)
	persisted := *e
	e.persisted = &persisted
	return err
}

func (es *` + EntityName + `s) point(e *` + EntityName + `) (err error) {
	if es.i >= len(es.es) || es.i < 0 {
		return fmt.Errorf("in ` + EntityName + `s.point: out of range")
	}
	*e = es.es[es.i]
	persisted := *e
	e.persisted = &persisted
	return nil
}
`
