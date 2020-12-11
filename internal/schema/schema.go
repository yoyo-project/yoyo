package schema

import (
	"github.com/dotvezz/yoyo/internal/datatype"
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
	PrimaryKey    bool `yaml:"primary_key"`
	AutoIncrement bool `yaml:"auto_increment"`
}

// Reference represents a relationship boetween tables.
// Not a SQL-native concept, more of an ORM-style design. Translates to foreign keys and constraints in SQL
type Reference struct {
	HasOne      bool
	HasMany     bool
	Optional    bool
	ColumnNames []string
	OnDelete    string
	OnUpdate    string
}

// Index represents a simple index on a column or columns
type Index struct {
	Columns []string
	Unique  bool
}
