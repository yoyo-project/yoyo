package schema

import (
	"reflect"
	"testing"
)

func TestTable_ExportedGoName(t1 *testing.T) {
	type fields struct {
		Name   string
		GoName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "from simple name",
			fields: fields{
				Name: "test",
			},
			want: "Test",
		},
		{
			name: "name with underscores",
			fields: fields{
				Name: "test_table",
			},
			want: "TestTable",
		},
		{
			name: "name already PascalCase",
			fields: fields{
				Name: "TestTable",
			},
			want: "TestTable",
		},
		{
			name: "name with dashes",
			fields: fields{
				Name: "test-table",
			},
			want: "TestTable",
		},
		{
			name: "simple GoName",
			fields: fields{
				GoName: "GoTable",
			},
			want: "GoTable",
		},
		{
			name: "GoName overrides Name",
			fields: fields{
				Name:   "no",
				GoName: "GoTable",
			},
			want: "GoTable",
		},
		{
			name: "GoName ignores dashes",
			fields: fields{
				Name:   "no",
				GoName: "Go-Table",
			},
			want: "GoTable",
		},
		{
			name: "GoName forces pascal",
			fields: fields{
				Name:   "no",
				GoName: "goTable",
			},
			want: "GoTable",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:   tt.fields.Name,
				GoName: tt.fields.GoName,
			}
			if got := t.ExportedGoName(); got != tt.want {
				t1.Errorf("ExportedGoName()\nwant %#v\n got %#v", tt.want, got)
			}
		})
	}
}

func TestTable_GetColumn(t1 *testing.T) {
	type fields struct {
		Columns []Column
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Column
		want1  bool
	}{
		{
			name: "match with single column",
			fields: fields{Columns: []Column{
				{Name: "col"},
			}},
			args:  args{name: "col"},
			want:  Column{Name: "col"},
			want1: true,
		},
		{
			name: "match with first of two columns",
			fields: fields{Columns: []Column{
				{Name: "yes"},
				{Name: "no"},
			}},
			args:  args{name: "yes"},
			want:  Column{Name: "yes"},
			want1: true,
		},
		{
			name: "match with second of two columns",
			fields: fields{Columns: []Column{
				{Name: "no"},
				{Name: "yes"},
			}},
			args:  args{name: "yes"},
			want:  Column{Name: "yes"},
			want1: true,
		},
		{
			name:   "miss on empty table",
			fields: fields{Columns: []Column{}},
			args:   args{name: "yes"},
			want1:  false,
		},
		{
			name: "miss on table with no matching column",
			fields: fields{Columns: []Column{
				{Name: "no"},
				{Name: "also_no"},
			}},
			args:  args{name: "yes"},
			want1: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Columns: tt.fields.Columns,
			}
			got, got1 := t.GetColumn(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetColumn()\nwant %#v\n got %v", tt.want, got)
			}
			if got1 != tt.want1 {
				t1.Errorf("GetColumn() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTable_QueryPackageName(t1 *testing.T) {
	type fields struct {
		Name   string
		GoName string
	}
	tests := []struct {
		name     string
		fields   fields
		wantName string
	}{
		{
			name: "from single word name",
			fields: fields{
				Name: "table",
			},
			wantName: "table",
		},
		{
			name: "from two-word underscored name",
			fields: fields{
				Name: "some_table",
			},
			wantName: "some_table",
		},
		{
			name: "from two-word Pascal name",
			fields: fields{
				Name: "SomeTable",
			},
			wantName: "some_table",
		},
		{
			name: "from two-word dashe-case name",
			fields: fields{
				Name: "some-table",
			},
			wantName: "some_table",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:   tt.fields.Name,
				GoName: tt.fields.GoName,
			}
			if gotName := t.QueryPackageName(); gotName != tt.wantName {
				t1.Errorf("QueryPackageName()\nwant= %v\n got %v", tt.wantName, gotName)
			}
		})
	}
}

func TestTable_PKColumns(t1 *testing.T) {
	type fields struct {
		Columns []Column
	}
	tests := []struct {
		name     string
		fields   fields
		wantCols []Column
	}{
		{
			name: "single pk column",
			fields: fields{Columns: []Column{
				{Name: "id", PrimaryKey: true},
			}},
			wantCols: []Column{
				{Name: "id", PrimaryKey: true},
			},
		},
		{
			name: "single pk out of two columns",
			fields: fields{Columns: []Column{
				{Name: "id", PrimaryKey: true},
				{Name: "col", PrimaryKey: false},
			}},
			wantCols: []Column{
				{Name: "id", PrimaryKey: true},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Columns: tt.fields.Columns,
			}
			if gotCols := t.PKColumns(); !reflect.DeepEqual(gotCols, tt.wantCols) {
				t1.Errorf("PKColumns()\nwant: %#v,\n got: %#v", tt.wantCols, gotCols)
			}
		})
	}
}
