package postgres

import (
	"fmt"
	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/schema"
	"strings"
)

// TypeString returns the string representation of a given datatype.Datatype for PostgreSQL
// An error will be returned if the datatype.Datatype is invalid or not supported by PostgreSQL
func (a adapter) TypeString(dt datatype.Datatype) (s string, err error) {
	if !a.SupportsDatatype(dt) {
		return "", fmt.Errorf("datatype %s is not supported in postgresql", dt)
	}
	switch dt {
	default:
		s, err = a.Base.TypeString(dt)
	}
	return s, err
}

// CreateTable generates a query to create a given table.
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
		sb.WriteString(fmt.Sprintf(`\nCONSTRAINT "%s-pk" PRIMARY KEY ("%s")`, tName, strings.Join(pks, `", "`)))
	}

	sb.WriteString("\n);")

	return sb.String()
}

// AddColumn generates a query that adds a column to an existing table
func (a adapter) AddColumn(tName, cName string, c schema.Column) string {
	return fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN %s;`, tName, a.generateColumn(cName, c))
}

// AddIndex returns a string query which adds the specified index to a table
func (a adapter) AddIndex(tName, iName string, i schema.Index) string {
	indexType := "INDEX"
	if i.Unique {
		indexType = "UNIQUE INDEX"
	}
	return fmt.Sprintf(`CREATE %s "%s" ON "%s" ("%s");`, indexType, tName, iName, strings.Join(i.Columns, `", "`))
}

// AddReference generates a query that adds columns and foreign keys for the given table, foreign table, and schema.Reference
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

	sw.WriteString(fmt.Sprintf(`ALTER TABLE "%s" ADD CONSTRAINT "reference_%s_%s_%s" FOREIGN KEY ("%s") REFERENCES "%s"("%s")`,
		tName, tName, ftName, strings.Join(fCols, "_"), strings.Join(lCols, `", "`), ftName, strings.Join(fCols, `", "`)))

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
		sb.WriteString(fmt.Sprintf(`"%s" %s(%s)`, cName, ts, strings.Join(c.Params, ", ")))
	} else {
		sb.WriteString(fmt.Sprintf(`"%s" %s`, cName, ts))
	}

	if c.Default != nil {
		sb.WriteString(` DEFAULT `)
		if c.Datatype.IsString() {
			sb.WriteString(fmt.Sprintf("'%s'", *c.Default))
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

	// TODO: Auto Increment / serial

	return sb.String()
}
