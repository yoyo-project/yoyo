package repository

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"
	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewEntityGenerator(packageName string, db schema.Database) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		var fields, referenceFields, scanFields, imports []string
		for _, c := range t.Columns {
			fields = append(fields, fmt.Sprintf("%s %s", c.ExportedGoName(), c.GoTypeString()))
			scanFields = append(scanFields, fmt.Sprintf("&e.%s", c.ExportedGoName()))
			if imp := c.RequiredImport(); imp != "" {
				imports = append(imports, imp)
			}
		}

		for _, r := range t.References {
			if r.HasOne {
				ft, _ := db.GetTable(r.TableName)
				for _, cn := range ft.PKColNames() {
					c, _ := ft.GetColumn(cn)
					referenceFields = append(referenceFields, fmt.Sprintf("%s%s %s", ft.ExportedGoName(), c.ExportedGoName(), c.GoTypeString()))
				}
			}
		}

		for _, t2 := range db.Tables {
			for _, r := range t2.References {
				if r.HasMany && r.TableName == t.Name {
					for _, c := range t2.PKColumns() {
						referenceFields = append(referenceFields, fmt.Sprintf("%s %s", c.ExportedGoName(), c.GoTypeString()))
					}
				}
			}
		}

		r := strings.NewReplacer(
			template.PackageName,
			packageName,
			template.EntityFields,
			strings.Join(fields, "\n	"),
			template.ScanFields,
			strings.Join(scanFields, ", "),
			template.Imports,
			strings.Join(sortedUnique(imports), "\n	"),
			template.EntityName,
			t.ExportedGoName(),
			template.ReferenceFields,
			strings.Join(referenceFields, "\n	"),
		)

		_, err := w.WriteString(r.Replace(template.EntityFile))
		return err
	}
}
