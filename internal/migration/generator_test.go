package migration

import (
	"strings"
	"testing"

	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
)

func TestNewSchemaGenerator(t *testing.T) {
	//type args struct {
	//	conn              *sql.DB
	//	dialect                 dialect.Base
	//	addMissingColumns func(name string, table schema.Table, sw io.StringWriter) error
	//	addMissingIndices func(name string, table schema.Table, sw io.StringWriter) error
	//	addIndices        func(name string, table schema.Table, sw io.StringWriter) error
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	want func(dbms schema.Database, w io.StringWriter) error
	//}{
	//	// TODO: Add test cases.
	//}
	//
	//generate := NewGenerator()
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := NewGenerator(tt.args.conn, tt.args.dialect, tt.args.addMissingColumns, tt.args.addMissingIndices, tt.args.addIndices); !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("NewGenerator() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func TestNewIndexAdder(t *testing.T) {
	type fields struct {
		options         uint8
		existingIndices []string
	}
	type args struct {
		tableName string
		table     schema.Table
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name: "addAll with no existing indices",
			fields: fields{
				options: AddAll,
			},
			args: args{
				tableName: "myTable",
				table:     schema.Table{},
			},
		},
		{
			name: "addAll with single primary key",
			fields: fields{
				options: AddAll,
			},
			args: args{
				tableName: "myTable",
				table: schema.Table{
					Columns: map[string]schema.Column{
						"id": {
							Datatype:      datatype.Integer,
							Unsigned:      true,
							Nullable:      false,
							AutoIncrement: true,
							PrimaryKey:    true,
						},
					},
					Indices: map[string]schema.Index{
						"primary": {
							Columns: []string{
								"id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := strings.Builder{}
			d := &mockDialect{}
			f := NewIndexAdder(d, tt.fields.options, func(_, index string) bool {
				for _, s := range tt.fields.existingIndices {
					if s == index {
						return true
					}
				}
				return false
			})

			err := f(tt.args.tableName, tt.args.table, &sb)
			if err != nil && !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("wanted eror: '%s', got %s", tt.wantErr, err)
			} else if err == nil && len(tt.wantErr) > 0 {
				t.Fatalf("wanted error '%s', got no error", tt.wantErr)
			}

			got := sb.String()
			if got != tt.want {
				t.Fatalf("Wanted string '%s', got '%s'", tt.want, got)
			}
		})
	}
}

func TestNewColumnAdder(t *testing.T) {
	type fields struct {
		options         uint8
		existingColumns []string
	}
	type args struct {
		tableName string
		table     schema.Table
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name: "addAll with no existingColumns",
			fields: fields{
				options: AddAll,
			},
			args: args{
				tableName: "myTable",
				table:     schema.Table{},
			},
		},
		{
			name: "addAll with single column",
			fields: fields{
				options: AddAll,
			},
			want: "\n",
			args: args{
				tableName: "myTable",
				table: schema.Table{
					Columns: map[string]schema.Column{
						"id": {
							Datatype: datatype.Integer,
							Unsigned: true,
						},
					},
				},
			},
		},
		{
			name: "addMissing with all existingColumns missing (identical to addAll)",
			want: "\n\n",
			fields: fields{
				options: AddAll,
			},
			args: args{
				tableName: "myTable",
				table: schema.Table{
					Columns: map[string]schema.Column{
						"id": {
							Datatype: datatype.Integer,
							Unsigned: true,
						},
						"blah": {
							Datatype: datatype.Text,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := strings.Builder{}
			d := &mockDialect{}
			f := NewColumnAdder(d, tt.fields.options, func(_, column string) bool {
				for _, s := range tt.fields.existingColumns {
					if s == column {
						return true
					}
				}
				return false
			})

			err := f(tt.args.tableName, tt.args.table, &sb)
			if err != nil && !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("wanted eror: '%s', got %s", tt.wantErr, err)
			} else if err == nil && len(tt.wantErr) > 0 {
				t.Fatalf("wanted error '%s', got no error", tt.wantErr)
			}

			got := sb.String()
			if got != tt.want {
				t.Fatalf("Wanted string '%s', got '%s'", tt.want, got)
			}
		})
	}
}
