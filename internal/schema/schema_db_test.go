package schema

import (
	"reflect"
	"testing"
)

func TestDatabase_GetTable(t *testing.T) {
	type fields struct {
		Tables []Table
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Table
		want1  bool
	}{
		{
			name: "present, single table",
			fields: fields{
				Tables: []Table{{Name: "yes"}},
			},
			args:  args{name: "yes"},
			want:  Table{Name: "yes"},
			want1: true,
		},
		{
			name: "present, first table",
			fields: fields{
				Tables: []Table{{Name: "yes"}, {Name: "no"}},
			},
			args:  args{name: "yes"},
			want:  Table{Name: "yes"},
			want1: true,
		},
		{
			name: "present, first table",
			fields: fields{
				Tables: []Table{{Name: "no"}, {Name: "yes"}},
			},
			args:  args{name: "yes"},
			want:  Table{Name: "yes"},
			want1: true,
		},
		{
			name: "missing",
			fields: fields{
				Tables: []Table{{Name: "no"}},
			},
			args:  args{name: "yes"},
			want:  Table{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Tables: tt.fields.Tables,
			}
			got, got1 := db.GetTable(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTable() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetTable() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
