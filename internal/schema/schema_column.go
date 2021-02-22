package schema

import "fmt"

func (c *Column) ExportedGoName() string {
	if c.GoName != "" {
		return pascal(c.GoName)
	}

	return pascal(c.name)
}

func (c *Column) SetName(in string) {
	c.name = in
}

func (c *Column) GoTypeString() string {
	s := c.Datatype.GoTypeString()

	if c.Unsigned == false && c.Datatype.IsSignable() && c.Datatype.HasGoUnsigned() {
		s = fmt.Sprintf("u%s", s)
	}

	return s
}

func (c *Column) RequiredImport() string {
	if c.Datatype.IsTime() {
		return `"time"`
	}

	return ""
}
