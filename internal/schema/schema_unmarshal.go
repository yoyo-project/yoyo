package schema

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
	"strings"
)

var invalidNameChars = regexp.MustCompile("[^a-zA-Z\\d_-]")
var scaleRemover = regexp.MustCompile("[^\\d,]")

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

	t.Indices = t2.Indices
	t.References = t2.References
	t.GoName = t2.GoName

	return t.validate()
}

func (c *Column) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c2 := new(struct {
		Datatype      string `yaml:"type"`
		Unsigned      bool
		Nullable      bool
		Default       *string
		Charset       string
		Collation     string
		PrimaryKey    bool   `yaml:"primary_key"`
		AutoIncrement bool   `yaml:"auto_increment"`
		GoName        string `yaml:"go_name"`
	})

	err := unmarshal(c2)
	if err != nil {
		return err
	}

	ps := scaleRemover.ReplaceAllString(c2.Datatype, "")
	ss := strings.Split(ps, ",")
	for i, s := range ss {
		if s == "" {
			continue
		}
		switch i {
		case 0:
			c.Precision, err = strconv.Atoi(s)
		case 1:
			c.Scale, err = strconv.Atoi(s)
		}
		if err != nil {
			return err
		}
	}

	err = yaml.Unmarshal([]byte(c2.Datatype), &c.Datatype)
	if err != nil {
		return err
	}
	c.Unsigned = c2.Unsigned
	c.Nullable = c2.Nullable
	c.Default = c2.Default
	c.Charset = c2.Charset
	c.Collation = c2.Collation
	c.PrimaryKey = c2.PrimaryKey
	c.AutoIncrement = c2.AutoIncrement
	c.GoName = c2.GoName

	return nil
}

// UnmarshalYAML provides an implementation for yaml/v2.Unmarshaler to parse a Reference definition
func (r *Reference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	r2 := new(struct {
		HasOne      bool `yaml:"has_one"`
		HasMany     bool `yaml:"has_many"`
		Required    bool
		ColumnNames []string `yaml:"column_names"`
		OnDelete    string   `yaml:"on_delete"`
		OnUpdate    string   `yaml:"on_update"`
	})

	err := unmarshal(r2)
	if err != nil {
		return err
	}

	r.HasOne = r2.HasOne
	r.HasMany = r2.HasMany
	r.Required = r2.Required
	r.ColumnNames = r2.ColumnNames
	r.OnUpdate = strings.TrimSpace(strings.ToUpper(r2.OnUpdate))
	r.OnDelete = strings.TrimSpace(strings.ToUpper(r2.OnDelete))

	return r.validate()
}
