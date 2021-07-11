package sqlite

import (
	"fmt"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/schema"
)

func (*adapter) SupportsDatatype(dt datatype.Datatype) bool {
	switch dt {
	case datatype.Boolean:
		return false
	}

	return true
}

func (*adapter) ValidateTable(t schema.Table) error {
	var ai bool
	for _, c := range t.PKColumns() {
		if c.AutoIncrement {
			ai = true
		}
	}
	if ai && len(t.PKColumns()) != 1 {
		return fmt.Errorf("cannot use AutoIncrement on a table with a compound primary key in sqlite")
	}
	return nil
}

func (*adapter) SupportsAutoIncrement() bool {
	return true
}
