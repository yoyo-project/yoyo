package template

const (
	TableName             = "$TABLE_NAME$"
	QueryPackageName      = "$ENTITY_PACKAGE_NAME$"
	EntityFieldPointers   = "$ENTITY_FIELD_POINTERS$"
	ColumnNames           = "$COLUMN_NAMES$"
	StatementPlaceholders = "$STATEMENT_PLACEHOLDERS$"
	PrimaryKeyCondition   = "$PRIMARY_KEY_CONDITION$"
	ColumnAssignments     = "$COLUMN_ASSIGNMENTS$"
	EntityName            = "$ENTITY_NAME$"
	InProperties          = "$IN_PROPERTIES$"
)

const RepositoryFile = `package ` + PackageName + `

import (
	"database/sql"
	"fmt"
	
	` + Imports + `
)

const (
	insert` + EntityName + ` = "INSERT INTO ` + TableName + `" +
		" ( ` + ColumnNames + `) " +
		" VALUES ( ` + StatementPlaceholders + ` );"
	update` + EntityName + ` = "UPDATE ` + TableName + `" +
		" SET ` + ColumnAssignments + `" +
		" WHERE ` + PrimaryKeyCondition + `
	select` + EntityName + ` = "SELECT ` + ColumnNames + ` FROM ` + TableName + `%s;"
	delete` + EntityName + ` = "DELETE FROM ` + TableName + `%s;"
)

type ` + QueryPackageName + `Repo struct {
	*repository
}

func (r *` + QueryPackageName + `Repo) FetchOne(query ` + QueryPackageName + `.Query) (ent ` + EntityName + `, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(select` + EntityName + `, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(` + EntityFieldPointers + `)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *` + QueryPackageName + `Repo) Search(query ` + QueryPackageName + `.Query) (es ` + EntityName + `s, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(select` + EntityName + `, conditions))
	if err != nil {
		return es, err
	}

	es.rs, err = stmt.Query(args...)

	return es, err
}

func (r *` + QueryPackageName + `Repo) Save(in ` + EntityName + `) (` + EntityName + `, error) {
	if in.persisted == nil {
		return r.insert(in)
	} else {
		return r.update(in)
	}
}

func (r *` + QueryPackageName + `Repo) insert(in ` + EntityName + `) (e ` + EntityName + `, err error) {
	var (
		stmt *sql.Stmt
		res  sql.Result
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(insert` + EntityName + `)
	if err != nil {
		return e, err
	}

	res, err = stmt.Exec(` + InProperties + `)
	if err != nil {
		return e, err
	}

	e = in
	var eid int64
	eid, err = res.LastInsertId()
	e.ID = uint(eid)
	if err != nil {
		return e, err
	}

	in = e
	e.persisted = &in

	return e, err
}

func (r *` + QueryPackageName + `Repo) update(in ` + EntityName + `) (e ` + EntityName + `, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(update` + EntityName + `)
	if err != nil {
		return e, err
	}

	_, err = stmt.Exec(in.ID, in.Name, in.FavoriteColor, in.Birthday, in.StateID, in.persisted.ID)
	if err != nil {
		return e, err
	}

	e = in
	in = e
	e.persisted = &in

	return e, err
}

func (r *` + QueryPackageName + `Repo) Delete(query ` + EntityName + `.Query) (err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
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
}

`
