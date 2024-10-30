package repositories

import (
	"database/sql"
	"fmt"
	
	"github.com/yoyo-project/yoyo/example/postgres/yoyo/repositories/query/person"
)

const (
	insertPerson = "INSERT INTO person" +
		" (name, nickname, favorite_color, age, fk_city_id) " +
		" VALUES ($1, $2, $3, $4, $5, $6);"
	updatePerson = "UPDATE person" +
		" SET id = $1, name = $2, nickname = $3, favorite_color = $4, age = $5, fk_city_id = $6 %s;"
	selectPerson = "SELECT id, name, nickname, favorite_color, age, fk_city_id FROM person %s;"
	deletePerson = "DELETE FROM person %s;"
)

type PersonRepository struct {
	*repository
}

func (r *PersonRepository) FetchOne(query person.Query) (ent Person, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectPerson, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(&ent.Id, &ent.Name, &ent.Nickname, &ent.FavoriteColor, &ent.Age, &ent.CityId)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *PersonRepository) Search(query person.Query) (es Persons, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectPerson, conditions))
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
			var ent Person
			err = rs.Scan(&ent.Id, &ent.Name, &ent.Nickname, &ent.FavoriteColor, &ent.Age, &ent.CityId)
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
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	stmt, err = r.prepare(insertPerson)
	if err != nil {
		return e, err
	}

	res, err = stmt.Exec(in.Id, in.Name, in.Nickname, in.FavoriteColor, in.Age, in.CityId)
	if err != nil {
		return e, err
	}

	e = in
	var eid int64
	eid, err = res.LastInsertId()
	e.Id = uint32(eid)
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
		if stmt != nil && r.tx == nil {
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

	fields := []interface{}{in.Id, in.Name, in.Nickname, in.FavoriteColor, in.Age, in.CityId}
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
		if stmt != nil && r.tx == nil {
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

