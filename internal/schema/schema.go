package schema

import (
	"github.com/yoyo-project/yoyo/internal/datatype"
)

// Database represents a single database/schema in a DBMS.
type Database struct {
	Dialect string
	Tables  map[string]Table
}

// Table represents a table in a database
type Table struct {
	Columns    map[string]Column
	Indices    map[string]Index
	References map[string]Reference // map[tableName]Reference
	GoName     string               `yaml:"go_name"`

	name string // used for internal yoyo logic
}

// Column represents a column in a table
type Column struct {
	Datatype      datatype.Datatype `yaml:"type"`
	Unsigned      bool
	Nullable      bool
	Default       *string
	Precision     int
	Scale         int
	Charset       string
	Collation     string
	PrimaryKey    bool   `yaml:"primary_key"`
	AutoIncrement bool   `yaml:"auto_increment"`
	GoName        string `yaml:"go_name"`

	name string // used for internal yoyo logic
}

// Reference represents a relationship boetween tables.
// Not a SQL-native concept, more of an ORM-style design. Translates to foreign keys and constraints in SQL
type Reference struct {
	HasOne      bool `yaml:"has_one"`
	HasMany     bool `yaml:"has_many"`
	Required    bool
	ColumnNames []string `yaml:"column_names"`
	OnDelete    string   `yaml:"on_delete"`
	OnUpdate    string   `yaml:"on_update"`
}

// Index represents a simple index on a column or columns
type Index struct {
	Columns []string
	Unique  bool
}
