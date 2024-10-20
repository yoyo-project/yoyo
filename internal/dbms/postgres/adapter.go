package postgres

import (
	"fmt"
	"github.com/yoyo-project/yoyo/internal/dbms/base"
	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
)

// NewAdapter returns an implementation of migration.Dialect for PostgreSQL
func NewAdapter() *adapter {
	return &adapter{
		Base: base.Base{
			Dialect: dialect.PostgreSQL,
		},
		validator: validator{},
	}
}

type adapter struct {
	base.Base
	validator
	reverser
}

func (a *adapter) PreparedStatementPlaceholders(count int) (ps []string) {
	for i := 1; i <= count; i++ {
		ps = append(ps, fmt.Sprintf("$%d", i))
	}
	return ps
}
