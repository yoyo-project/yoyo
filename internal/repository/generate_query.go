package repository

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"
	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewQueryFileGenerator(reposPath string, findPackagePath Finder) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		var methods, functions, imports []string
		for _, c := range t.Columns {
			ms, fs, is := template.GenerateQueryLogic(c.Name, c)
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
