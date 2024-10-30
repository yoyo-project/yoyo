package schema

import (
	"fmt"
	"strings"

	"github.com/yoyo-project/yoyo/internal/datatype"
)

// ExportedGoName returns the string that will be used for naming Exported types, functions, etc in generated Go code
func (c *Column) ExportedGoName() string {
	if c.GoName != "" {
		return pascal(c.GoName)
	}

	return pascal(c.Name)
}

// GoTypeString returns the string keyword of the column type's corresponding Go type
func (c *Column) GoTypeString() string {
	var s string
	if c.Nullable {
		s = c.Datatype.GoNullableTypeString()
	} else {
		s = c.Datatype.GoTypeString()
		if c.Unsigned && c.Datatype.IsSignable() && c.Datatype.HasGoUnsigned() {
			s = fmt.Sprintf("u%s", s)
		}
	}

	return s
}

// RequiredImport returns any packages that need to be imported to support the Go type of a column in generated  Go code
func (c *Column) RequiredImport(nullPath string) string {
	if c.Datatype.IsTime() && !c.Nullable {
		return `"time"`
	}

	if c.Datatype == datatype.UUID {
		return `"github.com/google/uuid"`
	}

	if c.Nullable && strings.HasPrefix(c.Datatype.GoNullableTypeString(), "nullable") {
		return `"` + nullPath + `"`
	}

	return ""
}
