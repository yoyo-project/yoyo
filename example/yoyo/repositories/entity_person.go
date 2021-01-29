package repositories

import (
	"database/sql"
	"fmt"
	"time"
)

type Person struct {
	ID            uint
	Name          string
	FavoriteColor string
	Birthday      time.Time

	StateID *int

	persisted *Person
}

type Persons struct {
	rs *sql.Rows
}

func (s *Persons) Next() bool {
	return s.rs.Next()
}

func (s *Persons) Scan(person *Person) (err error) {
	if person == nil {
		return fmt.Errorf("in States.Scan: passed a nil entity")
	}

	err = s.rs.Scan(&person.ID, &person.Name, &person.FavoriteColor)
	persisted := *person
	person.persisted = &persisted
	return err
}
