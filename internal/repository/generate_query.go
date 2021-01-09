package repository

import (
	"io"
	"strings"

	"github.com/dotvezz/yoyo/internal/repository/template"
	"github.com/dotvezz/yoyo/internal/schema"
	"github.com/dotvezz/yoyo/internal/yoyo"
)

func NewQueryFileGenerator(config yoyo.Config) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		var methods, functions, imports []string
		for cn, c := range t.Columns {
			ms, fs, is := template.GenerateQueryLogic(cn, c)
			methods = append(methods, ms...)
			functions = append(functions, fs...)
			imports = append(imports, is...)
		}

		r := strings.NewReplacer(
			template.PackageName,
			t.QueryPackageName(),
			template.StdlibImports,
			strings.Join(imports, "\n    "),
			template.RepositoriesPackage,
			config.Paths.Repositories,
		)

		_, err := w.WriteString(r.Replace(template.QueryFile))
		if err != nil {
			return err
		}

		_, err = w.WriteString(strings.Join(methods, "\n"))
		if err != nil {
			return err
		}

		_, err = w.WriteString(strings.Join(functions, "\n"))
		if err != nil {
			return err
		}

		return err
	}
}
