package schema

import (
	"github.com/dotvezz/yoyo/internal/datatype"
)

type Database struct {
	Dialect string
	Tables  map[string]Table
}

type Table struct {
	Columns    map[string]Column
	Indices    map[string]Index
	References map[string]Reference // map[tableName]Reference
}

//TODO: Handle serial/auto-increment
type Column struct {
	Datatype      datatype.Datatype `yaml:"type"`
	Unsigned      bool
	Nullable      bool
	Default       *string
	Precision     int
	Scale         int
	Charset       string
	Collation     string
	PrimaryKey    bool `yaml:"primary_key"`
	AutoIncrement bool `yaml:"auto_increment"`
}

type Reference struct {
	HasOne      bool
	HasMany     bool
	Optional    bool
	ColumnNames []string
	OnDelete    string
	OnUpdate    string
}

type Index struct {
	Columns []string
	Unique  bool
}
