package repository

import (
	"io"

	"github.com/dotvezz/yoyo/internal/repository/template"
)

func NewQueryNodeGenerator() SimpleWriteGenerator {
	return func(w io.StringWriter) error {
		_, err := w.WriteString(template.NodeFile)
		return err
	}
}
