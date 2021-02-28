package migration

import (
	"errors"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/schema"
)

const (
	Mock = "mockDialect"
)

type mockDialect struct{}

func (a *mockDialect) DataSourceName(host, username, schema, password, port string) string {
	return ""
}

func (a *mockDialect) AddReference(table string, rt schema.Table, i schema.Reference) string {
	return ""
}

func (a *mockDialect) ListReferences(table string) (string, interface{}) {
	return "", nil
}

func (*mockDialect) TypeString(dt datatype.Datatype) (string, error) {
	s := dt.String()

	if s == "NONE" {
		return "", errors.New("invalid type")
	}

	return s, nil
}

func (a *mockDialect) CreateTable(table string, t schema.Table) string {
	return ""
}

func (a *mockDialect) AddIndex(table, index string, i schema.Index) string {
	return ""
}

func (a *mockDialect) AddColumn(table, column string, c schema.Column) string {
	return ""
}

func (a *mockDialect) ListTables() (string, interface{}) {
	return "", 0
}

func (a *mockDialect) ListIndexes(table string) (string, interface{}) {
	return "", 0
}

func (a *mockDialect) ListColumns(table string) (string, interface{}) {
	return "", 0
}
