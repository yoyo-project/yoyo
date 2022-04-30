package repositories

import (
	"database/sql"
	"fmt"

	
)

type Person struct {
	Id int32
	Name string
	Nickname string
	FavoriteColor string
	Age float64

	CityId int32

	persisted *Person
}

type Persons struct {
	rs *sql.Rows
}

func (es *Persons) Next() bool {
	return es.rs.Next()
}

func (es *Persons) Scan(e *Person) (err error) {
	if e == nil {
		return fmt.Errorf("in Persons.Scan: passed a nil entity")
	}

	err = es.rs.Scan(&e.Id, &e.Name, &e.Nickname, &e.FavoriteColor, &e.Age)
	persisted := *e
	e.persisted = &persisted
	return err
}
