package template

const (
	Imports               = "$IMPORTS$"
	TableName             = "$TABLE_NAME$"
	QueryPackageName      = "$ENTITY_PACKAGE_NAME$"
	InsertColumnNames     = "$INSERT_COLUMN_NAMES$"
	SelectColumnNames     = "$SELECT_COLUMN_NAMES$"
	StatementPlaceholders = "$STATEMENT_PLACEHOLDERS$"
	ColumnAssignments     = "$COLUMN_ASSIGNMENTS$"
	EntityName            = "$ENTITY_NAME$"
	InFields              = "$IN_FIELDS$"
	PKCapture             = "$ID_CAPTURE$"
	FieldName             = "$PK_FIELD_NAME$"
	PKQuery               = "$PK_QUERY$"
	PKFields              = "$PK_FIELDS"
	SaveFuncs             = "$SAVE_FUNCS$"
)

const SaveWithPK = `func (r *` + EntityName + `Repository) Save(in ` + EntityName + `) (` + EntityName + `, error) {
	if in.persisted == nil {
		return r.insert(in)
	} else {
		return r.update(in)
	}
}

func (r *` + EntityName + `Repository) insert(in ` + EntityName + `) (e ` + EntityName + `, err error) {
	var (
		stmt *sql.Stmt
		res  sql.Result
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(insert` + EntityName + `)
	if err != nil {
		return e, err
	}

	res, err = stmt.Exec(` + InFields + `)
	if err != nil {
		return e, err
	}
` + PKCapture + `
	in = e
	e.persisted = &in

	return e, err
}

func (r *` + EntityName + `Repository) update(in ` + EntityName + `) (e ` + EntityName + `, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()
` + PKQuery + `
	stmt, err = r.prepare(fmt.Sprintf(update` + EntityName + `, q))
	if err != nil {
		return e, err
	}

	fields := []interface{}{` + InFields + `}
	_, err = stmt.Exec(append(fields, args...)...)
	if err != nil {
		return e, err
	}

	e = in
	in = e
	e.persisted = &in

	return e, err
}

func (r *` + EntityName + `Repository) Delete(query ` + QueryPackageName + `.Query) (err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(delete` + EntityName + `, conditions))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)

	return err
}`

const SaveWithoutPK = `func (r *` + EntityName + `Repository) Save(in ` + EntityName + `) (e ` + EntityName + `, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(insert` + EntityName + `)
	if err != nil {
		return e, err
	}

	_, err = stmt.Exec(` + InFields + `)
	if err != nil {
		return e, err
	}
` + PKCapture + `
	in = e
	e.persisted = &in

	return e, err
}`

const RepositoryFile = `package ` + PackageName + `

import (
	"database/sql"
	"fmt"
	
	` + Imports + `
)

const (
	insert` + EntityName + ` = "INSERT INTO ` + TableName + `" +
		" (` + InsertColumnNames + `) " +
		" VALUES (` + StatementPlaceholders + `);"
	update` + EntityName + ` = "UPDATE ` + TableName + `" +
		" SET ` + ColumnAssignments + ` %s;"
	select` + EntityName + ` = "SELECT ` + SelectColumnNames + ` FROM ` + TableName + ` %s;"
	delete` + EntityName + ` = "DELETE FROM ` + TableName + ` %s;"
)

type ` + EntityName + `Repository struct {
	*repository
}

func (r *` + EntityName + `Repository) FetchOne(query ` + QueryPackageName + `.Query) (ent ` + EntityName + `, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(select` + EntityName + `, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(` + ScanFields + `)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *` + EntityName + `Repository) Search(query ` + QueryPackageName + `.Query) (es ` + EntityName + `s, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(select` + EntityName + `, conditions))
	if err != nil {
		return es, err
	}

	// If we're in a transaction, take the full result set into memory to free up the sql connection's buffer
	if r.tx != nil {
		var rs *sql.Rows
		rs, err = stmt.Query()
		if err != nil {
			return es, err
		}

		for rs.Next() {
			var ent ` + EntityName + `
			err = rs.Scan(` + ScanFields + `)
			if err != nil {
				return es, err
			}
			es.es = append(es.es, ent)
		}

		es.i = -1

		return es, nil
	}

	es.rs, err = stmt.Query(args...)

	return es, err
}

` + SaveFuncs + `

`

const NoPKCapture = `
	_ = res
`

const SinglePKCaptureTemplate = `
	e = in
	var eid int64
	eid, err = res.LastInsertId()
	e.` + FieldName + ` = ` + Type + `(eid)
	if err != nil {
		return e, err
	}
`

const MultiPKCaptureTemplate = `
	e = in
	var eid int64
	eid, err = res.LastInsertId()
	e.Id = int32(eid)
	if err != nil {
		return e, err
	}
`

const PKQueryTemplate = `
	q, args := ` + QueryPackageName + `.Query{}.
		` + PKFields + `
		SQL()
`

const PKFieldTemplate = FieldName + "(in.persisted." + FieldName + ")."
