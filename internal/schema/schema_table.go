package schema

import (
	"unicode"
)

func (t *Table) ExportedGoName() string {
	if t.GoName != "" {
		return pascal(t.GoName)
	}

	return pascal(t.name)
}

func (t *Table) SetName(name string) {
	t.name = name
}

func (t *Table) TableName() string {
	return t.name
}

func (t *Table) QueryPackageName() string {
	if t.GoName != "" {
		return pascal(t.GoName)
	}

	name := pascal(t.name)
	name = string(append([]byte{byte(unicode.ToLower(rune(name[0])))}, name[1:]...))
	return name
}

func (t *Table) PKColNames() (cols []string) {
	for name, col := range t.Columns {
		if col.PrimaryKey {
			cols = append(cols, name)
		}
	}

	return cols
}

func (t *Table) PKColumns() (cols []Column) {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			cols = append(cols, col)
		}
	}

	return cols
}
