package validation

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/schema"
)

func ValidateDatabase(db schema.Database) (err error) {
	var validator Adapter
	validator, err = LoadValidator(db.Dialect)
	if err != nil {
		return fmt.Errorf("unable to load database validator: %w", err)
	}

	for _, t := range db.Tables {
		err = validator.ValidateTable(t)
		for _, c := range t.Columns {
			if !validator.SupportsDatatype(c.Datatype) {
				return fmt.Errorf("%s does not support datatype `%s` on `%s`.`%s`", db.Dialect, c.Datatype, t.Name, c.Name)
			}
			if c.AutoIncrement && !validator.SupportsAutoIncrement() {
				return fmt.Errorf("%s does not support AutoIncrement on `%s`.`%s`", db.Dialect, t.Name, c.Name)
			}
		}
	}

	return nil
}
