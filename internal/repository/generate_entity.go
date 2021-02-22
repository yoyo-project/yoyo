package repository

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"
	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewEntityGenerator(packageName string, ts map[string]schema.Table) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		var fields, referenceFields, scanFields, imports []string
		for _, c := range t.Columns {
			fields = append(fields, fmt.Sprintf("%s %s", c.ExportedGoName(), c.GoTypeString()))
			scanFields = append(scanFields, fmt.Sprintf("&e.%s", c.ExportedGoName()))
			if imp := c.RequiredImport(); imp != "" {
				imports = append(imports, imp)
			}
		}

		for rn, r := range t.References {
			if r.HasOne {
				ft := ts[rn]
				for _, cn := range ft.PKColNames() {
					c := ft.Columns[cn]
					referenceFields = append(referenceFields, fmt.Sprintf("%s%s %s", ft.ExportedGoName(), c.ExportedGoName(), c.GoTypeString()))
				}
			}
		}

		for _, t2 := range ts {
			for rn, r := range t2.References {
				if r.HasMany && rn == t.TableName() {
					for _, cn := range t2.PKColNames() {
						c := t2.Columns[cn]
						referenceFields = append(referenceFields, fmt.Sprintf("%s %s", c.ExportedGoName(), c.GoTypeString()))
					}
				}
			}
		}

		r := strings.NewReplacer(
			template.PackageName,
			packageName,
			template.EntityFields,
			strings.Join(sortedUnique(fields), "\n	"),
			template.ScanFields,
			strings.Join(sortedUnique(scanFields), ", "),
			template.Imports,
			strings.Join(sortedUnique(imports), "\n	"),
			template.EntityName,
			t.ExportedGoName(),
			template.ReferenceFields,
			strings.Join(sortedUnique(referenceFields), "\n	"),
		)

		_, err := w.WriteString(r.Replace(template.EntityFile))
		return err
	}
}
