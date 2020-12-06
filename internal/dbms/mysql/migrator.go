package mysql

import (
	"errors"
	"fmt"
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/dbms/base"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/schema"
	"strings"
)

func NewMigrator() *migrator {
	return &migrator{
		Base: base.Base{
			Dialect: dialect.MySQL,
		},
		validator: validator{},
	}
}

type migrator struct {
	base.Base
	validator
}

func (m *migrator) TypeString(dt datatype.Datatype) (s string, err error) {
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

// CreateTable generates a query to create a given table.
func (m *migrator) CreateTable(tName string, t schema.Table) string {
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

func (m *migrator) AddColumn(tName, cName string, c schema.Column) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", tName, m.generateColumn(cName, c))
}

// AddIndex returns a string query which adds the specified index to a table
func (m *migrator) AddIndex(tName, iName string, i schema.Index) string {
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

func (m *migrator) AddReference(table, referencedTable string, db schema.Database, r schema.Reference) (string, error) {
	sw := strings.Builder{}

	if r.HasMany { // swap the tables if it's a HasMany
		table, referencedTable = referencedTable, table
	}

	rt, ok := db.Tables[referencedTable]
	if !ok { // This should technically be caught by validation, but still
		return "", fmt.Errorf("referenced table `%s` does not exist in dbms definition", referencedTable)
	}
	var (
		fcols       []string
		fknames     []string
		fkname      string
		refColNames = r.ColumnNames
	)

	var hasPK bool

	for cname, col := range rt.Columns {
		if !col.PrimaryKey {
			continue
		}

		hasPK = true

		switch {
		case len(refColNames) > 0:
			fkname, refColNames = refColNames[0], refColNames[1:len(refColNames)]
		case len(r.ColumnName) > 0:
			fkname = r.ColumnName
		default:
			fkname = fmt.Sprintf("fk_%s_%s", referencedTable, cname)
		}

		fknames = append(fknames, fkname)
		fcols = append(fcols, cname)

		sw.WriteString(m.AddColumn(table, fkname, col))
		sw.WriteRune('\n')
	}

	if len(r.ColumnNames) > 0 && len(r.ColumnNames) != len(fcols) {
		return "", fmt.Errorf("cannot add reference from `%s` to `%s`: length of column_names does not match length of primary keys", table, referencedTable)
	}

	if !hasPK { // This should technically be caught by validation, but still
		return "", fmt.Errorf("referenced table `%s` does not have a primary key defined", referencedTable)
	}

	sw.WriteString(fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `reference_%s` FOREIGN KEY (`%s`) REFERENCES %s(%s)",
		table, referencedTable, strings.Join(fknames, ", "), referencedTable, strings.Join(fcols, ", "),
	))

	if r.OnDelete != "" {
		sw.WriteString(fmt.Sprintf(" ON DELETE %s", r.OnDelete))
	}

	if r.OnUpdate != "" {
		sw.WriteString(fmt.Sprintf(" ON UPDATE %s", r.OnUpdate))
	}

	return sw.String(), nil
}

func (m *migrator) generateColumn(cName string, c schema.Column) string {
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
