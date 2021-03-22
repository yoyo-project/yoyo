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
	rs *sql.Rows
}

func (es *` + EntityName + `s) Next() bool {
	return es.rs.Next()
}

func (es *` + EntityName + `s) Scan(e *` + EntityName + `) (err error) {
	if e == nil {
		return fmt.Errorf("in ` + EntityName + `s.Scan: passed a nil entity")
	}

	err = es.rs.Scan(` + ScanFields + `)
	persisted := *e
	e.persisted = &persisted
	return err
}
`
