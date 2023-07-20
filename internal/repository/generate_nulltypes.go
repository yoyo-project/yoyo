package repository

import (
	"github.com/yoyo-project/yoyo/internal/repository/template"
	"io"
)

func NewNullTypesFileGenerator() SimpleWriteGenerator {
	return func(w io.StringWriter) error {
		_, err := w.WriteString(template.NullTypeFile)
		return err
	}
}