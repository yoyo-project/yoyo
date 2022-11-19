package schema

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Currently, only cascade is supported by Yoyo
	setCascade = "CASCADE"

	actions = ":" + setCascade + ":"

	nameCharacterLimit = 30
)

func validateName(name string) error {
	invalid := disallowedNameChars.Match([]byte(name))
	if invalid {
		return fmt.Errorf("invalid characters in Name: %s", name)
	}
	if len(name) > nameCharacterLimit {
		return fmt.Errorf("Name '%s' is too long. Must be %d characters or less", name, nameCharacterLimit)
	}
	return nil
}

func (c *Column) validate() error {
	if err := validateName(c.Name); err != nil {
		return err
	}

	if c.Unsigned && !c.Datatype.IsSignable() {
		return fmt.Errorf("unsigned flag is set but '%s' does not use a signing bit", c.Datatype)
	}

	if c.AutoIncrement && !c.Datatype.IsInt() {
		return fmt.Errorf("auto_ncrement is set but '%s' is not incrementable (only ints allowed)", c.Datatype)
	}

	if c.AutoIncrement && c.Default != nil {
		return fmt.Errorf("auto_increment is set but the column has a default value")
	}

	if c.AutoIncrement && !c.PrimaryKey {
		return fmt.Errorf("auto_increment is set but the column is not primary key")
	}

	if c.Default != nil {
		if _, err := strconv.ParseFloat(*c.Default, 64); err != nil && c.Datatype.IsNumeric() {
			return fmt.Errorf("non-numeric default '%s' used for numeric type '%s'", *c.Default, c.Datatype)
		}
	}

	if c.Datatype.RequiresParams() && len(c.Params) == 0 {
		return fmt.Errorf("datatype '%s' requires at least one parameter", c.Datatype)
	}

	return nil
}

func (i *Index) validate() error {
	if err := validateName(i.Name); err != nil {
		return err
	}

	if len(i.Columns) == 0 {
		return fmt.Errorf("index must have at least one column")
	}

	return nil
}

func (r *Reference) validate() error {
	if err := validateName(r.GoName); err != nil {
		return err
	}
	if r.HasMany == r.HasOne {
		return fmt.Errorf("reference must be either HasOne or HasMany")
	}

	if !strings.Contains(actions, r.OnUpdate) {
		return fmt.Errorf("unknown action for OnUpdate: '%s'", r.OnUpdate)
	}
	if !strings.Contains(actions, r.OnDelete) {
		return fmt.Errorf("unknown action for OnUpdate: '%s'", r.OnUpdate)
	}

	return nil
}

func (t *Table) validate() (err error) {
	if err = validateName(t.Name); err != nil {
		return err
	}

	if len(t.Columns) == 0 {
		return fmt.Errorf("must have at least one column")
	}

	cNames := make(map[string]bool)
	for _, col := range t.Columns {
		if err = col.validate(); err != nil {
			return fmt.Errorf("column '%s' validation error: %w", col.Name, err)
		}

		if _, ok := cNames[col.Name]; ok {
			return fmt.Errorf("duplicate column name '%s'", col.Name)
		} else {
			cNames[col.Name] = true
		}
	}

	for _, i := range t.Indices {
		if err = i.validate(); err != nil {
			return fmt.Errorf("index '%s' validation error: %w", i.Name, err)
		}
		for _, icn := range i.Columns {
			var colExists bool
			for _, c := range t.Columns {
				if icn == c.Name {
					colExists = true
					break
				}
			}
			if !colExists {
				return fmt.Errorf("index '%s' validation error: column '%s' referenced but doesn't exist in table def", i.Name, icn)
			}
		}
	}

	for _, r := range t.References {
		if err = r.validate(); err != nil {
			return fmt.Errorf("%w on reference for table %s", err, r.TableName)
		}
	}

	return nil
}

func (db *Database) validate() (err error) {
	tNames := make(map[string]bool)
	for _, t := range db.Tables {
		if err = t.validate(); err != nil {
			return fmt.Errorf("%w for table `%s`", err, t.Name)
		}

		if _, ok := tNames[t.Name]; ok {
			return fmt.Errorf("duplicate table name '%s'", t.Name)
		} else {
			tNames[t.Name] = true
		}

		for _, r := range t.References {
			var ft Table
			rName := r.TableName
			for _, table := range db.Tables {
				if table.Name == rName {
					ft = table
				}
			}

			if ft.Name == "" {
				return fmt.Errorf("reference table `%s` does not exist in schema yaml", rName)
			}

			if r.HasMany {
				ft, rName = t, t.Name
			}

			pkCount := 0
			for _, c := range ft.Columns {
				if c.PrimaryKey {
					pkCount++
				}
			}

			if len(r.ColumnNames) > 0 && len(r.ColumnNames) != pkCount {
				return fmt.Errorf("cannot add reference from `%s` to `%s`: length of column_names does not match length of primary keys", t.Name, rName)
			}

			if pkCount == 0 {
				return fmt.Errorf("table `%s` has no primary key but needs it for a defined reference", rName)
			}
		}
	}

	return nil
}
