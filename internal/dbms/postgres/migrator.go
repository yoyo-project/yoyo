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
	}
}

type postgres struct {
	base.Base
}

func (d *postgres) TypeString(dt datatype.Datatype) (s string, err error) {
	if dt&datatype.PostgreSQL != datatype.PostgreSQL {
		return "", fmt.Errorf("datatype %s is not supported in postgresql", dt)
	}
	switch dt {
	default:
		s, err = d.Base.TypeString(dt)
	}
	return s, err
}
