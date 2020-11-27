package base

import (
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
)

type Base struct {
	Dialect string
}

func (d *Base) TypeString(dt datatype.Datatype) (string, error) {
	s := dt.String()

	if s == "NONE" {
		return "", d.InvalidDatatype(dt)
	}

	return s, nil
}

func (d *Base) AddReference(table, reference string, db schema.Database, i schema.Reference) (string, error) {
	panic("implement me")
}

func (d *Base) CreateTable(table string, t schema.Table) string {
	panic("implement me")
}

func (d *Base) AddIndex(table, index string, i schema.Index) string {
	panic("implement me")
}

func (d *Base) AddColumn(table, column string, c schema.Column) string {
	panic("implement me")
}

func (d *Base) ListTables() ([]string, error) {
	panic("implement me")
}

func (d *Base) ListColumns(table string) ([]string, error) {
	panic("implement me")
}

func (d *Base) ListIndices(table string) ([]string, error) {
	panic("implement me")
}

func (d *Base) ListReferences(table string) ([]string, error) {
	panic("implement me")
}

func (d *Base) GetColumn(table, column string) (schema.Column, error) {
	panic("implement me")
}

func (d *Base) GetIndex(table, index string) (schema.Index, error) {
	panic("implement me")
}

func (d *Base) GetReference(table, reference string) (schema.Reference, error) {
	panic("implement me")
}
