package repositories

import (
	"database/sql"
	"fmt"

	
)

type Person struct {
	FavoriteColor string
	Id int32
	Name string

	

	persisted *Person
}

type Persons struct {
	rs *sql.Rows
}

func (es *Persons) Next() bool {
	return es.rs.Next()
}

func (es *Persons) Scan(e * Person) (err error) {
	if e == nil {
		return fmt.Errorf("in Persons.Scan: passed a nil entity")
	}

	err = es.rs.Scan(&e.FavoriteColor, &e.Id, &e.Name)
	persisted := *e
	e.persisted = &persisted
	return err
}
