package mysql

import (
	"database/sql"

	"github.com/yoyo-project/yoyo/internal/dbms/base"
	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
)

type adapter struct {
	db *sql.DB
	base.Base
}

// NewAdapter returns an implementation of migration.Dialect for MySQL
func NewAdapter() adapter {
	return adapter{
		Base: base.Base{
			Dialect: dialect.MySQL,
		},
	}
}
