package postgres

import (
	"fmt"
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/dbms/base"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
)

func NewMigrator() *postgres {
	return &postgres{
		Base: base.Base{
			Dialect: dialect.PostgreSQL,
		},
		validator: validator{},
	}
}

type postgres struct {
	base.Base
	validator
}

func (d *postgres) TypeString(dt datatype.Datatype) (s string, err error) {
	if !d.SupportsDatatype(dt) {
		return "", fmt.Errorf("datatype %s is not supported in postgresql", dt)
	}
	switch dt {
	default:
		s, err = d.Base.TypeString(dt)
	}
	return s, err
}
