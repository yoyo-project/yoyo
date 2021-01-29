package repository

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/dotvezz/yoyo/internal/repository/template"
	"github.com/dotvezz/yoyo/internal/schema"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

func NewEntityRepositoryGenerator(cfg yoyo.Config, adapter Adapter) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		_, packageName := filepath.Split(cfg.Paths.Repositories)
		cNames := make([]string, 0)
		for columnName := range t.Columns {
			cNames = append(cNames, columnName)
		}

		r := strings.NewReplacer(
			template.PackageName,
			packageName,
			template.Imports,
			"",
			template.EntityName,
			t.ExportedGoName(),
			template.TableName,
			t.TableName(),
			template.ColumnNames,
			strings.Join(cNames, ", "),
			template.StatementPlaceholders,
			strings.Join(adapter.PreparedStatementPlaceholders(len(cNames)), ", "),
			template.ColumnAssignments,
			"",
			template.PrimaryKeyCondition,
			"",
			template.QueryPackageName,
			t.QueryPackageName(),
		)

		r.Replace(template.RepositoryFile)

		return nil
	}
}
