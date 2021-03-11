package mysql

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/yoyo-project/yoyo/internal/dbms/base"
)

func Test_adapter_PreparedStatementPlaceholders(t *testing.T) {
	type fields struct {
		db        *sql.DB
		Base      base.Base
		validator validator
	}
	type args struct {
		count int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "single column",
			args: args{count: 1},
			want: []string{"?"},
		},
		{
			name: "two columns",
			args: args{count: 2},
			want: []string{"?", "?"},
		},
		{
			name: "three columns",
			args: args{count: 3},
			want: []string{"?", "?", "?"},
		},
		{
			name: "zero columns",
			args: args{count: 0},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &adapter{
				db:        tt.fields.db,
				Base:      tt.fields.Base,
				validator: tt.fields.validator,
			}
			if got := a.PreparedStatementPlaceholders(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreparedStatementPlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}
