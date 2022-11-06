package repositories

import (
	"database/sql"
	"fmt"
	
	"github.com/yoyo-project/yoyo/example/mysql/yoyo/repositories/query/no_pk_table"
)

const (
	insertNoPkTable = "INSERT INTO no_pk_table" +
		" (col) " +
		" VALUES (?);"
	updateNoPkTable = "UPDATE no_pk_table" +
		" SET col = ? %s;"
	selectNoPkTable = "SELECT col FROM no_pk_table %s;"
	deleteNoPkTable = "DELETE FROM no_pk_table %s;"
)

type NoPkTableRepository struct {
	*repository
}

func (r *NoPkTableRepository) FetchOne(query no_pk_table.Query) (ent NoPkTable, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectNoPkTable, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(&ent.Col)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *NoPkTableRepository) Search(query no_pk_table.Query) (es NoPkTables, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectNoPkTable, conditions))
	if err != nil {
		return es, err
	}

	es.rs, err = stmt.Query(args...)

	return es, err
}

func (r *NoPkTableRepository) Save(in NoPkTable) (e NoPkTable, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(insertNoPkTable)
	if err != nil {
		return e, err
	}

	_, err = stmt.Exec(in.Col)
	if err != nil {
		return e, err
	}

	in = e
	e.persisted = &in

	return e, err
}

