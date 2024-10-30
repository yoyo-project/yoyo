package repositories

import (
	"database/sql"
	"fmt"

	"github.com/yoyo-project/yoyo/example/postgres/yoyo/repositories/nullable"
)

type Person struct {
	Id uint32
	Name string
	Nickname string
	FavoriteColor nullable.String
	Age float64

	CityId uint32

	persisted *Person
}

type Persons struct {
	// If we're not in a transaction, then Person saves memory by wrapping a *sql.Rows to scan from the connection
	// buffer on-demand.
	// This uses less application memory but more connections to the DBMS.
	rs *sql.Rows

	// If we are in a transaction, then Person reads the entire result set to memory to clear the buffer and allow
	// other queries to run on the goroutine.
	// This uses more application memory but fewer connections to the DBMS.
	i  int
	es []Person
}

func (es *Persons) Next() bool {
	if es.rs != nil {
		return es.rs.Next()
	} else {
		es.i++
		return es.i < len(es.es)
	}
}

func (es *Persons) Scan(e *Person) (err error) {
	if e == nil {
		return fmt.Errorf("in Persons.Scan: passed a nil entity")
	}

	if es.rs != nil {
		return es.scan(e)
	}

	return es.point(e)
}

func (es *Persons) scan(e *Person) (err error) {
	err = es.rs.Scan(&e.Id, &e.Name, &e.Nickname, &e.FavoriteColor, &e.Age, &e.CityId)
	persisted := *e
	e.persisted = &persisted
	return err
}

func (es *Persons) point(e *Person) (err error) {
	if es.i >= len(es.es) || es.i < 0 {
		return fmt.Errorf("in Persons.point: out of range")
	}
	*e = es.es[es.i]
	persisted := *e
	e.persisted = &persisted
	return nil
}
