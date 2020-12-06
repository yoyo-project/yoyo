package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dotvezz/yoyo/internal/reverse"
	"github.com/dotvezz/yoyo/internal/schema"
	_ "github.com/lib/pq"
)

func InitNewReverser(open func(driver, dsn string) (*sql.DB, error)) func(host, user, dbname, password, port string) (reverse.Reverser, error) {
	return func(host, user, dbname, password, port string) (reverse.Reverser, error) {
		reverser := reverser{}

		var err error
		reverser.db, err = open("postgresql", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
		if err != nil {
			return nil, fmt.Errorf("unable to open database connection for mysql reverser: %w", err)
		}

		return &reverser, nil
	}
}

type reverser struct {
	db *sql.DB
}

func (r reverser) ListTables() ([]string, error) {
	panic("implement me")
}

func (r reverser) ListColumns(table string) ([]string, error) {
	panic("implement me")
}

func (r reverser) ListIndices(table string) ([]string, error) {
	panic("implement me")
}

func (r reverser) ListReferences(table string) ([]string, error) {
	panic("implement me")
}

func (r reverser) GetColumn(table, column string) (schema.Column, error) {
	panic("implement me")
}

func (r reverser) GetIndex(table, column string) (schema.Index, error) {
	panic("implement me")
}

func (r reverser) GetReference(table, column string) (schema.Reference, error) {
	panic("implement me")
}
