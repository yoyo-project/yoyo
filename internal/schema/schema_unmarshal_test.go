package schema

import (
	"github.com/yoyo-project/yoyo/internal/datatype"
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

func TestDatabase_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yml     string
		wantDB  Database
		wantErr bool
	}{
		{
			name: "no tables",
			yml:  "dialect: mysql",
			wantDB: Database{
				Dialect: "mysql",
			},
		},
		{
			name: "with single table",
			yml: `
dialect: mysql
tables:
  primary:
    columns:
      id:
        type: int`,
			wantDB: Database{
				Dialect: "mysql",
				Tables:  []Table{{Name: "primary", Columns: []Column{{Name: "id", Datatype: datatype.Integer}}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := Database{}
			if err := yaml.Unmarshal([]byte(tt.yml), &db); err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(db, tt.wantDB) {
				t.Errorf("\nWant %#v,\n got %#v", tt.wantDB, db)
			}
		})
	}
}

func TestReference_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yml     string
		wantRef Reference
		wantErr bool
	}{
		{
			yml: "has_one: true",
			wantRef: Reference{
				HasOne: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reference{}
			if err := yaml.Unmarshal([]byte(tt.yml), &r); err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(r, tt.wantRef) {
				t.Errorf("Want %#v, got %#v", tt.wantRef, r)
			}
		})
	}
}

func TestTable_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yml     string
		want    Table
		wantErr bool
	}{
		{
			yml: `
columns:
  id:
    type: INT
    primary_key: true
  id2:
    type: INT
    primary_key: true`,
			want: Table{
				Columns: []Column{
					{
						Name:       "id",
						Datatype:   datatype.Integer,
						PrimaryKey: true,
					},
					{
						Name:       "id2",
						Datatype:   datatype.Integer,
						PrimaryKey: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Table{}
			if err := yaml.Unmarshal([]byte(tt.yml), &r); err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("\nWant %#v,\n got %#v", tt.want, r)
			}
		})
	}
}

func TestColumn_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yml     string
		want    Column
		wantErr bool
	}{
		{
			name: "int pk",
			yml: `
    type: INT
    primary_key: true`,
			want: Column{
				Datatype:   datatype.Integer,
				PrimaryKey: true,
			},
		},
		{
			name: "enum",
			yml: `type: enum("red", "blue")`,
			want: Column{
				Datatype:   datatype.Enum,
				Params: []string{"red", "blue"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Column{}
			if err := yaml.Unmarshal([]byte(tt.yml), &r); err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("\nWant %#v,\n got %#v", tt.want, r)
			}
		})
	}

}
