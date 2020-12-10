package migration

import (
	"errors"
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
)

const (
	Mock = "mockDialect"
)

type mockDialect struct{}

func (d *mockDialect) DataSourceName(host, username, schema, password, port string) string {
	return ""
}

func (d *mockDialect) AddReference(table, referencedTable string, rt schema.Table, i schema.Reference) string {
	return ""
}

func (d *mockDialect) ListReferences(table string) (string, interface{}) {
	return "", nil
}

func (*mockDialect) TypeString(dt datatype.Datatype) (string, error) {
	s := dt.String()

	if s == "NONE" {
		return "", errors.New("invalid type")
	}

	return s, nil
}

func (d *mockDialect) CreateTable(table string, t schema.Table) string {
	return ""
}

func (d *mockDialect) AddIndex(table, index string, i schema.Index) string {
	return ""
}

func (d *mockDialect) AddColumn(table, column string, c schema.Column) string {
	return ""
}

func (d *mockDialect) ListTables() (string, interface{}) {
	return "", 0
}

func (d *mockDialect) ListIndexes(table string) (string, interface{}) {
	return "", 0
}

func (d *mockDialect) ListColumns(table string) (string, interface{}) {
	return "", 0
}
