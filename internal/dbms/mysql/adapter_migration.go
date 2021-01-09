package mysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/dbms/base"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/schema"
)

// NewAdapter returns an implementation of migration.Dialect for MySQL
func NewAdapter() *adapter {
	return &adapter{
		Base: base.Base{
			Dialect: dialect.MySQL,
		},
		validator: validator{},
	}
}

// TypeString returns the string representation of a given datatype.Datatype for MySQL
// An error will be returned if the datatype.Datatype is invalid or not supported by MySQL
func (m *adapter) TypeString(dt datatype.Datatype) (s string, err error) {
	switch dt {
	case datatype.Integer:
		s = "INT"
	default:
		s, err = m.Base.TypeString(dt)
	}
	if err == nil && !m.validator.SupportsDatatype(dt) {
		err = errors.New("unsupported datatype")
	}
	return s, err
}

// CreateTable returns a query string that create a given table.
func (m *adapter) CreateTable(tName string, t schema.Table) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("CREATE TABLE `%s` (\n", tName))

	var (
		first = true
		pks   []string
	)
	for colName, c := range t.Columns {
		if !first {
			sb.WriteString(",\n")
		} else {
			first = false
		}
		sb.WriteString("    ")
		sb.WriteString(m.generateColumn(colName, c))
		if c.PrimaryKey {
			pks = append(pks, colName)
		}
	}

	if len(pks) > 0 {
		sb.WriteString(fmt.Sprintf("\n    PRIMARY KEY (`%s`)", strings.Join(pks, ",")))
	}

	sb.WriteString("\n);")

	return sb.String()
}

// AddColumn returns a string query which adds a column to an existing table
func (m *adapter) AddColumn(tName, cName string, c schema.Column) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", tName, m.generateColumn(cName, c))
}

// AddIndex returns a string query which adds the specified index to an existing table
func (m *adapter) AddIndex(tName, iName string, i schema.Index) string {
	var indexType string

	switch {
	case i.Unique:
		indexType = "UNIQUE INDEX"
	default:
		indexType = "INDEX"
	}

	cols := strings.Builder{}
	firstCol := true
	for _, col := range i.Columns {
		if !firstCol {
			cols.WriteString(", ")
		}
		firstCol = false
		cols.WriteString(fmt.Sprintf("`%s`", col))
	}

	return fmt.Sprintf("ALTER TABLE `%s` ADD %s `%s` (%s);", tName, indexType, iName, cols.String())
}

// AddReference returns a query string that adds columns and foreign keys for the given table, foreign table, and schema.Reference
func (m *adapter) AddReference(tName, ftName string, fTable schema.Table, r schema.Reference) string {
	var (
		fCols = fTable.PKColNames()
		lCols = r.ColNames(ftName, fTable)
		sw    = strings.Builder{}
	)

	for i := range lCols {
		fCol := fTable.Columns[fCols[i]]
		fCol.Nullable = !r.Required
		sw.WriteString(m.AddColumn(tName, lCols[i], fCol))
		sw.WriteRune('\n')
	}

	sw.WriteString(fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `reference_%s` FOREIGN KEY (`%s`) REFERENCES %s(`%s`)",
		tName, ftName, strings.Join(lCols, "`, `"), ftName, strings.Join(fCols, "`, `")))

	if r.OnDelete != "" {
		sw.WriteString(fmt.Sprintf(" ON DELETE %s", r.OnDelete))
	}

	if r.OnUpdate != "" {
		sw.WriteString(fmt.Sprintf(" ON UPDATE %s", r.OnUpdate))
	}

	sw.WriteRune(';')

	return sw.String()
}

func (m *adapter) generateColumn(cName string, c schema.Column) string {
	sb := strings.Builder{}
	ts, _ := m.TypeString(c.Datatype)

	if c.Datatype.RequiresScale() {
		sb.WriteString(fmt.Sprintf("`%s` %s(%d, %d)", cName, ts, c.Scale, c.Precision))
	} else {
		if c.Scale > 0 {
			if c.Precision > 0 {
				sb.WriteString(fmt.Sprintf("`%s` %s(%d, %d)", cName, ts, c.Scale, c.Precision))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` %s(%d)", cName, ts, c.Scale))
			}
		} else {
			sb.WriteString(fmt.Sprintf("`%s` %s", cName, ts))
		}
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
