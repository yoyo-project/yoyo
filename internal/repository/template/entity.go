package template

const (
	EntityFields    = "$ENTITY_FIELDS$"
	ReferenceFields = "$REFERENCE_FIELDS$"
)

const EntityFile = `package repositories

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

func (s *` + EntityName + `s) Next() bool {
	return s.rs.Next()
}

func (s *` + EntityName + `s) Scan(e * ` + EntityName + `) (err error) {
	if e == nil {
		return fmt.Errorf("in ` + EntityName + `s.Scan: passed a nil entity")
	}

	err = s.rs.Scan(&e.ID, &e.Name, &e.FavoriteColor)
	persisted := *e
	e.persisted = &persisted
	return err
}
`
