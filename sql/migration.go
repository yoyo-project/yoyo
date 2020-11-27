package sql

import "database/sql"

type Migration struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}
