package schema

import (
	"fmt"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"strconv"
	"strings"
)

const (
	// Currently, only cascade is supported by Yoyo
	setCascade = "CASCADE"

	actions = ":" + setCascade + ":"

	asc    = "ASC"
	desc   = "DESC"
	orders = ":" + asc + ":" + desc + ":"

	nameCharacterLimit = 30
	dialects           = ":mysql:postgresql:"
)

func validateName(name string) error {
	invalid := invalidNameChars.Match([]byte(name))
	if invalid {
		return fmt.Errorf("invalid characters in name: %s", name)
	}
	if len(name) > nameCharacterLimit {
		return fmt.Errorf("name '%s' is too long. Must be %d characters or less", name, nameCharacterLimit)
	}
	return nil
}

func (c *Column) validate() error {
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

	if c.Datatype.RequiresScale() && c.Scale < 1 {
		return fmt.Errorf("datatype '%s' requires a scale value", c.Datatype)
	}

	return nil
}

func (i *Index) validate() error {
	if len(i.Columns) == 0 {
		return fmt.Errorf("index must have at least one column")
	}

	return nil
}

func (r *Reference) validate() error {
	if r.HasMany == r.HasOne {
		return fmt.Errorf("reference must be either HasOne or HasMany")
	}

	if !strings.Contains(actions, r.OnUpdate) {
		return fmt.Errorf("unknown action for OnUpdate: '%s'", r.OnUpdate)
	}
	if !strings.Contains(actions, r.OnDelete) {
		return fmt.Errorf("unknown action for OnUpdate: '%s'", r.OnUpdate)
	}

	if len(r.ColumnNames) > 0 && len(r.ColumnName) > 0 {
		return fmt.Errorf("cannot set both column_name and column_names for reference")
	}

	return nil
}

func (t *Table) validate() error {
	if len(t.Columns) == 0 {
		return fmt.Errorf("must have at least one column")
	}

	var err error

	for name, col := range t.Columns {
		if err = validateName(name); err != nil {
			return err
		}
		if err = col.validate(); err != nil {
			return fmt.Errorf("column '%s' validation error: %w", name, err)
		}
	}

	for name, i := range t.Indices {
		if err = validateName(name); err != nil {
			return err
		}
		if err = i.validate(); err != nil {
			return fmt.Errorf("index '%s' validation error: %w", name, err)
		}

		for _, icn := range i.Columns {
			var colExists bool
			for cn, _ := range t.Columns {
				if icn == cn {
					colExists = true
					break
				}
			}
			if !colExists {
				return fmt.Errorf("index '%s' validation error: column '%s' referenced but doesn't exist in table def", name, icn)
			}
		}
	}

	for ft, r := range t.References {
		if err = r.validate(); err != nil {
			return fmt.Errorf("%w on reference for table %s", err, ft)
		}
	}

	return nil
}

func (db *Database) validate() error {
	if !strings.Contains(dialects, fmt.Sprintf(":%s:", db.Dialect)) {
		return fmt.Errorf("unknown dialect: %s", db.Dialect)
	}

	var err error
	for tName, t := range db.Tables {
		if err = validateName(tName); err != nil {
			return err
		}

		if err = t.validate(); err != nil {
			return fmt.Errorf("%w for table `%s`", err, tName)
		}
		var ai bool

		for cName, c := range t.Columns {
			if c.AutoIncrement {
				if db.Dialect != dialect.MySQL {
					return fmt.Errorf("only mysql is supported for autoIncrement (column `%s`.`%s`)", tName, cName)
				}
				if ai {
					return fmt.Errorf("only one autoincrement column allowed")
				}
				ai = true
			}
		}

		for rName, r := range t.References {
			ft, ok := db.Tables[rName]

			if !ok {
				return fmt.Errorf("reference table `%s` does not exist in schema yaml", rName)
			}

			if r.HasMany {
				ft, rName = t, tName
			}

			var hasPK bool

			for _, c := range ft.Columns {
				if c.PrimaryKey {
					hasPK = true
				}
			}

			if !hasPK {
				return fmt.Errorf("table `%s` has no primary key but needs it for a defined reference", rName)
			}
		}
	}

	return nil
}
