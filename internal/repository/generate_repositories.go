package repository

import (
	"io"

	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewRepositoriesGenerator() WriteGenerator {
	return func(db schema.Database, w io.StringWriter) error {

		return nil
	}
}
