package repositories

import (
	"database/sql"
	"fmt"
	
	"github.com/yoyo-project/yoyo/example/mysql/yoyo/repositories/query/city"
)

const (
	insertCity = "INSERT INTO city" +
		" (name) " +
		" VALUES (?, ?);"
	updateCity = "UPDATE city" +
		" SET id = ?, name = ? %s;"
	selectCity = "SELECT id, name FROM city %s;"
	deleteCity = "DELETE FROM city %s;"
)

type CityRepository struct {
	*repository
}

func (r *CityRepository) FetchOne(query city.Query) (ent City, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectCity, conditions))
	if err != nil {
		return
	}

	row := stmt.QueryRow(args...)

	err = row.Scan(&ent.Id, &ent.Name)

	persisted := ent
	ent.persisted = &persisted

	return ent, err
}

func (r *CityRepository) Search(query city.Query) (es Citys, err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(selectCity, conditions))
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
			var ent City
			err = rs.Scan(&ent.Id, &ent.Name)
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

func (r *CityRepository) Save(in City) (City, error) {
	if in.persisted == nil {
		return r.insert(in)
	} else {
		return r.update(in)
	}
}

func (r *CityRepository) insert(in City) (e City, err error) {
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

	stmt, err = r.prepare(insertCity)
	if err != nil {
		return e, err
	}

	res, err = stmt.Exec(in.Id, in.Name)
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

func (r *CityRepository) update(in City) (e City, err error) {
	var (
		stmt *sql.Stmt
	)
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	q, args := city.Query{}.
		Id(in.persisted.Id).
		SQL()

	stmt, err = r.prepare(fmt.Sprintf(updateCity, q))
	if err != nil {
		return e, err
	}

	fields := []interface{}{in.Id, in.Name}
	_, err = stmt.Exec(append(fields, args...)...)
	if err != nil {
		return e, err
	}

	e = in
	in = e
	e.persisted = &in

	return e, err
}

func (r *CityRepository) Delete(query city.Query) (err error) {
	var stmt *sql.Stmt
	// ensure the *sql.Stmt is closed after we're done with it
	defer func() {
		if stmt != nil && r.tx == nil {
			_ = stmt.Close()
		}
	}()

	conditions, args := query.SQL()
	stmt, err = r.prepare(fmt.Sprintf(deleteCity, conditions))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)

	return err
}

