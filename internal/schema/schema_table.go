package schema

import (
	"unicode"
)

// ExportedGoName returns a string with a name to use to represent the table in Exported go types or variables.
// If a GoName is explicitly set already, the returned value will be that forced into PascalCase. If not, it will be
// the default name forced into PascalCase.
func (t *Table) ExportedGoName() string {
	if t.GoName != "" {
		return pascal(t.GoName)
	}

	return pascal(t.Name)
}

// GetColumn returns a column matching the given name if present. If a matching column is found, the returned bool is true.
// If a matching column is not found, the returned bool is false.
func (t *Table) GetColumn(name string) (Column, bool) {
	for _, col := range t.Columns {
		if col.Name == name {
			return col, true
		}
	}
	return Column{}, false
}

// QueryPackageName returns a string to use for this table's query package
func (t *Table) QueryPackageName() string {
	if t.GoName != "" {
		return pascal(t.GoName)
	}

	name := pascal(t.Name)
	name = string(append([]byte{byte(unicode.ToLower(rune(name[0])))}, name[1:]...))
	return name
}

// PKColNames returns a list of column names that represent this table's Primary Key.
// This is usually a slice with a single value, but may be empty if there is no PK, and may be several if it's a
// compound PK
func (t *Table) PKColNames() (cols []string) {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			cols = append(cols, col.Name)
		}
	}

	return cols
}

// PKColumns returns a list of Columns that represent this table's Primary Key.
// This is usually a slice with a single value, but may be empty if there is no PK, and may be several if it's a
// compound PK
func (t *Table) PKColumns() (cols []Column) {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			cols = append(cols, col)
		}
	}

	return cols
}
