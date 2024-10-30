package repository

import (
	"io"
	"strconv"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"
)

func NewQueryNodeGenerator(a Adapter) SimpleWriteGenerator {
	return func(w io.StringWriter) error {
		statement, add := a.PreparedStatementPlaceholderDef()
		identifierQuote := strconv.Quote(string(a.IdentifierQuoteRune()))
		identifierQuote = identifierQuote[1 : len(identifierQuote)-1]
		r := strings.NewReplacer(
			template.IdentifierQuote, identifierQuote,
			template.PlaceholderStatement, statement,
			template.PlaceholderAdd, strconv.Itoa(add),
		)
		_, err := w.WriteString(r.Replace(template.NodeFile))
		return err
	}
}
