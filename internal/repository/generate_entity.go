package repository

import (
	"fmt"
	"io"
	"strings"

	"github.com/dotvezz/yoyo/internal/repository/template"
	"github.com/dotvezz/yoyo/internal/schema"
)

func NewEntityGenerator(ts map[string]schema.Table) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) error {
		var fields, referenceFields []string
		for _, c := range t.Columns {
			fields = append(fields, fmt.Sprintf("%s %s", c.ExportedGoName(), c.GoTypeString()))
		}

		for rn, r := range t.References {
			if r.HasOne {
				ft := ts[rn]
				for _, cn := range ft.PKColNames() {
					c := ft.Columns[cn]
					referenceFields = append(referenceFields, fmt.Sprintf("%s %s", c.ExportedGoName(), c.GoTypeString()))
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
			template.EntityFields,
			strings.Join(fields, "\n    "),
			template.EntityName,
			t.ExportedGoName(),
			template.ReferenceFields,
			strings.Join(referenceFields, "\n    "),
		)

		_, err := w.WriteString(r.Replace(template.EntityFile))
		return err
	}
}
