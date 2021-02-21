package repositories

import (
	"database/sql"
	"fmt"

	
)

type City struct {
	Id int32
	Name string

	

	persisted *City
}

type Citys struct {
	rs *sql.Rows
}

func (es *Citys) Next() bool {
	return es.rs.Next()
}

func (es *Citys) Scan(e * City) (err error) {
	if e == nil {
		return fmt.Errorf("in Citys.Scan: passed a nil entity")
	}

	err = es.rs.Scan(&e.Id, &e.Name)
	persisted := *e
	e.persisted = &persisted
	return err
}
