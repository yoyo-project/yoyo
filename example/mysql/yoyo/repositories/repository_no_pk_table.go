package repositories

import (
	"database/sql"
	"fmt"
	
	"github.com/yoyo-project/yoyo/example/mysql/yoyo/repositories/query/no_pk_table"
)

const (
	insertNoPkTable = "INSERT INTO no_pk_table" +
		" (col, col2) " +
		" VALUES (?, ?);"
	updateNoPkTable = "UPDATE no_pk_table" +
		" SET col = ?, col2 = ? %s;"
	selectNoPkTable = "SELECT col, col2 FROM no_pk_table %s;"
	deleteNoPkTable = "DELETE FROM no_pk_table %s;"
)

type NoPkTableRepository struct {
	*repository
}

func (r *NoPkTableRepository) FetchOne(query no_pk_table.Query) (ent NoPkTable, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectNoPkTable, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(&ent.Col, &ent.Col2)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *NoPkTableRepository) Search(query no_pk_table.Query) (es NoPkTables, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectNoPkTable, conditions))
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
			var ent NoPkTable
			err = rs.Scan(&ent.Col, &ent.Col2)
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

func (r *NoPkTableRepository) Save(in NoPkTable) (e NoPkTable, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(insertNoPkTable)
	if err != nil {
		return e, err
	}

	_, err = stmt.Exec(in.Col, in.Col2)
	if err != nil {
		return e, err
	}

	in = e
	e.persisted = &in

	return e, err
}

