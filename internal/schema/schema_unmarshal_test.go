package schema

import (
	"reflect"
	"testing"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"gopkg.in/yaml.v3"
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
		{
			name: "with invalid table",
			yml: `
dialect: mysql
tables:
  primary: no`,
			wantDB: Database{
				Dialect: "mysql",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := Database{}
			if err := yaml.Unmarshal([]byte(tt.yml), &db); err != nil && !tt.wantErr || err == nil && tt.wantErr {
				t.Errorf("Got error %s, want error %v", err, tt.wantErr)
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
			name: "has one",
			yml:  "has_one: true",
			wantRef: Reference{
				HasOne: true,
			},
		},
		{
			name: "has many",
			yml:  "has_many: true",
			wantRef: Reference{
				HasMany: true,
			},
		},
		{
			name: "required",
			yml:  "required: true\nhas_one: true",
			wantRef: Reference{
				HasOne:   true,
				Required: true,
			},
		},
		{
			name: "columns",
			yml:  "columns:\n  - col1\n  - col2\nhas_one: true",
			wantRef: Reference{
				HasOne:      true,
				ColumnNames: []string{"col1", "col2"},
			},
		},
		{
			name: "on delete",
			yml:  "has_one: true\non_delete: CASCADE",
			wantRef: Reference{
				HasOne:   true,
				OnDelete: "CASCADE",
			},
		},
		{
			name: "on update",
			yml:  "has_one: true\non_update: CASCADE",
			wantRef: Reference{
				HasOne:   true,
				OnUpdate: "CASCADE",
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
			name: "goname",
			yml: `
go_name: Table
columns:
  id:
    type: INT
    primary_key: true`,
			want: Table{
				GoName: "Table",
				Columns: []Column{
					{
						Name:       "id",
						Datatype:   datatype.Integer,
						PrimaryKey: true,
					},
				},
			},
		},
		{
			name: "one primary key",
			yml: `
columns:
  id:
    type: INT
    primary_key: true`,
			want: Table{
				Columns: []Column{
					{
						Name:       "id",
						Datatype:   datatype.Integer,
						PrimaryKey: true,
					},
				},
			},
		},
		{
			name: "two primary keys",
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
		{
			name: "index on single column",
			yml: `
columns:
  col:
    type: INT
indices:
  - name: idx_col
    columns:
      - col`,
			want: Table{
				Columns: []Column{
					{
						Name:     "col",
						Datatype: datatype.Integer,
					},
				},
				Indices: []Index{
					{
						Name:    "idx_col",
						Columns: []string{"col"},
					},
				},
			},
		},
		{
			name: "reference on single column",
			yml: `
columns:
  col:
    type: INT
references:
  foreign_table:
    has_many: true`,
			want: Table{
				Columns: []Column{
					{
						Name:     "col",
						Datatype: datatype.Integer,
					},
				},
				References: []Reference{
					{
						TableName: "foreign_table",
						HasMany:   true,
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
			yml:  `type: enum("red", "blue")`,
			want: Column{
				Datatype: datatype.Enum,
				Params:   []string{"\"red\"", "\"blue\""},
			},
		},
		{
			name: "unsigned int",
			yml: `
type: INT
unsigned: true`,
			want: Column{
				Datatype: datatype.Integer,
				Unsigned: true,
			},
		},
		{
			name: "nullable int",
			yml: `
type: INT
nullable: true`,
			want: Column{
				Datatype: datatype.Integer,
				Nullable: true,
			},
		},
		{
			name: "nullable int",
			yml: `
type: INT
nullable: true`,
			want: Column{
				Datatype: datatype.Integer,
				Nullable: true,
			},
		},
		{
			name: "int with goname",
			yml: `
type: INT
go_name: Col`,
			want: Column{
				Datatype: datatype.Integer,
				GoName:   "Col",
			},
		},
		{
			name: "auto_increment int",
			yml: `
type: INT
primary_key: true
auto_increment: true`,
			want: Column{
				Datatype:      datatype.Integer,
				PrimaryKey:    true,
				AutoIncrement: true,
			},
		},
		{
			name: "auto_increment int",
			yml: `
type: INT
primary_key: true
auto_increment: true`,
			want: Column{
				Datatype:      datatype.Integer,
				PrimaryKey:    true,
				AutoIncrement: true,
			},
		},
		{
			name: "string with collation and charset",
			yml: `
type: text
charset: utf8-mb4
collation: latin`,
			want: Column{
				Datatype:  datatype.Text,
				Charset:   "utf8-mb4",
				Collation: "latin",
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
