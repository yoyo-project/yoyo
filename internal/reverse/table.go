package reverse

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/schema"
)

type TableSearcher func(table, x string) bool

func InitHasColumn(getColumn func(table, column string) (schema.Column, error)) TableSearcher {
	return func(table, col string) bool {
		_, err := getColumn(table, col)
		return err == nil
	}
}

func InitHasIndex(getIndex func(table, column string) (schema.Index, error)) TableSearcher {
	return func(table, col string) bool {
		_, err := getIndex(table, col)
		return err == nil
	}
}

func InitHasReference(getReference func(table, reference string) (schema.Reference, error)) TableSearcher {
	return func(table, col string) bool {
		_, err := getReference(table, col)
		return err == nil
	}
}

func InitHasTable(listTables func() ([]string, error)) func(table string) (bool, error) {
	return func(table string) (bool, error) {
		tables, err := listTables()
		if err != nil {
			return false, fmt.Errorf("unable to list tables: %w", err)
		}

		for _, t := range tables {
			if t == table {
				return true, nil
			}
		}

		return false, nil
	}
}
