package migration

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/yoyo-project/yoyo/internal/reverse"
	"github.com/yoyo-project/yoyo/internal/yoyo"

	"github.com/yoyo-project/yoyo/internal/datatype"
	"github.com/yoyo-project/yoyo/internal/schema"
)

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
			name: "AddAll with no existing indices",
			fields: fields{
				options: AddAll,
			},
			args: args{
				tableName: "myTable",
				table:     schema.Table{},
			},
		},
		{
			name: "AddMissing with existing index",
			fields: fields{
				options:         AddMissing,
				existingIndices: []string{"somewhere"},
			},
			args: args{
				tableName: "myTable",
				table: schema.Table{
					Indices: []schema.Index{
						{
							Name: "somewhere",
						},
					},
				},
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
					Columns: []schema.Column{
						{
							Name:          "id",
							Datatype:      datatype.Integer,
							Unsigned:      true,
							Nullable:      false,
							AutoIncrement: true,
							PrimaryKey:    true,
						},
					},
					Indices: []schema.Index{
						{
							Name: "primary",
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
			d := &mockAdapter{}
			f := NewIndexAdder(d, tt.fields.options, func(_, index string) bool {
				for _, s := range tt.fields.existingIndices {
					if s == index {
						return true
					}
				}
				return false
			})

			err := f(tt.args.table, &sb)
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
			name: "addMissing with no existingColumns",
			fields: fields{
				options: AddMissing,
			},
			args: args{
				tableName: "myTable",
				table: schema.Table{
					Columns: []schema.Column{
						{
							Name: "asdf",
						},
					},
				},
			},
			want: "\n",
		},
		{
			name: "addMissing with existing column",
			fields: fields{
				options:         AddMissing,
				existingColumns: []string{"asdf"},
			},
			args: args{
				tableName: "myTable",
				table: schema.Table{
					Columns: []schema.Column{
						{
							Name: "asdf",
						},
					},
				},
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
					Columns: []schema.Column{
						{
							Name:     "id",
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
					Columns: []schema.Column{
						{
							Name:     "id",
							Datatype: datatype.Integer,
							Unsigned: true,
						},
						{
							Name:     "blah",
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
			a := &mockAdapter{}
			f := NewColumnAdder(a, tt.fields.options, func(_, column string) bool {
				for _, s := range tt.fields.existingColumns {
					if s == column {
						return true
					}
				}
				return false
			})

			err := f(tt.args.table, &sb)
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

func TestNewRefAdder(t *testing.T) {
	type args struct {
		localTable string
		refs       []schema.Reference
	}
	type fields struct {
		options      uint8
		existingRefs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name: "AddAll one column",
			fields: fields{
				options: AddAll,
			},
			args: args{
				localTable: "one",
				refs: []schema.Reference{
					{
						TableName: "one",
						HasOne:    true,
					},
				},
			},
			want: "\n",
		},
		{
			name: "AddAll one column nonexistant table",
			fields: fields{
				options: AddAll,
			},
			args: args{
				localTable: "one",
				refs: []schema.Reference{
					{
						TableName: "no",
						HasOne:    true,
					},
				},
			},
			wantErr: "does not exist",
		},
		{
			name: "AddMissing with missing ref",
			fields: fields{
				options: AddMissing,
			},
			args: args{
				localTable: "one",
				refs: []schema.Reference{
					{
						TableName: "one",
						HasMany:   true,
					},
				},
			},
			want: "\n",
		},
		{
			name: "AddMissing with no missing refs",
			fields: fields{
				options:      AddMissing,
				existingRefs: []string{"one"},
			},
			args: args{
				localTable: "one",
				refs: []schema.Reference{
					{
						TableName: "one",
						HasOne:    true,
					},
				},
			},
		},
		{
			name: "addAll two refs",
			fields: fields{
				options: AddAll,
			},
			args: args{
				localTable: "one",
				refs: []schema.Reference{
					{
						TableName: "one",
						HasOne:    true,
					},
					{
						TableName: "two",
						HasOne:    true,
					},
				},
			},
			want: "\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := strings.Builder{}
			a := &mockAdapter{}
			db := schema.Database{
				Tables: []schema.Table{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
				},
			}

			f := NewRefAdder(a, db, tt.fields.options, func(_, column string) bool {
				for _, s := range tt.fields.existingRefs {
					if s == column {
						return true
					}
				}
				return false
			})

			err := f(tt.args.localTable, tt.args.refs, &sb)

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

func TestNewGenerator(t *testing.T) {
	const (
		callCreateTable       = "callCreateTable"
		callAddMissingColumns = "callAddMissingColumns"
		callAddMissingIndices = "callAddMissingIndices"
		callAddAllIndices     = "callAddAllIndices"
		callHasTable          = "callHasTable"
		callAddMissingRefs    = "callAddMissingRefs"
		callAddAllRefs        = "callAddAllRefs"
	)

	type args struct {
		db schema.Database
	}
	tests := []struct {
		name             string
		args             args
		wantCallsInOrder []string // the call order
		existingTables   []string
		errorOnCall      int
		wantErr          bool
	}{
		{
			name: "simple table on new database",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			wantCallsInOrder: []string{
				callHasTable,
				callCreateTable,
				callAddAllIndices,
				callAddAllRefs,
			},
		},
		{
			name: "no new tables",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			existingTables: []string{"table"},
			wantCallsInOrder: []string{
				callHasTable,
				callAddMissingColumns,
				callAddMissingIndices,
				callAddMissingRefs,
			},
		},
		{
			name: "error finding existing tables",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			wantCallsInOrder: []string{
				callHasTable,
			},
			wantErr:     true,
			errorOnCall: 1,
		},
		{
			name: "error creating table",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			wantCallsInOrder: []string{
				callHasTable,
				callCreateTable,
			},
			wantErr:     true,
			errorOnCall: 2,
		},
		{
			name: "error adding all indices",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			wantCallsInOrder: []string{
				callHasTable,
				callCreateTable,
				callAddAllIndices,
			},
			wantErr:     true,
			errorOnCall: 3,
		},
		{
			name: "error adding missing cols",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			existingTables: []string{"table"},
			wantCallsInOrder: []string{
				callHasTable,
				callAddMissingColumns,
			},
			wantErr:     true,
			errorOnCall: 2,
		},
		{
			name: "error adding missing indices",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			existingTables: []string{"table"},
			wantCallsInOrder: []string{
				callHasTable,
				callAddMissingColumns,
				callAddMissingIndices,
			},
			wantErr:     true,
			errorOnCall: 3,
		},
		{
			name: "error adding missing refs",
			args: args{db: schema.Database{
				Tables: []schema.Table{
					{
						Name: "table",
						Columns: []schema.Column{
							{
								Name: "column",
							},
						},
					},
				},
			}},
			existingTables: []string{"table"},
			wantCallsInOrder: []string{
				callHasTable,
				callAddMissingColumns,
				callAddMissingIndices,
				callAddMissingRefs,
			},
			wantErr:     true,
			errorOnCall: 4,
		},
	}

	for _, tt := range tests {
		var gotCallsInOrder []string
		var (
			w           io.StringWriter = &strings.Builder{}
			createTable TableGenerator  = func(table schema.Table, sw io.StringWriter) (err error) {
				gotCallsInOrder = append(gotCallsInOrder, callCreateTable)
				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return err
			}
			addMissingColumns TableGenerator = func(table schema.Table, sw io.StringWriter) (err error) {
				gotCallsInOrder = append(gotCallsInOrder, callAddMissingColumns)
				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return err
			}
			addMissingIndices TableGenerator = func(table schema.Table, sw io.StringWriter) (err error) {
				gotCallsInOrder = append(gotCallsInOrder, callAddMissingIndices)
				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return err
			}
			addAllIndices TableGenerator = func(table schema.Table, sw io.StringWriter) (err error) {
				gotCallsInOrder = append(gotCallsInOrder, callAddAllIndices)
				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return err
			}
			hasTable StringSearcher = func(s string) (res bool, err error) {
				gotCallsInOrder = append(gotCallsInOrder, callHasTable)
				for i := range tt.existingTables {
					if tt.existingTables[i] == s {
						res = true
						break
					}
				}

				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return res, err
			}
			addMissingRefs RefGenerator = func(localTable string, refs []schema.Reference, sw io.StringWriter) (err error) {
				gotCallsInOrder = append(gotCallsInOrder, callAddMissingRefs)
				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return err
			}
			addAllRefs RefGenerator = func(localTable string, refs []schema.Reference, sw io.StringWriter) (err error) {
				gotCallsInOrder = append(gotCallsInOrder, callAddAllRefs)
				if tt.errorOnCall == len(gotCallsInOrder) {
					err = errors.New("err")
				}
				return err
			}
		)
		t.Run(tt.name, func(t *testing.T) {
			f := NewGenerator(
				createTable,
				addMissingColumns,
				addMissingIndices,
				addAllIndices,
				hasTable,
				addMissingRefs,
				addAllRefs,
			)

			gotErr := f(tt.args.db, w)

			if !reflect.DeepEqual(tt.wantCallsInOrder, gotCallsInOrder) {
				t.Errorf("NewGenerator()\nwant %#v\n got %#v", tt.wantCallsInOrder, gotCallsInOrder)
			}

			if (gotErr != nil) != tt.wantErr {
				t.Errorf("NewGenerator()\nwant error %#v\n got %#v", tt.wantErr, gotErr)
			}
		})
	}
}

func TestInitGeneratorLoader(t *testing.T) {
	type fields struct {
		initReverseAdapter   func(dia string) (adapter reverse.Adapter, err error)
		initMigrationAdapter func(dia string) (a Adapter, err error)
	}
	type args struct {
		config yoyo.Config
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    GeneratorLoader
		wantErr bool
	}{
		{
			name: "simple no errors",
			fields: fields{
				initReverseAdapter:   func(string) (reverse.Adapter, error) { return mockReverseAdapter{}, nil },
				initMigrationAdapter: func(string) (Adapter, error) { return nil, nil },
			},
		},
		{
			name: "error initing reverser",
			fields: fields{
				initReverseAdapter:   func(string) (reverse.Adapter, error) { return nil, errors.New("blah") },
				initMigrationAdapter: func(string) (Adapter, error) { return nil, nil },
			},
			wantErr: true,
		},
		{
			name: "error initing migrator",
			fields: fields{
				initReverseAdapter:   func(string) (reverse.Adapter, error) { return mockReverseAdapter{}, nil },
				initMigrationAdapter: func(string) (Adapter, error) { return nil, errors.New("blah") },
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newGenerator := func(
				TableGenerator,
				TableGenerator,
				TableGenerator,
				TableGenerator,
				StringSearcher,
				RefGenerator,
				RefGenerator,
			) Generator {
				return nil
			}
			f := InitGeneratorLoader(tt.fields.initReverseAdapter, tt.fields.initMigrationAdapter, newGenerator)

			_, err := f(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitGeneratorLoader() want error: %#v, got %#v", tt.wantErr, err)
			}
		})
	}
}

func TestNewTableAdder(t *testing.T) {
	type args struct {
		a Adapter
	}
	tests := []struct {
		name string
		args args
		want TableGenerator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTableAdder(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTableAdder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTableAdder1(t *testing.T) {
	type args struct {
		t schema.Table
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "simple create table",
			args: args{
				t: schema.Table{},
			},
			want: "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewTableAdder(&mockAdapter{})
			sb := strings.Builder{}

			if gotErr := f(tt.args.t, &sb); tt.wantErr != (gotErr != nil) {
				t.Errorf("NewTableAdder() want error %v, got %#v", tt.wantErr, gotErr)
			}

			if got := tt.want; got != sb.String() {
				t.Errorf("NewTableAdder()\nwant %#v\n got %#v", tt.want, got)
			}

		})
	}
}
