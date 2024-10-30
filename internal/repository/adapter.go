package repository

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
	"github.com/yoyo-project/yoyo/internal/dbms/mysql"
	"github.com/yoyo-project/yoyo/internal/dbms/postgres"
	"github.com/yoyo-project/yoyo/internal/schema"
	"github.com/yoyo-project/yoyo/internal/yoyo"
)

type AdapterLoader func(dia string) (Adapter, error)
type DatabaseReader func(config yoyo.Config) (db schema.Database, err error)
type AdapterBuilder func(host, userName, dbName, password, port string) (Adapter, error)

// Adapter is the yoyo interface for creating repository code
type Adapter interface {
	PreparedStatementPlaceholders(count int) []string
	PreparedStatementPlaceholderDef() (string, int)
	IdentifierQuoteRune() rune
	StringQuoteRune() rune
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
