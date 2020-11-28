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
	}
}

type migrator struct {
	base.Base
}

func (d *migrator) TypeString(dt datatype.Datatype) (s string, err error) {
	if dt&datatype.MySQL != datatype.MySQL {
		return "", errors.New("unsupported datatype")
	}
	switch dt {
	case datatype.Integer:
		s = "INT"
	default:
		s, err = d.Base.TypeString(dt)
	}
	return s, err
}

// CreateTable generates a query to create a given table.
func (d *migrator) CreateTable(table string, t schema.Table) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("CREATE TABLE `%s` (\n", table))

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
		sb.WriteString(d.generateColumn(colName, c))
		if c.PrimaryKey {
			pks = append(pks, colName)
		}
	}

	if len(pks) > 0 {
		sb.WriteString(fmt.Sprintf("\nPRIMARY KEY (%s)", strings.Join(pks, ",")))
	}

	sb.WriteString("\n);")

	return sb.String()
}

func (d *migrator) AddColumn(table, column string, c schema.Column) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table, d.generateColumn(column, c))
}

// AddIndex returns a string query which adds the specified index to a table
func (d *migrator) AddIndex(table, index string, i schema.Index) string {
	var indexType string

	switch {
	case i.Unique:
		indexType = " UNIQUE INDEX"
	default:
		indexType = " INDEX"
	}

	cols := strings.Builder{}
	firstCol := true
	for _, col := range i.Columns {
		if !firstCol {
			cols.WriteRune(',')
		}
		firstCol = false
		cols.WriteString(fmt.Sprintf("%s", col))
	}

	return fmt.Sprintf("ALTER TABLE `%s` ADD %s `%s` (%s);\n", table, indexType, index, cols.String())
}

func (d *migrator) AddReference(table, referencedTable string, db schema.Database, r schema.Reference) (string, error) {
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

		sw.WriteString(d.AddColumn(table, fkname, col))
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

func (d *migrator) generateColumn(name string, c schema.Column) string {
	sb := strings.Builder{}

	if c.Datatype.RequiresScale() {
		sb.WriteString(fmt.Sprintf("`%s` %s(%d, %d)", name, c.Datatype, c.Scale, c.Precision))
	} else {
		if c.Scale > 0 {
			if c.Precision > 0 {
				sb.WriteString(fmt.Sprintf("`%s` %s(%d, %d)", name, c.Datatype, c.Scale, c.Precision))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` %s(%d)", name, c.Datatype, c.Scale))
			}
		} else {
			sb.WriteString(fmt.Sprintf("`%s` %s", name, c.Datatype))
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
			sb.WriteString(fmt.Sprintf(`"%s" `, *c.Default))
		} else {
			sb.WriteString(fmt.Sprintf("%s ", *c.Default))
		}
	} else if c.Nullable {
		sb.WriteString(` DEFAULT NULL`)
	}

	if !c.Nullable {
		sb.WriteString(" NOT")
	}
	sb.WriteString(" NULL")

	return sb.String()
}
