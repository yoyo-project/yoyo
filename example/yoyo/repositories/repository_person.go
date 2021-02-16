package repositories

import (
	"database/sql"
	"fmt"

	"github.com/yoyo-project/yoyo/example/yoyo/repositories/query/person"
)

const (
	insertPerson = "INSERT INTO person" +
		" (id, name, favorite_color, birthday, fk_state_id) " +
		" VALUES (?, ?, ?, ?, ?);"
	updatePerson = "UPDATE person" +
		" SET id = ?, name = ?, favorite_color = ?, birthday = ?, fk_state_id = ?" +
		" WHERE id = ?;"
	selectPerson = "SELECT * FROM person%s;"
	deletePerson = "DELETE FROM person%s;"
)

type personRepo struct {
	*repository
}

func (r *personRepo) FetchOne(search person.Query) (p Person, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := search.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectPerson, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(&p)

	persisted := p
	p.persisted = &persisted

	return p, err
}

func (r *personRepo) Search(query person.Query) (ps Persons, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectPerson, conditions))
	if err != nil {
		return ps, err
	}

	ps.rs, err = stmt.Query(args...)

	return ps, err
}

func (r *personRepo) Save(in Person) (Person, error) {
	if in.persisted == nil {
		return r.insert(in)
	} else {
		return r.update(in)
	}
}

func (r *personRepo) insert(in Person) (rp Person, err error) {
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

	stmt, err = r.prepare(insertPerson)
	if err != nil {
		return rp, err
	}

	res, err = stmt.Exec(in.ID, in.Name, in.FavoriteColor)
	if err != nil {
		return rp, err
	}

	rp = in
	var tid int64
	tid, err = res.LastInsertId()
	rp.ID = uint(tid)
	if err != nil {
		return rp, err
	}

	in = rp
	rp.persisted = &in

	return rp, err
}

func (r *personRepo) update(in Person) (out Person, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(updatePerson)
	if err != nil {
		return out, err
	}

	_, err = stmt.Exec(in.ID, in.Name, in.FavoriteColor, in.Birthday, in.StateID, in.persisted.ID)
	if err != nil {
		return out, err
	}

	out = in
	in = out
	out.persisted = &in

	return out, err
}

func (r *personRepo) Delete(query person.Query) (err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(deletePerson, conditions))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)

	return err
}
