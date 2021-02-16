package mysql

import (
	"database/sql"

	"github.com/yoyo-project/yoyo/internal/dbms/base"
)

type adapter struct {
	db *sql.DB
	base.Base
	validator
}
