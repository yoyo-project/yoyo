package repository

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"
	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewEntityRepositoryGenerator(packageName string, adapter Adapter, reposPath string, packagePath Finder) EntityGenerator {
	return func(t schema.Table, w io.StringWriter) (err error) {
		var pkNames, cNames, scanFields, inFields, pkFields, colAssignments []string
		for columnName, col := range t.Columns {
			if !col.PrimaryKey {
				cNames = append(cNames, columnName)
			} else {
				pkFields = append(pkFields, strings.ReplaceAll(template.PKFieldTemplate, template.FieldName, col.ExportedGoName()))
				pkNames = append(pkNames, columnName)
			}
			scanFields = append(scanFields, fmt.Sprintf("&ent.%s", col.ExportedGoName()))
			inFields = append(inFields, fmt.Sprintf("in.%s", col.ExportedGoName()))
		}

		var queryImportPath string
		queryImportPath, err = packagePath(fmt.Sprintf("%s/query/%s", reposPath, t.QueryPackageName()))
		if err != nil {
			return fmt.Errorf("unable to generate repository: %w", err)
		}

		var pkCapture, pkCapTemplate string
		pkReplacer := strings.NewReplacer()

		switch len(t.PKColumns()) {
		case 0:
			// Do nothing
		case 1:
			col := t.PKColumns()[0]
			pkCapTemplate = template.SinglePKCaptureTemplate
			pkReplacer = strings.NewReplacer(
				template.FieldName,
				col.ExportedGoName(),
				template.Type,
				col.GoTypeString(),
			)
		default:
			pkCapTemplate = template.MultiPKCaptureTemplate
			pkReplacer = strings.NewReplacer()
		}

		pkCapture = pkReplacer.Replace(pkCapTemplate)

		pkQueryReplacer := strings.NewReplacer(
			template.QueryPackageName,
			t.QueryPackageName(),
			template.PKFields,
			strings.Join(pkFields, "\n		"),
		)

		pkQuery := pkQueryReplacer.Replace(template.PKQueryTemplate)

		preparedStatementPlaceholders := adapter.PreparedStatementPlaceholders(len(cNames))
		for i, colName := range cNames {
			colAssignments = append(colAssignments, fmt.Sprintf("%s = %s", colName, preparedStatementPlaceholders[i]))
		}

		r := strings.NewReplacer(
			template.PackageName,
			packageName,
			template.Imports,
			fmt.Sprintf(`"%s"`, queryImportPath),
			template.ScanFields,
			strings.Join(scanFields, ", "),
			template.InFields,
			strings.Join(inFields, ", "),
			template.EntityName,
			t.ExportedGoName(),
			template.TableName,
			t.TableName(),
			template.ColumnNames,
			strings.Join(cNames, ", "),
			template.StatementPlaceholders,
			strings.Join(preparedStatementPlaceholders, ", "),
			template.PKCapture,
			pkCapture,
			template.PKQuery,
			pkQuery,
			template.ColumnAssignments,
			strings.Join(colAssignments, ", "),
			template.QueryPackageName,
			t.QueryPackageName(),
		)

		_, err = w.WriteString(r.Replace(template.RepositoryFile))

		return err
	}
}
