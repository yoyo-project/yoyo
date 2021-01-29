package schema

import (
	"regexp"
	"strings"
)

var invalidNameChars = regexp.MustCompile("[^a-zA-Z\\d_-]")

// UnmarshalYAML provides an implementation for yaml/v2.Unmarshaler to parse a Database definition
func (db *Database) UnmarshalYAML(unmarshal func(interface{}) error) error {
	db2 := new(struct {
		Tables  map[string]Table
		Dialect string
	})

	err := unmarshal(db2)
	if err != nil {
		return err
	}

	db.Dialect = strings.TrimSpace(strings.ToLower(db2.Dialect))

	db.Tables = make(map[string]Table)
	for tn, t := range db2.Tables {
		t.name = tn
		db.Tables[tn] = t
	}

	return db.validate()
}

func (t *Table) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t2 := new(struct {
		Columns    map[string]Column
		Indices    map[string]Index
		References map[string]Reference // map[tableName]Reference
		GoName     string               `yaml:"go_name"`
	})

	err := unmarshal(t2)
	if err != nil {
		return err
	}

	t.Columns = make(map[string]Column)
	for cn, c := range t2.Columns {
		c.name = cn
		t.Columns[cn] = c
	}

	return t.validate()
}

// UnmarshalYAML provides an implementation for yaml/v2.Unmarshaler to parse a Reference definition
func (r *Reference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	r2 := new(struct {
		HasOne   bool
		HasMany  bool
		OnDelete string
		OnUpdate string
	})

	err := unmarshal(r2)
	if err != nil {
		return err
	}

	r.OnUpdate = strings.TrimSpace(strings.ToUpper(r2.OnUpdate))
	r.OnDelete = strings.TrimSpace(strings.ToUpper(r2.OnDelete))
	r.HasMany = r2.HasMany
	r.HasOne = r2.HasOne

	return r.validate()
}
