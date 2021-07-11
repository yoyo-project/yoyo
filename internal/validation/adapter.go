package validation

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
	"github.com/yoyo-project/yoyo/internal/dbms/mysql"
	"github.com/yoyo-project/yoyo/internal/dbms/postgres"
	"github.com/yoyo-project/yoyo/internal/dbms/sqlite"
	"github.com/yoyo-project/yoyo/internal/schema"
)

// Adapter provides an interface for DBMS-specific validations
type Adapter interface {
	// SupportsDatatype takes a datatype.Datatype and returns true if the underlying DBMS supports it
	SupportsDatatype(datatype datatype.Datatype) bool
	// ValidateTable takes a schema.Table and returns an error if the underlying DBMS doesn't support it
	ValidateTable(col schema.Table) error
	// SupportsAutoIncrement returns true if the underlying DBMS supports AutoIncrement
	SupportsAutoIncrement() bool
}

func LoadValidator(name string) (a Adapter, err error) {
	switch name {
	case dialect.MySQL:
		a = mysql.NewAdapter()
	case dialect.PostgreSQL:
		a = postgres.NewAdapter()
	case dialect.SQLite:
		a = sqlite.NewAdapter()
	default:
		err = fmt.Errorf("unknown dialect `%s`", name)
	}

	return a, err
}
