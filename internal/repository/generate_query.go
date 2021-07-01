package repository

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"
	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewQueryFileGenerator(reposPath string, findPackagePath Finder, db schema.Database) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		var methods, functions, imports []string
		for _, c := range t.Columns {
			ms, fs, is := template.GenerateQueryLogic(c.Name, c)
			methods = append(methods, ms...)
			functions = append(functions, fs...)
			imports = append(imports, is...)
		}

		for _, r := range t.References {
			if r.HasMany {
				continue // Skip HasMany references which require a join
			}

			ft, ok := db.GetTable(r.TableName)
			if !ok {
				return fmt.Errorf("unable to generate queries for table %s, missing foreign table %s", t.Name, r.TableName)
			}

			var ms, fs, is []string

			for i, n := range r.ColNames(ft) {
				c := ft.PKColumns()[i]
				// Override the GoName in order to generate correct method/function names
				c.GoName = ft.ExportedGoName() + c.ExportedGoName()
				ms, fs, is = template.GenerateQueryLogic(n, c)
			}

			methods = append(methods, ms...)
			functions = append(functions, fs...)
			imports = append(imports, is...)
		}

		importPath, err := findPackagePath(reposPath + "/")
		if err != nil {
			return fmt.Errorf("unable to generate query file: %w", err)
		}

		r := strings.NewReplacer(
			template.PackageName,
			t.QueryPackageName(),
			template.StdlibImports,
			strings.Join(sortedUnique(imports), "\n	"),
			template.RepositoriesPackage,
			importPath,
		)

		_, err = w.WriteString(r.Replace(template.QueryFile))
		if err != nil {
			return err
		}

		_, err = w.WriteString(strings.Join(sortedUnique(methods), "\n"))
		if err != nil {
			return err
		}

		_, err = w.WriteString(strings.Join(sortedUnique(functions), "\n"))
		if err != nil {
			return err
		}

		return err
	}
}

func sortedUnique(in []string) (out []string) {
	m := make(map[string]bool)
	for i := range in {
		if _, ok := m[in[i]]; ok {
			continue
		}
		out = append(out, in[i])
		m[in[i]] = true
	}
	sort.Strings(out)
	return out
}
