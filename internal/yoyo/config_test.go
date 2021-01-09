package yoyo

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
)

func cfg(wd string) Config {
	conf := Config{
		Paths: Paths{
			Migrations:   fmt.Sprintf("%s/%s", wd, defaultMigrationsPath),
			Repositories: fmt.Sprintf("%s/%s", wd, defaultRepositoryPath),
		},
		Schema: schema.Database{
			Dialect: "mysql",
			Tables: map[string]schema.Table{
				"primary": {
					Columns: map[string]schema.Column{
						"id": {
							Datatype: datatype.Integer,
							Default: func() *string {
								var s = "0"
								return &s
							}(),
						},
						"secondary_id": {
							Datatype: datatype.Integer,
						},
					},
				},
				"secondary": {
					Columns: map[string]schema.Column{
						"id": {
							Datatype: datatype.Integer,
						},
					},
				},
			},
		},
	}

	for tn, t := range conf.Schema.Tables {
		t.SetName(tn)
		conf.Schema.Tables[tn] = t

		for cn, c := range t.Columns {
			c.SetName(cn)
			t.Columns[cn] = c
		}
	}

	return conf
}

func TestLoad(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Error("unable to start tests: could not determine cwd")
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	tests := []struct {
		name    string
		wantCfg Config
		wantErr bool
		wd      string
	}{
		{
			name:    "Normal Condition",
			wantCfg: cfg(originalWD),
			wantErr: false,
			wd:      originalWD,
		},
		{
			name:    "WD In Project Child",
			wantCfg: cfg(originalWD),
			wantErr: false,
			wd:      fmt.Sprintf("%s/test_child", originalWD),
		},
		{
			name:    "WD In Project Parent",
			wantCfg: Config{},
			wantErr: true,
			wd:      fmt.Sprintf("%s/..", originalWD),
		},
	}
	for _, tt := range tests {
		if tt.wd != "" {
			_ = os.Chdir(tt.wd)
			wd, _ := os.Getwd()
			fmt.Sprintf("Working Directory: %s", wd)
		} else {
			_ = os.Chdir(originalWD)
		}
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("LoadConfig() gotCfg = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}

}
