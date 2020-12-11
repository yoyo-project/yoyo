package base

import (
	"github.com/dotvezz/yoyo/internal/datatype"
)

// Base is a partial implementation of migration.Dialect. It provides the TypeString method.
type Base struct {
	Dialect string
}

// TypeString is implemented because, with most datatype.Datatype values in most SQL dialects, the datatype.Datatype String()
// method will probably be correct. Dialect-specific implementations can extend this TypeString method with any specific
// exceptions. For one example, check the mysql.migrator's TypeString method.
func (d *Base) TypeString(dt datatype.Datatype) (string, error) {
	s := dt.String()

	if s == "NONE" {
		return "", datatype.ErrUnknownDatatype
	}

	return s, nil
}
