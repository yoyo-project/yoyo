package postgres

import (
	"fmt"
	"github.com/yoyo-project/yoyo/internal/dbms/base"
	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
)

// NewAdapter returns an implementation of migration.Dialect for PostgreSQL
func NewAdapter() adapter {
	return adapter{
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

func (adapter) PreparedStatementPlaceholders(count int) []string {
	ps := make([]string, count)
	for i := 0; i < count; i++ {
		ps[i] = fmt.Sprintf("$%d", i+1)
	}
	return ps
}

func (adapter) PreparedStatementPlaceholderDef() (string, int) {
	return "$%d", 1
}

func (adapter) IdentifierQuoteRune() rune {
	return '"'
}

func (adapter) StringQuoteRune() rune {
	return '\''
}
