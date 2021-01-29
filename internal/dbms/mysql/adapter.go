package mysql

import (
	"database/sql"

	"github.com/dotvezz/yoyo/internal/dbms/base"
)

type adapter struct {
	db *sql.DB
	base.Base
	validator
}
