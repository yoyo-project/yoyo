package repositories

import (
	"context"
	"database/sql"
)

type TransactFunc func(func() error, ...TransactOptions) error

type TransactOptions struct {
	sql.TxOptions

	Context context.Context
}

type Repositories struct {
	*NoPkTableRepository
	*CityRepository
	*PersonRepository
}

func InitRepositories(db *sql.DB) (Repositories, TransactFunc) {
	baseRepo := &repository{db: db}
	return Repositories{
		CityRepository: &CityRepository{baseRepo},
		NoPkTableRepository: &NoPkTableRepository{baseRepo},
		PersonRepository: &PersonRepository{baseRepo},
	}, initTransact(baseRepo)
}

type repository struct {
	db *sql.DB
	tx *sql.Tx
}

func (r repository) prepare(query string) (*sql.Stmt, error) {
	if r.tx != nil {
		return r.tx.Prepare(query)
	} else {
		return r.db.Prepare(query)
	}
}

func initTransact(r *repository) TransactFunc {
	return func(f func() error, options ...TransactOptions) (err error) {
		var opts *sql.TxOptions
		ctx := context.Background()
		if len(options) > 0 {
			opts = &options[0].TxOptions
			if options[0].Context != nil {
				ctx = options[0].Context
			}
		}
		r.tx, err = r.db.BeginTx(ctx, opts)
		defer func() {
			if err != nil {
				err = r.tx.Rollback()
			} else {
				err = r.tx.Commit()
			}

			r.tx = nil
		}()

		if err == nil {
			err = f()
		}

		return
	}
}
