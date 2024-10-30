package mysql

import "github.com/yoyo-project/yoyo/internal/schema"

func (adapter) IdentifierQuoteRune() rune {
	return '`'
}

func (adapter) StringQuoteRune() rune {
	return '"'
}

func (adapter) ValidateTable(t schema.Table) error {
	return nil
}

func (adapter) PreparedStatementPlaceholders(count int) []string {
	out := make([]string, count)
	for i := range out {
		out[i] = "?"
	}
	return out
}

func (adapter) PreparedStatementPlaceholderDef() (string, int) {
	return "?", 0
}
