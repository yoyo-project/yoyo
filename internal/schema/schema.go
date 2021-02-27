package schema

import (
	"github.com/yoyo-project/yoyo/internal/datatype"
)

// Database represents a single database/schema in a DBMS.
type Database struct {
	Dialect string
	Tables  []Table
}

// Table represents a table in a database
type Table struct {
	Name       string
	GoName     string
	Columns    []Column
	Indices    []Index
	References []Reference
}

// Column represents a column in a table
type Column struct {
	Name          string
	GoName        string
	Datatype      datatype.Datatype
	Unsigned      bool
	Nullable      bool
	Default       *string
	Params        []string
	Charset       string
	Collation     string
	PrimaryKey    bool
	AutoIncrement bool
}

// Reference represents a relationship boetween tables.
// Not a SQL-native concept, more of an ORM-style design. Translates to foreign keys and constraints in SQL
type Reference struct {
	TableName   string
	HasOne      bool
	HasMany     bool
	Required    bool
	ColumnNames []string
	OnDelete    string
	OnUpdate    string
}

// Index represents a simple index on a column or columns
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}
