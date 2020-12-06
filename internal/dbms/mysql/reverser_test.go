package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
	"reflect"
	"strings"
	"testing"
)

func TestInitNewReverser(t *testing.T) {
	type args struct {
		host     string
		user     string
		dbname   string
		password string
		port     string
	}
	tests := []struct {
		name      string
		open      func(driver, dsn string) (*sql.DB, error)
		args      args
		wantError bool
	}{
		{
			name: "all blank",
			open: func(_, _ string) (*sql.DB, error) {
				db, _, _ := sqlmock.New()
				return db, nil
			},
		},
		{
			name: "with port",
			args: args{
				port: "123",
			},
			open: func(_, _ string) (*sql.DB, error) {
				db, _, _ := sqlmock.New()
				return db, nil
			},
		},
		{
			name: "with error",
			open: func(_, _ string) (*sql.DB, error) {
				return nil, errors.New("blah blah")
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitNewReverser(tt.open)
			_, err := got(tt.args.host, tt.args.user, tt.args.dbname, tt.args.password, tt.args.port)

			if err != nil && tt.wantError == false {
				t.Errorf("got unexpected error: %s", err)
			}
		})
	}
}

func Test_reverser_ListTables(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr string
	}{
		{
			name: "single table",
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SHOW TABLES").
						WillReturnRows(mock.NewRows([]string{"Tables_in_db"}).
							AddRow("table"))
					return db
				}(),
			},
			want: []string{"table"},
		},
		{
			name: "four tables",
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SHOW TABLES").
						WillReturnRows(mock.NewRows([]string{"Tables_in_db"}).
							AddRow("table1").
							AddRow("table2").
							AddRow("table3").
							AddRow("table4"))
					return db
				}(),
			},
			want: []string{"table1", "table2", "table3", "table4"},
		},
		{
			name: "zero tables",
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SHOW TABLES").
						WillReturnRows(mock.NewRows([]string{"Tables_in_db"}))
					return db
				}(),
			},
		},
		{
			name: "wrong number of columns",
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SHOW TABLES").
						WillReturnRows(mock.NewRows([]string{"Tables_in_db", "bonus_col"}).
							AddRow("table", "bonus_val"))
					return db
				}(),
			},
			wantErr: "unable to scan table results",
		},
		{
			name: "query error",
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectQuery("SHOW TABLES").
						WillReturnError(fmt.Errorf("oh no it broke"))
					return db
				}(),
			},
			wantErr: "unable to list tables",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.ListTables()
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverser_ListIndices(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		table string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    []string
		wantErr string
	}{
		{
			name: "single index",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listIndicesQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"INDEX_NAME"}).
							AddRow("index"))
					return db
				}(),
			},
			want: []string{"index"},
		},
		{
			name: "three indices",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listIndicesQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"INDEX_NAME"}).
							AddRow("index1").AddRow("index2").AddRow("index3"))
					return db
				}(),
			},
			want: []string{"index1", "index2", "index3"},
		},
		{
			name: "wrong number of columns",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listIndicesQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"INDEX_NAME", "bonus_col"}).
							AddRow("table", "bonus_val"))
					return db
				}(),
			},
			wantErr: "unable to scan",
		},
		{
			name: "query error",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listIndicesQuery, "table")).
						WillReturnError(fmt.Errorf("oh no it broke"))
					return db
				}(),
			},
			wantErr: "unable to list indices",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.ListIndices(tt.args.table)
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverser_ListColumns(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		table string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    []string
		wantErr string
	}{
		{
			name: "single column",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listColumnsQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME"}).
							AddRow("col"))
					return db
				}(),
			},
			want: []string{"col"},
		},
		{
			name: "wrong number of columns",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listColumnsQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "bonus_col"}).
							AddRow("col", "bonus_val"))
					return db
				}(),
			},
			wantErr: "unable to scan",
		},
		{
			name: "query error",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listColumnsQuery, "table")).
						WillReturnError(fmt.Errorf("oh no it broke"))
					return db
				}(),
			},
			wantErr: "unable to list indices",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.ListColumns(tt.args.table)
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverser_ListReferences(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		table string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    []string
		wantErr string
	}{
		{
			name: "single reference",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listReferencesQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"REFERENCED_TABLE_NAME"}).
							AddRow("table2"))
					return db
				}(),
			},
			want: []string{"table2"},
		},
		{
			name: "wrong number of columns",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listReferencesQuery, "table")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "bonus_col"}).
							AddRow("table2", "bonus_val"))
					return db
				}(),
			},
			wantErr: "unable to scan",
		},
		{
			name: "query error",
			args: args{
				table: "table",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(listReferencesQuery, "table")).
						WillReturnError(fmt.Errorf("oh no it broke"))
					return db
				}(),
			},
			wantErr: "unable to list indices",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.ListReferences(tt.args.table)
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverser_GetColumn(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		table   string
		colName string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    schema.Column
		wantErr string
	}{
		{
			name: "basic int id column",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("id", "INT(11)", "NO", "PRI", nil, "auto_increment"))
					return db
				}(),
			},
			want: schema.Column{
				Datatype:      datatype.Integer,
				Precision:     11,
				PrimaryKey:    true,
				AutoIncrement: true,
			},
		},
		{
			name: "unsigned int",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "INT(11) UNSIGNED", "NO", "", nil, ""))
					return db
				}(),
			},
			want: schema.Column{
				Datatype:  datatype.Integer,
				Precision: 11,
				Unsigned:  true,
			},
		},
		{
			name: "DECIMAL(5,3)",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "DECIMAL(5,3)", "NO", "", nil, ""))
					return db
				}(),
			},
			want: schema.Column{
				Datatype:  datatype.Decimal,
				Precision: 5,
				Scale:     3,
			},
		},
		{
			name: "unsigned int",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "DECIMAL(5,3)", "NO", "", nil, ""))
					return db
				}(),
			},
			want: schema.Column{
				Datatype:  datatype.Decimal,
				Precision: 5,
				Scale:     3,
			},
		},
		{
			name: "query error",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnError(fmt.Errorf("uh oh"))
					return db
				}(),
			},
			wantErr: "unable to get column information for",
		},
		{
			name: "no result",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}))
					return db
				}(),
			},
			wantErr: "unable to get column, empty result",
		},
		{
			name: "too may columns",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra", "BORK"}).
							AddRow("col", "DECIMAL(5,3)", "NO", "", nil, "", "BORK"))
					return db
				}(),
			},
			wantErr: "unable to scan result",
		},
		{
			name: "too may rows",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "DECIMAL(5,3)", "NO", "", nil, "").
							AddRow("col", "DECIMAL(5,3)", "NO", "", nil, "").
							CloseError(fmt.Errorf("oh no")))
					return db
				}(),
			},
			wantErr: "unable to close",
		},
		{
			name: "unknown datatype",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "WAT", "NO", "", nil, ""))
					return db
				}(),
			},
			wantErr: "unable to determine datatype",
		},
		{
			name: "malformed precision/length",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "DECIMAL(asd,1)", "NO", "", nil, ""))
					return db
				}(),
			},
			wantErr: "unable to determine precision",
		},
		{
			name: "malformed scale",
			args: args{
				table:   "table",
				colName: "id",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getColumnQuery, "table", "id")).
						WillReturnRows(mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
							AddRow("col", "DECIMAL(4,asd)", "NO", "", nil, ""))
					return db
				}(),
			},
			wantErr: "unable to determine scale",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.GetColumn(tt.args.table, tt.args.colName)
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverser_GetIndex(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		table     string
		indexName string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    schema.Index
		wantErr string
	}{
		{
			name: "basic non-unique index",
			args: args{
				table:     "table",
				indexName: "col",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getIndexQuery, "table", "col")).
						WillReturnRows(mock.NewRows([]string{"NOT NON_UNIQUE", "COLUMN_NAME"}).
							AddRow(0, "col"))
					return db
				}(),
			},
			want: schema.Index{
				Unique:  false,
				Columns: []string{"col"},
			},
		},
		{
			name: "two-column non-unique index",
			args: args{
				table:     "table",
				indexName: "col",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getIndexQuery, "table", "col")).
						WillReturnRows(mock.NewRows([]string{"NOT NON_UNIQUE", "COLUMN_NAME"}).
							AddRow(0, "col").AddRow(0, "col2"))
					return db
				}(),
			},
			want: schema.Index{
				Unique:  false,
				Columns: []string{"col", "col2"},
			},
		},
		{
			name: "basic unique index",
			args: args{
				table:     "table",
				indexName: "col",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getIndexQuery, "table", "col")).
						WillReturnRows(mock.NewRows([]string{"NOT NON_UNIQUE", "COLUMN_NAME"}).
							AddRow(1, "col"))
					return db
				}(),
			},
			want: schema.Index{
				Unique:  true,
				Columns: []string{"col"},
			},
		},
		{
			name: "query error",
			args: args{
				table:     "table",
				indexName: "col",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getIndexQuery, "table", "col")).
						WillReturnError(fmt.Errorf("oh no"))
					return db
				}(),
			},
			wantErr: "unable to get information for index",
		},
		{
			name: "query error",
			args: args{
				table:     "table",
				indexName: "col",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getIndexQuery, "table", "col")).
						WillReturnRows(mock.NewRows([]string{"NOT NON_UNIQUE", "COLUMN_NAME", "EXTRA_COL"}).
							AddRow(1, "col", "aaaaah"))
					return db
				}(),
			},
			wantErr: "unable to scan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.GetIndex(tt.args.table, tt.args.indexName)
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverser_GetReference(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		table         string
		referenceName string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    schema.Reference
		wantErr string
	}{
		{
			name: "non-optional reference with single column",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnRows(mock.NewRows([]string{"UPDATE_RULE", "DELETE_RULE", "CONSTRAINT_NAME"}).
							AddRow("RESTRICT", "RESTRICT", "foreign_fk"))
					mock.ExpectQuery(fmt.Sprintf(getReferenceColumnsQuery, "table", "foreign", "foreign_fk")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "IS_NULLABLE"}).
							AddRow("foreign_id", 0))
					return db
				}(),
			},
			want: schema.Reference{
				OnUpdate:   "RESTRICT",
				OnDelete:   "RESTRICT",
				ColumnName: "foreign_id",
				Optional:   false,
			},
		},
		{
			name: "non-optional reference with two columns",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnRows(mock.NewRows([]string{"UPDATE_RULE", "DELETE_RULE", "CONSTRAINT_NAME"}).
							AddRow("RESTRICT", "RESTRICT", "foreign_fk"))
					mock.ExpectQuery(fmt.Sprintf(getReferenceColumnsQuery, "table", "foreign", "foreign_fk")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "IS_NULLABLE"}).
							AddRow("foreign_id", 0).
							AddRow("foreign_id2", 0))
					return db
				}(),
			},
			want: schema.Reference{
				OnUpdate:    "RESTRICT",
				OnDelete:    "RESTRICT",
				ColumnNames: []string{"foreign_id", "foreign_id2"},
				Optional:    false,
			},
		},
		{
			name: "non-optional reference with single column",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnRows(mock.NewRows([]string{"UPDATE_RULE", "DELETE_RULE", "CONSTRAINT_NAME"}).
							AddRow("RESTRICT", "RESTRICT", "foreign_fk"))
					mock.ExpectQuery(fmt.Sprintf(getReferenceColumnsQuery, "table", "foreign", "foreign_fk")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "IS_NULLABLE"}).
							AddRow("foreign_id", 1))
					return db
				}(),
			},
			want: schema.Reference{
				OnUpdate:   "RESTRICT",
				OnDelete:   "RESTRICT",
				ColumnName: "foreign_id",
				Optional:   true,
			},
		},
		{
			name: "optional reference with two columns",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnRows(mock.NewRows([]string{"UPDATE_RULE", "DELETE_RULE", "CONSTRAINT_NAME"}).
							AddRow("RESTRICT", "RESTRICT", "foreign_fk"))
					mock.ExpectQuery(fmt.Sprintf(getReferenceColumnsQuery, "table", "foreign", "foreign_fk")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "IS_NULLABLE"}).
							AddRow("foreign_id", 1).
							AddRow("foreign_id2", 1))
					return db
				}(),
			},
			want: schema.Reference{
				OnUpdate:    "RESTRICT",
				OnDelete:    "RESTRICT",
				ColumnNames: []string{"foreign_id", "foreign_id2"},
				Optional:    true,
			},
		},
		{
			name: "no columns found",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnRows(mock.NewRows([]string{"UPDATE_RULE", "DELETE_RULE", "CONSTRAINT_NAME"}).
							AddRow("RESTRICT", "RESTRICT", "foreign_fk"))
					mock.ExpectQuery(fmt.Sprintf(getReferenceColumnsQuery, "table", "foreign", "foreign_fk")).
						WillReturnRows(mock.NewRows([]string{"COLUMN_NAME", "IS_NULLABLE"}))
					return db
				}(),
			},
			wantErr: "unable to find any reference columns for table",
		},
		{
			name: "no reference information found",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnError(fmt.Errorf("oh no"))
					return db
				}(),
			},
			wantErr: "unable to get reference information",
		},
		{
			name: "no column information found",
			args: args{
				table:         "table",
				referenceName: "foreign",
			},
			fields: fields{
				db: func() *sql.DB {
					db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					mock.ExpectQuery(fmt.Sprintf(getReferenceQuery, "table", "foreign")).
						WillReturnRows(mock.NewRows([]string{"UPDATE_RULE", "DELETE_RULE", "CONSTRAINT_NAME"}).
							AddRow("RESTRICT", "RESTRICT", "foreign_fk"))
					mock.ExpectQuery(fmt.Sprintf(getReferenceColumnsQuery, "table", "foreign", "foreign_fk")).
						WillReturnError(fmt.Errorf("oh no"))
					return db
				}(),
			},
			wantErr: "unable to get reference columns",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &reverser{
				db: tt.fields.db,
			}
			got, err := d.GetReference(tt.args.table, tt.args.referenceName)
			if (err != nil) && len(tt.wantErr) != 0 && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr == "" {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(tt.wantErr) != 0 {
				t.Errorf("ListTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTables() got = %v, want %v", got, tt.want)
			}
		})
	}
}
