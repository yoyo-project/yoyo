package repository

import (
	"reflect"
	"testing"

	"github.com/yoyo-project/yoyo/internal/dbms/dialect"
	"github.com/yoyo-project/yoyo/internal/dbms/mysql"
	"github.com/yoyo-project/yoyo/internal/dbms/postgres"
)

func TestLoadAdapter(t *testing.T) {
	type args struct {
		dia string
	}
	tests := []struct {
		name        string
		args        args
		wantAdapter Adapter
		wantErr     bool
	}{
		{
			name:        "mysql",
			args:        args{dia: dialect.MySQL},
			wantAdapter: mysql.NewAdapter(),
			wantErr:     false,
		},
		{
			name:        "postgres",
			args:        args{dia: dialect.PostgreSQL},
			wantAdapter: postgres.NewAdapter(),
			wantErr:     false,
		},
		{
			name:    "n/a",
			args:    args{dia: "n/a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdapter, err := LoadAdapter(tt.args.dia)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAdapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAdapter, tt.wantAdapter) {
				t.Errorf("LoadAdapter() gotAdapter = %v, want %v", gotAdapter, tt.wantAdapter)
			}
		})
	}
}
