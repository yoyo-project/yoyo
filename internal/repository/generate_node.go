package repository

import (
	"io"

	"github.com/yoyo-project/yoyo/internal/repository/template"
)

func NewQueryNodeGenerator() SimpleWriteGenerator {
	return func(w io.StringWriter) error {
		_, err := w.WriteString(template.NodeFile)
		return err
	}
}
