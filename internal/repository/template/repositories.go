package template

const (
	Imports              = "$IMPORTS"
	StructFields         = "$STRUCT_FIELDS$"
	RepositoryInterfaces = "$REPOSITORY_INTERFACES$"
	RepoInits            = "$REPO_INITS$"
)

const RepositoriesFile = `package ` + PackageName + `

import (
	"database/sql"

` + Imports + `
)

type Transact func(func() error) error

type Repositories struct {
` + StructFields + `
}

` + RepositoryInterfaces + `

func InitRepositories(db *sql.DB) (Repositories, Transact) {
	baseRepo := &repository{db: db}
	return Repositories{
` + RepoInits + `
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
`
