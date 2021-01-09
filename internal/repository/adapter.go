package repository

import (
	"fmt"

	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/dbms/mysql"
	"github.com/dotvezz/yoyo/internal/dbms/postgres"
	"github.com/dotvezz/yoyo/internal/schema"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

type AdapterLoader func(dia string) (Adapter, error)
type DatabaseReader func(config yoyo.Config) (db schema.Database, err error)
type AdapterBuilder func(host, userName, dbName, password, port string) (Adapter, error)

// Adapter is the yoyo interface for creating repository code
type Adapter interface {
	PreparedStatementPlaceholders(count int) []string
}

func LoadAdapter(dia string) (adapter Adapter, err error) {
	switch dia {
	case dialect.MySQL:
		adapter = mysql.NewAdapter()
	case dialect.PostgreSQL:
		adapter = postgres.NewAdapter()
	default:
		err = fmt.Errorf("unknown dialect `%s`", dia)
	}

	return adapter, err
}
