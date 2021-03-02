package migration

import (
	"errors"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/schema"
)

const (
	Mock = "mockAdapter"
)

type mockAdapter struct{}

func (a *mockAdapter) DataSourceName(host, username, schema, password, port string) string {
	return ""
}

func (a *mockAdapter) AddReference(table string, rt schema.Table, i schema.Reference) string {
	return ""
}

func (a *mockAdapter) ListReferences(table string) (string, interface{}) {
	return "", nil
}

func (*mockAdapter) TypeString(dt datatype.Datatype) (string, error) {
	s := dt.String()

	if s == "NONE" {
		return "", errors.New("invalid type")
	}

	return s, nil
}

func (a *mockAdapter) CreateTable(table string, t schema.Table) string {
	return ""
}

func (a *mockAdapter) AddIndex(table, index string, i schema.Index) string {
	return ""
}

func (a *mockAdapter) AddColumn(table, column string, c schema.Column) string {
	return ""
}

func (a *mockAdapter) ListTables() (string, interface{}) {
	return "", 0
}

func (a *mockAdapter) ListIndexes(table string) (string, interface{}) {
	return "", 0
}

func (a *mockAdapter) ListColumns(table string) (string, interface{}) {
	return "", 0
}

type mockReverseAdapter struct {
}

func (m mockReverseAdapter) ListTables() ([]string, error) {
	panic("implement me")
}

func (m mockReverseAdapter) ListColumns(table string) ([]string, error) {
	panic("implement me")
}

func (m mockReverseAdapter) ListIndices(table string) ([]string, error) {
	panic("implement me")
}

func (m mockReverseAdapter) ListReferences(table string) ([]string, error) {
	panic("implement me")
}

func (m mockReverseAdapter) GetColumn(table, column string) (schema.Column, error) {
	panic("implement me")
}

func (m mockReverseAdapter) GetIndex(table, column string) (schema.Index, error) {
	panic("implement me")
}

func (m mockReverseAdapter) GetReference(table, column string) (schema.Reference, error) {
	panic("implement me")
}
