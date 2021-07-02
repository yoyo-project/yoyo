package schema

import "fmt"

// ExportedGoName returns the string that will be used for naming Exported types, functions, etc in generated Go code
func (c *Column) ExportedGoName() string {
	if c.GoName != "" {
		return pascal(c.GoName)
	}

	return pascal(c.Name)
}

// GoTypeString returns the string keyword of the column type's corresponding Go type
func (c *Column) GoTypeString() string {
	s := c.Datatype.GoTypeString()

	if c.Unsigned == false && c.Datatype.IsSignable() && c.Datatype.HasGoUnsigned() {
		s = fmt.Sprintf("u%s", s)
	}

	return s
}

// RequiredImport returns any packages that need to be imported to support the Go type of a column in generated  Go code
func (c *Column) RequiredImport() string {
	if c.Datatype.IsTime() {
		return `"time"`
	}

	return ""
}
