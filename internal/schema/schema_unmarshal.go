package schema

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

var invalidNameChars = regexp.MustCompile("[^a-zA-Z\\d_-]")
var paramIsolator = regexp.MustCompile("(^.*?(\\(|$)|[\"\\)\\s])")

func (db *Database) UnmarshalYAML(value *yaml.Node) (err error) {
	for i, n := range value.Content {
		switch n.Value {
		case "dialect":
			err = value.Content[i+1].Decode(&db.Dialect)
		case "tables":
			tabsNode := value.Content[i+1]
			for ti, tn := range tabsNode.Content {
				if tn.Tag == "!!str" {
					t := Table{}
					t.Name = tn.Value
					err = tabsNode.Content[ti+1].Decode(&t)
					if err != nil {
						return fmt.Errorf("unable to unmarshal database: %w", err)
					}
					db.Tables = append(db.Tables, t)
				}
			}
		}

		if err != nil {
			return fmt.Errorf("unable to unmarshal database: %w", err)
		}
	}

	return db.validate()
}

func (t *Table) UnmarshalYAML(value *yaml.Node) (err error) {
	for i, n := range value.Content {
		switch n.Value {
		case "columns":
			colsNode := value.Content[i+1]
			for ci, cn := range colsNode.Content {
				if cn.Tag == "!!str" {
					c := Column{}
					c.Name = cn.Value
					err = colsNode.Content[ci+1].Decode(&c)
					if err != nil {
						return fmt.Errorf("unable to unmarshal table: %w", err)
					}
					t.Columns = append(t.Columns, c)
				}
			}
		case "indices", "indexes":
			colsNode := value.Content[i+1]
			for ci, in := range colsNode.Content {
				if in.Tag == "!!str" {
					ix := Index{}
					ix.Name = in.Value
					err = colsNode.Content[ci+1].Decode(&ix)
					if err != nil {
						return fmt.Errorf("unable to unmarshal table: %w", err)
					}
					t.Indices = append(t.Indices, ix)
				}
			}
		case "references":
			colsNode := value.Content[i+1]
			for ri, rn := range colsNode.Content {
				if rn.Tag == "!!str" {
					r := Reference{}
					r.TableName = rn.Value
					err = colsNode.Content[ri+1].Decode(&r)
					if err != nil {
						return fmt.Errorf("unable to unmarshal table: %w", err)
					}
					t.References = append(t.References, r)
				}
			}
		case "go_name":
			err = value.Content[i+1].Decode(&t.GoName)
		}

		if err != nil {
			return fmt.Errorf("unable to unmarshal table: %w", err)
		}
	}

	return t.validate()
}

func (c *Column) UnmarshalYAML(value *yaml.Node) (err error) {
	for i, n := range value.Content {
		switch n.Value {
		case "type", "datatype":
			err = value.Content[i+1].Decode(&c.Datatype)
			if err != nil {
				return fmt.Errorf("unable to unmarshal column: %w", err)
			}
			var str string
			err = value.Content[i+1].Decode(&str)
			if err != nil {
				return fmt.Errorf("unable to unmarshal column: %w", err)
			}
			ps := paramIsolator.ReplaceAllString(str, "")
			if len(ps) > 0 {
				c.Params = strings.Split(ps, ",")
			}
		case "unsigned":
			err = value.Content[i+1].Decode(&c.Unsigned)
		case "nullable":
			err = value.Content[i+1].Decode(&c.Nullable)
		case "default":
			err = value.Content[i+1].Decode(&c.Default)
		case "charset":
			err = value.Content[i+1].Decode(&c.Charset)
		case "collation":
			err = value.Content[i+1].Decode(&c.Collation)
		case "primary", "primary_key":
			err = value.Content[i+1].Decode(&c.PrimaryKey)
		case "auto_increment":
			err = value.Content[i+1].Decode(&c.AutoIncrement)
		case "go_name":
			err = value.Content[i+1].Decode(&c.GoName)
		}

		if err != nil {
			return fmt.Errorf("unable to unmarshal column: %w", err)
		}
	}

	return c.validate()
}

func (r *Reference) UnmarshalYAML(value *yaml.Node) (err error) {
	for i, n := range value.Content {
		switch n.Value {
		case "has_one":
			err = value.Content[i+1].Decode(&r.HasOne)
		case "has_many":
			err = value.Content[i+1].Decode(&r.HasMany)
		case "required":
			err = value.Content[i+1].Decode(&r.Required)
		case "column_names", "columns":
			err = value.Content[i+1].Decode(r.ColumnNames)
		case "on_delete":
			err = value.Content[i+1].Decode(&r.OnDelete)
		case "on_update":
			err = value.Content[i+1].Decode(&r.OnUpdate)
		}

		if err != nil {
			return fmt.Errorf("unable to unmarshal column: %w", err)
		}
	}

	return r.validate()
}
