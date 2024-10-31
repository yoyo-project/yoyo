package mysql

import (
	"fmt"
	"strings"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/schema"
)

// TypeString returns the string representation of a given datatype.Datatype for MySQL
// An error will be returned if the datatype.Datatype is invalid or not supported by MySQL
func (a adapter) TypeString(dt datatype.Datatype) (s string, err error) {
	switch dt {
	case datatype.Integer:
		s = "INT"
	case datatype.UUID:
		s = "VARCHAR(36)"
	default:
		s, err = a.Base.TypeString(dt)
	}
	if err == nil && !a.SupportsDatatype(dt) {
		err = fmt.Errorf("unsupported datatype")
	}
	return s, err
}

// CreateTable returns a query string that create a given table.
func (a adapter) CreateTable(tName string, t schema.Table) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("CREATE TABLE `%s` (\n", tName))

	var (
		first = true
		pks   []string
	)
	for _, c := range t.Columns {
		if !first {
			sb.WriteString("\n")
		} else {
			first = false
		}
		sb.WriteString("    ")
		sb.WriteString(a.generateColumn(c.Name, c))
		sb.WriteRune(',')
		if c.PrimaryKey {
			pks = append(pks, c.Name)
		}
	}

	if len(pks) > 0 {
		sb.WriteString(fmt.Sprintf("\n    PRIMARY KEY (`%s`)", strings.Join(pks, ",")))
	}

	sb.WriteString("\n);")

	return sb.String()
}

// AddColumn returns a string query which adds a column to an existing table
func (a adapter) AddColumn(tName, cName string, c schema.Column) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", tName, a.generateColumn(cName, c))
}

// AddIndex returns a string query which adds the specified index to an existing table
func (a adapter) AddIndex(tName, iName string, i schema.Index) string {
	indexType := "INDEX"
	if i.Unique {
		indexType = "UNIQUE INDEX"
	}
	return fmt.Sprintf("ALTER TABLE `%s` ADD %s `%s` (`%s`);", tName, indexType, iName, strings.Join(i.Columns, "`, `"))
}

// AddReference returns a query string that adds columns and foreign keys for the given table, foreign table, and schema.Reference
func (a adapter) AddReference(tName string, fTable schema.Table, r schema.Reference) string {
	var (
		fCols  = fTable.PKColNames()
		lCols  = r.ColNames(fTable)
		sw     = strings.Builder{}
		ftName = fTable.Name
	)

	for i, lColName := range lCols {
		fCol, _ := fTable.GetColumn(fCols[i])

		// Remove possibly invalid properties of fCol
		fCol.AutoIncrement = false
		fCol.PrimaryKey = false

		// Set properties of fCol to be correct for the current operation
		fCol.Nullable = !r.Required

		sw.WriteString(a.AddColumn(tName, lColName, fCol)) // use fCol because the column's definition needs to match
		sw.WriteRune('\n')
	}

	sw.WriteString(fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `reference_%s_%s_%s` FOREIGN KEY (`%s`) REFERENCES %s(`%s`)",
		tName, tName, ftName, strings.Join(fCols, "_"), strings.Join(lCols, "`, `"), ftName, strings.Join(fCols, "`, `")))

	if r.OnDelete != "" {
		sw.WriteString(fmt.Sprintf(" ON DELETE %s", r.OnDelete))
	}

	if r.OnUpdate != "" {
		sw.WriteString(fmt.Sprintf(" ON UPDATE %s", r.OnUpdate))
	}

	sw.WriteRune(';')

	return sw.String()
}

func (a adapter) generateColumn(cName string, c schema.Column) string {
	sb := strings.Builder{}
	ts, _ := a.TypeString(c.Datatype)

	if len(c.Params) > 0 {
		sb.WriteString(fmt.Sprintf("`%s` %s(%s)", cName, ts, strings.Join(c.Params, ", ")))
	} else {
		sb.WriteString(fmt.Sprintf("`%s` %s", cName, ts))
	}

	if c.Datatype.IsSignable() {
		if c.Unsigned {
			sb.WriteString(" UNSIGNED")
		} else {
			sb.WriteString(" SIGNED")
		}
	}

	if c.Default != nil {
		sb.WriteString(` DEFAULT `)
		if c.Datatype.IsString() {
			sb.WriteString(fmt.Sprintf(`"%s"`, *c.Default))
		} else {
			sb.WriteString(fmt.Sprintf("%s", *c.Default))
		}
	} else if c.Nullable {
		sb.WriteString(` DEFAULT NULL`)
	}

	if !c.Nullable {
		sb.WriteString(" NOT")
	}
	sb.WriteString(" NULL")

	if c.AutoIncrement {
		sb.WriteString(" AUTO_INCREMENT")
	}

	return sb.String()
}
