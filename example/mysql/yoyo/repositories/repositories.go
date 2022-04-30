package repositories

import (
	"database/sql"
)

type TransactFunc func(func() error) error

type Repositories struct {
	*CityRepository
	*PersonRepository
}

func InitRepositories(db *sql.DB) (Repositories, TransactFunc) {
	baseRepo := &repository{db: db}
	return Repositories{
		CityRepository: &CityRepository{baseRepo},
		PersonRepository: &PersonRepository{baseRepo},
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

func initTransact(r *repository) TransactFunc {
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
