package migration

import (
	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
	"github.com/yoyo-project/yoyo/internal/dbms/mysql"
	"github.com/yoyo-project/yoyo/internal/dbms/postgres"
	"reflect"
	"testing"
)

func TestLoadAdapter(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantA   Adapter
		wantErr bool
	}{
		{
			name:  dialect.MySQL,
			args:  args{name: dialect.MySQL},
			wantA: mysql.NewAdapter(),
		},
		{
			name:  dialect.PostgreSQL,
			args:  args{name: dialect.PostgreSQL},
			wantA: postgres.NewAdapter(),
		},
		{
			name:    "nothing",
			args:    args{name: "nothing"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := LoadAdapter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAdapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantA) {
				t.Errorf("LoadAdapter() gotD = %v, want %v", gotD, tt.wantA)
			}
		})
	}
}
