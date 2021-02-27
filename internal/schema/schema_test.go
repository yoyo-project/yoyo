package schema

import (
	"testing"

	"github.com/yoyo-project/yoyo/internal/datatype"

	yaml "gopkg.in/yaml.v3"
)

var TestDatabase = Database{
	Dialect: "mysql",
	Tables: []Table{
		{
			Name: "primary",
			Columns: []Column{
				{
					Name:     "id",
					Datatype: datatype.Integer,
				},
				{
					Name:     "id2",
					Datatype: datatype.Integer,
				},
			},
			References: []Reference{
				{
					TableName: "secondary",
					HasOne:    true,
				},
			},
		},
		{
			Name: "secondary",
			Columns: []Column{
				{
					Name:     "id",
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
      id2:
        type: integer
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
