package repositories

import (
	"database/sql"

	"github.com/dotvezz/yoyo/example/yoyo/repositories/query/person"
	"github.com/dotvezz/yoyo/example/yoyo/repositories/query/state"
)

type Transact func(func() error) error

type Repositories struct {
	PersonRepository
	StateRepository
}

type PersonRepository interface {
	FetchOne(person.Query) (Person, error)
	Search(person.Query) (Persons, error)
	Save(Person) (Person, error)
	Delete(person.Query) error
}

type StateRepository interface {
	FetchOne(state.Query) (State, error)
	Search(state.Query) (States, error)
	Save(State) (State, error)
	Delete(state.Query)
}

func InitRepositories(db *sql.DB) (Repositories, Transact) {
	baseRepo := &repository{db: db}
	return Repositories{
		PersonRepository: &personRepo{baseRepo},
		StateRepository:  &stateRepo{baseRepo},
	}, initTransact(baseRepo)
}

type repository struct {
	db   *sql.DB
	tx   *sql.Tx
	isTx bool
}

func (r repository) prepare(query string) (*sql.Stmt, error) {
	if r.isTx {
		return r.tx.Prepare(query)
	} else {
		return r.db.Prepare(query)
	}
}

func initTransact(r *repository) Transact {
	return func(f func() error) (err error) {
		r.tx, err = r.db.Begin()
		r.isTx = true
		defer func() {
			r.isTx = false
			if err != nil {
				err = r.tx.Rollback()
			} else {
				err = r.tx.Commit()
			}
		}()

		if err == nil {
			err = f()
		}

		return
	}
}
