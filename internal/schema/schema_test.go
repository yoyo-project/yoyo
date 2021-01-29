package schema

import (
	"testing"

	"github.com/dotvezz/yoyo/internal/datatype"

	yaml "gopkg.in/yaml.v2"
)

var TestDatabase = Database{
	Dialect: "mysql",
	Tables: map[string]Table{
		"primary": {
			Columns: map[string]Column{
				"id": {
					Datatype: datatype.Integer,
				},
				"secondary_id": {
					Datatype: datatype.Integer,
				},
			},
			References: map[string]Reference{
				"secondary": {
					HasOne: true,
				},
			},
		},
		"secondary": {
			Columns: map[string]Column{
				"id": {
					Datatype: datatype.Integer,
				},
			},
		},
	},
}

const testYAML = `
---
dialect: mysql
tables:
  table1:
    columns:
      id:
        type: integer
        default: 0
  table2:
    columns: 
      id:
        type: integer
`

// TestUnmarshalYAML doesn't test any logic in Yoyo, but is really just a
//// dumb helper for confirming the correct yaml structure.
func TestUnmarshalYAML(t *testing.T) {
	r := Database{}
	err := yaml.Unmarshal([]byte(testYAML), &r)
	if err != nil {
		t.Error(err)
	}
}

// TestMarshal doesn't test any logic in Yoyo, but is really just a
// dumb helper for confirming the correct yaml structure.
func TestMarshal(t *testing.T) {
	_, err := yaml.Marshal(TestDatabase)
	if err != nil {
		t.Error(err)
	}

}
