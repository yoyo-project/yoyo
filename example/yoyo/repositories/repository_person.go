package repositories

import (
	"database/sql"
	"fmt"
	
	"github.com/yoyo-project/yoyo/example/yoyo/repositories/query/person"
)

const (
	insertPerson = "INSERT INTO person" +
		" (name, favorite_color, age) " +
		" VALUES (?, ?, ?);"
	updatePerson = "UPDATE person" +
		" SET name = ?, favorite_color = ?, age = ? %s;"
	selectPerson = "SELECT name, favorite_color, age FROM person %s;"
	deletePerson = "DELETE FROM person %s;"
)

type PersonRepository struct {
	*repository
}

func (r *PersonRepository) FetchOne(query person.Query) (ent Person, err error) {
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
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(&ent.Id, &ent.Name, &ent.FavoriteColor, &ent.Age)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *PersonRepository) Search(query person.Query) (es Persons, err error) {
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
		return es, err
	}

	es.rs, err = stmt.Query(args...)

	return es, err
}

func (r *PersonRepository) Save(in Person) (Person, error) {
	if in.persisted == nil {
		return r.insert(in)
	} else {
		return r.update(in)
	}
}

func (r *PersonRepository) insert(in Person) (e Person, err error) {
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
		return e, err
	}

	res, err = stmt.Exec(in.Id, in.Name, in.FavoriteColor, in.Age)
	if err != nil {
		return e, err
	}

	e = in
	var eid int64
	eid, err = res.LastInsertId()
	e.Id = int32(eid)
	if err != nil {
		return e, err
	}

	in = e
	e.persisted = &in

	return e, err
}

func (r *PersonRepository) update(in Person) (e Person, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil {
			_ = stmt.Close()
		}
	}()

	q, args := person.Query{}.
		Id(in.persisted.Id).
		SQL()

	stmt, err = r.prepare(fmt.Sprintf(updatePerson, q))
	if err != nil {
		return e, err
	}

	fields := []interface{}{in.Id, in.Name, in.FavoriteColor, in.Age}
	_, err = stmt.Exec(append(fields, args...)...)
	if err != nil {
		return e, err
	}

	e = in
	in = e
	e.persisted = &in

	return e, err
}

func (r *PersonRepository) Delete(query person.Query) (err error) {
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

