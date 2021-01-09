package repository

import (
	"io"

	"github.com/dotvezz/yoyo/internal/schema"
)

func NewRepositoriesGenerator() WriteGenerator {
	return func(db schema.Database, w io.StringWriter) error {

		return nil
	}
}
