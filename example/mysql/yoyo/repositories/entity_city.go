package repositories

import (
	"database/sql"
	"fmt"

	
)

type City struct {
	Id uint32
	Name string

	

	persisted *City
}

type Citys struct {
	// If we're not in a transaction, then City saves memory by wrapping a *sql.Rows to scan from the connection
	// buffer on-demand.
	// This uses less application memory but more connections to the DBMS.
	rs *sql.Rows

	// If we are in a transaction, then City reads the entire result set to memory to clear the buffer and allow
	// other queries to run on the goroutine.
	// This uses more application memory but fewer connections to the DBMS.
	i  int
	es []City
}

func (es *Citys) Next() bool {
	if es.rs != nil {
		return es.rs.Next()
	} else {
		es.i++
		return es.i < len(es.es)
	}
}

func (es *Citys) Scan(e *City) (err error) {
	if e == nil {
		return fmt.Errorf("in Citys.Scan: passed a nil entity")
	}

	if es.rs != nil {
		return es.scan(e)
	}

	return es.point(e)
}

func (es *Citys) scan(e *City) (err error) {
	err = es.rs.Scan(&e.Id, &e.Name)
	persisted := *e
	e.persisted = &persisted
	return err
}

func (es *Citys) point(e *City) (err error) {
	if es.i >= len(es.es) || es.i < 0 {
		return fmt.Errorf("in Citys.point: out of range")
	}
	*e = es.es[es.i]
	persisted := *e
	e.persisted = &persisted
	return nil
}
